package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cosmos/cosmos-sdk/store/prefix"

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

	paramSpace params.Subspace
}

// NewKeeper creates new instances of the PDV Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace params.Subspace) Keeper {
	ps := paramSpace.WithKeyTable(types.ParamKeyTable())

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: ps,
	}
}

// Sets the entire PDV metadata struct for an address
func (k Keeper) SetPDV(ctx sdk.Context, address string, pdv types.PDV) {
	if pdv.Owner.Empty() {
		return
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)
	marshalled := k.cdc.MustMarshalBinaryBare(pdv)
	store.Set([]byte(address), marshalled)

	index := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexPrefix)
	indexKey := append(pdv.Owner, utils.Uint64ToBytes(pdv.Timestamp)...)
	index.Set(indexKey, k.cdc.MustMarshalBinaryBare(statsItem{
		Address: pdv.Address,
		Type:    pdv.Type,
	}))
}

// Gets the entire PDV metadata struct for an address
func (k Keeper) GetPDV(ctx sdk.Context, address string) types.PDV {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix)

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
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix).Has([]byte(address))
}

// Get an iterator over all PDVs in which the keys are the address and the values are the PDV
func (k Keeper) GetPDVsIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(prefix.NewStore(ctx.KVStore(k.storeKey), types.StorePrefix), nil)
}

func (k *Keeper) ListPDV(ctx sdk.Context, owner sdk.AccAddress, from *uint64, limit uint) []types.PDV {
	index := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndexPrefix)

	it := sdk.KVStoreReversePrefixIterator(prefix.NewStore(index, owner), nil)
	defer it.Close()

	if from != nil {
		t := utils.Uint64ToBytes(*from)
		for ; it.Valid() && bytes.Compare(it.Key(), t) > -1; it.Next() {
		}
	}

	out := make([]types.PDV, 0)

	for i := uint(0); i < limit && it.Valid(); i++ {
		var si statsItem
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &si)

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

// GetCerberusAddr returns the current Cerberus address
func (k *Keeper) GetCerberusAddr(ctx sdk.Context) string {
	var addr string
	k.paramSpace.GetIfExists(ctx, types.ParamCerberusAddressKey, &addr)
	return addr
}

// SetCerberusAddr sets the Cerberus address
func (k *Keeper) SetCerberusAddr(ctx sdk.Context, addr string) {
	k.paramSpace.Set(ctx, types.ParamCerberusAddressKey, &addr)
}
