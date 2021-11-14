package operations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/operations/keeper"
	"github.com/Decentr-net/decentr/x/operations/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if k.SetParams(ctx, types.DefaultParams()); genState.Params != nil {
		k.SetParams(ctx, *genState.Params)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	return &types.GenesisState{Params: &params}
}
