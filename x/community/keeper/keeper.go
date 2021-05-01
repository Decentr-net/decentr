package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/community/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	tokens   TokenKeeper

	paramSpace params.Subspace
}

// NewKeeper creates new instances of the community Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace params.Subspace, tokens TokenKeeper) Keeper {
	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		tokens:     tokens,
		paramSpace: ps,
	}
}

// GetModerators returns the current moderators
func (k *Keeper) GetModerators(ctx sdk.Context) []string {
	var moderators []string
	k.paramSpace.GetIfExists(ctx, types.ParamModerators, &moderators)
	return moderators
}

// SetModerators sets the moderators
func (k *Keeper) SetModerators(ctx sdk.Context, moderators []string) {
	k.paramSpace.Set(ctx, types.ParamModerators, &moderators)
}

func (k *Keeper) GetFixedGasParams(ctx sdk.Context) types.FixedGasParams {
	var out types.FixedGasParams
	k.paramSpace.GetIfExists(ctx, types.ParamsFixedGas, &out)
	return out
}

func (k *Keeper) SetFixedGasParams(ctx sdk.Context, out types.FixedGasParams) {
	k.paramSpace.Set(ctx, types.ParamsFixedGas, &out)
}
