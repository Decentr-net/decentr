package keeper

import (
	"github.com/Decentr-net/decentr/x/pdv/types"
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

// Sets the entire PDV metadata struct for an address
func (k Keeper) SetPDV(ctx sdk.Context, address string, pdv types.PDV) {
	if pdv.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)

	store.Set([]byte(address), k.cdc.MustMarshalBinaryBare(pdv))
}

// Gets the entire PDV metadata struct for an address
func (k Keeper) GetPDV(ctx sdk.Context, address string) types.PDV {
	store := ctx.KVStore(k.storeKey)

	if !k.IsHashPresent(ctx, address) {
		return types.PDV{}
	}

	bz := store.Get([]byte(address))

	var pdv types.PDV
	k.cdc.MustUnmarshalBinaryBare(bz, &pdv)
	return pdv
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetPDV(ctx, name).Owner
}

// Check if the address is present in the store or not
func (k Keeper) IsHashPresent(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(address))
}

// Get an iterator over all PDVs in which the keys are the address and the values are the PDV
func (k Keeper) GetPDVsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}