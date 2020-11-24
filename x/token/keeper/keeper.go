package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var totalSupplyKey = []byte("totalSupply")

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	stats    Stats
}

// NewKeeper creates new instances of the token Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, stats Stats) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		stats:    stats,
	}
}

// AddTokens adds token to the given owner
func (k Keeper) AddTokens(ctx sdk.Context, owner sdk.AccAddress, timestamp uint64, amount sdk.Int) {
	balance := k.GetBalance(ctx, owner)
	balance = balance.Add(amount)
	ctx.KVStore(k.storeKey).Set(owner, k.cdc.MustMarshalBinaryBare(balance))
	k.addTotalSupply(ctx, amount)
	if err := k.stats.AddToken(owner, timestamp, amount); err != nil {
		ctx.Logger().Error("failed to add tokens to stats", "err", err.Error())
	}
}

// addTotalSupply increase or decrease total supply with the given amount of tokens
func (k Keeper) addTotalSupply(ctx sdk.Context, amount sdk.Int) {
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
