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

// GetSupervisors returns the current supervisors
func (k *Keeper) GetSupervisors(ctx sdk.Context) []string {
	var moderators []string
	k.paramSpace.GetIfExists(ctx, types.ParamSupervisors, &moderators)
	return moderators
}

// SetSupervisors sets the Cerberus owners
func (k *Keeper) SetSupervisors(ctx sdk.Context, supervisors []string) {
	k.paramSpace.Set(ctx, types.ParamSupervisors, &supervisors)
}

func (k *Keeper) GetFixedGasParams(ctx sdk.Context) types.FixedGasParams {
	var out types.FixedGasParams
	k.paramSpace.GetIfExists(ctx, types.ParamFixedGas, &out)
	return out
}

func (k *Keeper) SetFixedGasParams(ctx sdk.Context, out types.FixedGasParams) {
	k.paramSpace.Set(ctx, types.ParamFixedGas, &out)
}
