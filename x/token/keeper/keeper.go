package keeper

import (
	"github.com/Decentr-net/decentr/x/token/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var totalSupplyKey = []byte("totalSupply")

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the token Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// AddTokens adds token to the given owner
func (k Keeper) AddTokens(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int) {
	balance := k.GetBalance(ctx, owner)

	total := k.GetTotalSupply(ctx)
	maxSupply := sdk.NewIntFromUint64(uint64(types.MaxSupply))
	if maxSupply.LT(total.Add(amount)) {
		// Max supply is reached, stop the emission
		amount = maxSupply.Sub(total)
	}

	balance = balance.Add(amount)
	ctx.KVStore(k.storeKey).Set(owner, k.cdc.MustMarshalBinaryBare(balance))
	k.IncreaseTotalSupply(ctx, amount)
}

// IncreaseTotalSupply increase total supply with the given amount of tokens
func (k Keeper) IncreaseTotalSupply(ctx sdk.Context, amount sdk.Int) {
	balance := k.GetTotalSupply(ctx)
	balance = balance.Add(amount)
	ctx.KVStore(k.storeKey).Set(totalSupplyKey, k.cdc.MustMarshalBinaryBare(balance))
}

// GetBalance returns token balance for the given owner
func (k Keeper) GetBalance(ctx sdk.Context, owner sdk.AccAddress) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	if store.Has(owner) {
		balance := store.Get(owner)
		var amount sdk.Int
		k.cdc.MustUnmarshalBinaryBare(balance, &amount)
		return amount
	}
	return sdk.ZeroInt()
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

// Get an iterator over all balances in which the keys are the accounts and the values are their balance
func (k Keeper) GetBalanceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
