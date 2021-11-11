package keeper

import (
	"context"

	"github.com/Decentr-net/decentr/x/operations/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	keeper      Keeper
	tokenKeeper types.TokenKeeper
}

// NewQueryServer returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper, tokenKeeper types.TokenKeeper) types.QueryServer {
	return &queryServer{
		keeper:      keeper,
		tokenKeeper: tokenKeeper,
	}
}

func (s queryServer) MinGasPrice(
	goCtx context.Context,
	_ *types.MinGasPriceRequest,
) (*types.MinGasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MinGasPriceResponse{
		MinGasPrice: s.keeper.GetParams(ctx).MinGasPrice,
	}, nil
}

func (s queryServer) IsAccountBanned(
	goCtx context.Context,
	r *types.IsAccountBannedRequest,
) (*types.IsAccountBannedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.IsAccountBannedResponse{
		IsBanned: s.tokenKeeper.IsBanned(ctx, r.Address),
	}, nil
}
