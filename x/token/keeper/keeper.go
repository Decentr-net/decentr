package keeper

import (
	"fmt"

	"github.com/Decentr-net/decentr/x/token/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/tendermint/tendermint/libs/log"
)

var totalSupplyKey = []byte("totalSupply")

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey      sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc           *codec.Codec // The wire codec for binary encoding/decoding.
	accountKeeper keeper.AccountKeeper
}

// NewKeeper creates new instances of the token Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accountKeeper keeper.AccountKeeper) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		accountKeeper: accountKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// AddTokens adds token to the given owner
func (k Keeper) AddTokens(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int) {
	var balance sdk.Int

	if k.IsInitialBalanceSet(ctx, owner) {
		balance = k.GetBalance(ctx, owner).Add(amount)
	} else {
		k.Logger(ctx).Info(
			"account initialized with init balance",
			"account", owner.String())
		balance = initialTokenBalance().Add(amount)
	}

	k.SetBalance(ctx, owner, balance)
}

// addTotalSupply increase or decrease total supply with the given amount of tokens
func (k Keeper) addTotalSupply(ctx sdk.Context, amount sdk.Int) {
	balance := k.GetTotalSupply(ctx)
	balance = balance.Add(amount)
	ctx.KVStore(k.storeKey).Set(totalSupplyKey, k.cdc.MustMarshalBinaryBare(balance))
}

// GetBalance returns token balance for the given address
func (k Keeper) GetBalance(ctx sdk.Context, address sdk.AccAddress) sdk.Int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)
	if store.Has(address) {
		balance := store.Get(address)
		var amount sdk.Int
		k.cdc.MustUnmarshalBinaryBare(balance, &amount)
		return amount
	}

	// if account exists, but initial token is not set
	if k.accountKeeper.GetAccount(ctx, address) != nil {
		return initialTokenBalance()
	}

	return sdk.ZeroInt()
}

func initialTokenBalance() sdk.Int {
	return sdk.NewInt(1 * types.Denominator)
}

func (k Keeper) IsInitialBalanceSet(ctx sdk.Context, owner sdk.AccAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)
	return store.Has(owner)
}

// SetBalance set balance to the given user
func (k Keeper) SetBalance(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)
	store.Set(owner, k.cdc.MustMarshalBinaryBare(amount))
	k.addTotalSupply(ctx, amount)
}

// GetBalanceIterator gets an iterator over all balances in which the keys are the accounts and the values are their balance
func (k Keeper) GetBalanceIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)
	return sdk.KVStorePrefixIterator(store, nil)
}

// GetTotalSupply returns total token supply
func (k Keeper) GetTotalSupply(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	totalSupply := store.Get(totalSupplyKey)
	if totalSupply == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.Int
	k.cdc.MustUnmarshalBinaryBare(totalSupply, &amount)
	return amount
}
