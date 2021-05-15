package ante

import (
	"github.com/Decentr-net/decentr/x/operations/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MinGasPriceDecorator struct {
	keeper keeper.Keeper
}

func NewMinGasPriceDecorator(keeper keeper.Keeper) *MinGasPriceDecorator {
	return &MinGasPriceDecorator{keeper: keeper}
}

func (mgp MinGasPriceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	ctx = ctx.WithMinGasPrices(sdk.NewDecCoins(mgp.keeper.GetMinGasPrice(ctx)))
	return next(ctx, tx, simulate)
}
