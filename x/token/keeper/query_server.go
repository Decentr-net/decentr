package keeper

import (
	"context"

	"github.com/Decentr-net/decentr/config"
	"github.com/Decentr-net/decentr/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	keeper             Keeper
	distributionKeeper types.DistributionKeeper
}

// NewQueryServer returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper, distributionKeeper types.DistributionKeeper) types.QueryServer {
	return &queryServer{
		keeper:             keeper,
		distributionKeeper: distributionKeeper,
	}
}

func (s queryServer) Balance(goCtx context.Context, r *types.BalanceRequest) (*types.BalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.BalanceResponse{
		Balance:      sdk.DecProto{Dec: s.keeper.GetBalance(ctx, r.Address)},
		BalanceDelta: sdk.DecProto{Dec: s.keeper.GetBalanceDelta(ctx, r.Address)},
		IsBanned:     s.keeper.IsBanned(ctx, r.Address),
	}, nil
}

func (s queryServer) Pool(goCtx context.Context, _ *types.PoolRequest) (*types.PoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	interval := s.keeper.GetParams(ctx).RewardsBlockInterval

	size := sdk.DecCoin{
		Denom:  config.DefaultBondDenom,
		Amount: s.distributionKeeper.GetFeePoolCommunityCoins(ctx).AmountOf(config.DefaultBondDenom),
	}
	totalDelta := sdk.DecProto{Dec: s.keeper.GetAccumulatedDelta(ctx)}
	nextDistributionHeight := interval * (uint64(ctx.BlockHeight())/interval + 1)

	return &types.PoolResponse{
		Size_:                  size,
		TotalDelta:             totalDelta,
		NextDistributionHeight: nextDistributionHeight,
	}, nil
}
