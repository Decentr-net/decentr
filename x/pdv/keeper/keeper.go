package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/pdv/types"
)

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

// GetCerberusOwners returns the current Cerberus owners
func (k *Keeper) GetCerberusOwners(ctx sdk.Context) []string {
	var moderators []string
	k.paramSpace.GetIfExists(ctx, types.ParamCerberusOwnersKey, &moderators)
	return moderators
}

// SetCerberusOwners sets the Cerberus owners
func (k *Keeper) SetCerberusOwners(ctx sdk.Context, moderators []string) {
	k.paramSpace.Set(ctx, types.ParamCerberusOwnersKey, &moderators)
}
