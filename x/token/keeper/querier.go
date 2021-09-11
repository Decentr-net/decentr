package keeper

import (
	"github.com/Decentr-net/decentr/x/token/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	QueryBalance = "balance"
	QueryPool    = "pool"
	QueryHistory = "history"
)

type Balance struct {
	Balance      sdk.Dec `json:"balance"`
	BalanceDelta sdk.Dec `json:"balanceDelta"`
	IsBanned     bool    `json:"isBanned,omitempty"`
}

type Pool struct {
	Size                   sdk.DecCoins `json:"size"`
	TotalDelta             sdk.Dec      `json:"totalDelta"`
	NextDistributionHeight int64        `json:"nextDistributionHeight"`
}

// NewQuerier creates a new querier for token clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, path[1:], req, keeper)
		case QueryPool:
			return queryPool(ctx, path[1:], req, keeper)
		case QueryHistory:
			return queryHistory(ctx, path[1:], req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown token query endpoint")
		}
	}
}

// nolint: unparam
func queryBalance(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	bz, err := keeper.cdc.MarshalBinaryBare(Balance{
		Balance:      keeper.GetBalance(ctx, owner).ToDec().QuoInt(types.Denominator),
		BalanceDelta: keeper.GetBalanceDelta(ctx, owner).ToDec().QuoInt(types.Denominator),
		IsBanned:     keeper.IsBanned(ctx, owner),
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// nolint: unparam
func queryPool(ctx sdk.Context, _ []string, _ abci.RequestQuery, keeper Keeper) ([]byte, error) {
	interval := keeper.GetParams(ctx).RewardsBlockInterval

	bz, err := keeper.cdc.MarshalBinaryBare(Pool{
		Size:                   keeper.distributionKeeper.GetFeePoolCommunityCoins(ctx),
		TotalDelta:             keeper.GetBalanceDelta(ctx, types.AccumulatedDelta).ToDec().QuoInt(types.Denominator),
		NextDistributionHeight: interval * (ctx.BlockHeight()/interval + 1),
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// nolint: unparam
func queryHistory(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	bz, err := keeper.cdc.MarshalBinaryBare(keeper.GetRewardsDistributionHistory(ctx, owner))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
