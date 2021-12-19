package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/operations/types"
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
