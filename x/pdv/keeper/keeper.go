package keeper

import (
	"bytes"

	"github.com/Decentr-net/decentr/x/pdv/types"
	"github.com/Decentr-net/decentr/x/utils"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenKeeper interface {
	AddTokens(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Int, description []byte)
}

type statsItem struct {
	Address string        `json:"address"`
	Type    types.PDVType `json:"type"`
}

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	tokens   TokenKeeper
}

// NewKeeper creates new instances of the PDV Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, tokens TokenKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		tokens:   tokens,
	}
}

// Sets the entire PDV metadata struct for an address
func (k Keeper) SetPDV(ctx sdk.Context, address string, pdv types.PDV) {
	if pdv.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)

	marshalled := k.cdc.MustMarshalBinaryBare(pdv)

	store.Set(getStoreKey([]byte(address)), marshalled)
	indexKey := append(pdv.Owner.Bytes(), utils.Uint64ToBytes(pdv.Timestamp)...)
	store.Set(getIndexKey(indexKey), k.cdc.MustMarshalBinaryBare(statsItem{
		Address: pdv.Address,
		Type:    pdv.Type,
	}))

	t := sdk.NewInt(0)
	switch pdv.Type {
	case types.PDVTypeCookie:
		t = sdk.NewInt(1)
	}

	k.tokens.AddTokens(ctx, pdv.Owner, t, utils.GetHash(pdv.Address))
}

// Gets the entire PDV metadata struct for an address
func (k Keeper) GetPDV(ctx sdk.Context, address string) types.PDV {
	store := ctx.KVStore(k.storeKey)

	if !k.IsHashPresent(ctx, address) {
		return types.PDV{}
	}

	bz := store.Get(getStoreKey([]byte(address)))

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
	return store.Has(getStoreKey([]byte(address)))
}

// Get an iterator over all PDVs in which the keys are the address and the values are the PDV
func (k Keeper) GetPDVsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.StorePrefix)
}

func (k *Keeper) ListPDV(ctx sdk.Context, owner sdk.AccAddress, from *uint64, limit uint) []types.PDV {
	store := ctx.KVStore(k.storeKey)

	it := sdk.KVStoreReversePrefixIterator(store, getIndexKey(owner))

	if from != nil {
		t := utils.Uint64ToBytes(*from)
		for ; it.Valid(); it.Next() {
			if bytes.Compare(it.Value(), t) == -1 {
				break
			}
		}
	}

	out := make([]types.PDV, 0)

	for i := uint(0); i < limit && it.Valid(); i++ {
		var si statsItem
		k.cdc.MustUnmarshalJSON(it.Value(), &si)

		out = append(out, types.PDV{
			Timestamp: utils.BytesToUint64(it.Key()),
			Address:   si.Address,
			Owner:     owner,
			Type:      si.Type,
		})

		it.Next()
	}

	return out
}

func getStoreKey(key []byte) []byte {
	return append(types.StorePrefix, key...)
}

func getIndexKey(key []byte) []byte {
	return append(types.IndexPrefix, key...)
}
