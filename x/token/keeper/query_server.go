package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/token/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	keeper Keeper
}

// NewQueryServer returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) types.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

func (s queryServer) Balance(goCtx context.Context, r *types.BalanceRequest) (*types.BalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.BalanceResponse{
		Balance: sdk.DecProto{Dec: s.keeper.GetBalance(ctx, r.Address)},
	}, nil
}
