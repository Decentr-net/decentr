package keeper

import (
	"github.com/Decentr-net/decentr/x/decentr/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	CoinKeeper types.BankKeeper
	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the PDV Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, coinKeeper types.BankKeeper) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		CoinKeeper: coinKeeper,
	}
}

// Sets the entire PDV metadata struct for a name
func (k Keeper) SetPDV(ctx sdk.Context, name string, pdv types.PDV) {
	if pdv.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)

	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(pdv))
}

// Gets the entire PDV metadata struct for a name
func (k Keeper) GetPDV(ctx sdk.Context, name string) types.PDV {
	store := ctx.KVStore(k.storeKey)

	if !k.IsNamePresent(ctx, name) {
		return types.PDV{}
	}

	bz := store.Get([]byte(name))

	var pdv types.PDV
	k.cdc.MustUnmarshalBinaryBare(bz, &pdv)
	return pdv
}

// Check if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

// Get an iterator over all names in which the keys are the names and the values are the PDV
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
