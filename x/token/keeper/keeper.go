package keeper

import (
	"github.com/Decentr-net/decentr/x/token/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

	total := k.GetTotalBalance(ctx)
	maxSupply := sdk.NewIntFromUint64(uint64(types.MaxSupply))
	if maxSupply.LT(total.Add(amount)) {
		// Max supply is reached, stop the emission
		amount = maxSupply.Sub(total)
	}

	balance = balance.Add(amount)
	ctx.KVStore(k.storeKey).Set(owner, k.cdc.MustMarshalBinaryBare(balance))
}

// Gets balance
func (k Keeper) GetBalance(ctx sdk.Context, owner sdk.AccAddress) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	if store.Has(owner) {
		bz := store.Get(owner)

		var amount sdk.Int
		k.cdc.MustUnmarshalBinaryBare(bz, &amount)
		return amount
	}
	return sdk.ZeroInt()
}

// Gets total balance
func (k Keeper) GetTotalBalance(ctx sdk.Context) sdk.Int {
	total := sdk.ZeroInt()

	iterator := k.GetBalanceIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		owner := iterator.Key()
		total = total.Add(k.GetBalance(ctx, owner))
	}

	return total
}

// Get an iterator over all balances in which the keys are the accounts and the values are their balance
func (k Keeper) GetBalanceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
