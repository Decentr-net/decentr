package ante

import (
	"github.com/Decentr-net/decentr/x/community"
	"github.com/Decentr-net/decentr/x/operations"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FixedGasTxDecorator struct {
	config map[string]func(ctx sdk.Context) sdk.Gas
}

func NewFixedGasTxDecorator(pk operations.Keeper, ck community.Keeper) FixedGasTxDecorator {
	config := map[string]func(ctx sdk.Context) sdk.Gas{
		operations.MsgResetAccount{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return pk.GetParams(ctx).FixedGas.ResetAccount
		},
		operations.MsgBanAccount{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return pk.GetParams(ctx).FixedGas.BanAccount
		},
		operations.MsgDistributeRewards{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return pk.GetParams(ctx).FixedGas.DistributeRewards
		},
		community.MsgCreatePost{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetParams(ctx).FixedGas.CreatePost
		},
		community.MsgDeletePost{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetParams(ctx).FixedGas.DeletePost
		},
		community.MsgSetLike{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetParams(ctx).FixedGas.SetLike
		},
		community.MsgFollow{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetParams(ctx).FixedGas.Follow
		},
		community.MsgUnfollow{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetParams(ctx).FixedGas.Unfollow
		},
	}

	return FixedGasTxDecorator{
		config: config,
	}
}

func (fgm FixedGasTxDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		if fixedGas, ok := fgm.config[msg.Type()]; ok {
			limit := ctx.GasMeter().Limit()

			// pass infinite gas meter since fixedGas requires gas to read parameters from keeper
			// this gas should be skipped
			consumed := fixedGas(ctx.WithGasMeter(sdk.NewInfiniteGasMeter()))

			// prepare new context
			ctx := ctx.WithGasMeter(NewFixedGasMeter(consumed, limit))

			if consumed == 0 {
				// special case: consumed gas is zero, what could be for DistributeRewards trx
				// set min gas price to zero, otherwise sdkerrors.ErrInsufficientFee occurs
				zeroDecCoins := sdk.NewDecCoins()
				return next(ctx.WithMinGasPrices(zeroDecCoins), tx, simulate)
			}

			return next(ctx, tx, simulate)
		}
	}

	return next(ctx, tx, simulate)
}

type fixedGasMeter struct {
	limit    sdk.Gas
	consumed sdk.Gas
}

// NewFixedGasMeter returns a reference to a new basicGasMeter.
func NewFixedGasMeter(consumed, limit sdk.Gas) sdk.GasMeter {
	return &fixedGasMeter{
		limit:    limit,
		consumed: consumed,
	}
}

func (g *fixedGasMeter) GasConsumed() sdk.Gas {
	return g.consumed
}

func (g *fixedGasMeter) Limit() sdk.Gas {
	return g.limit
}

func (g *fixedGasMeter) GasConsumedToLimit() sdk.Gas {
	if g.IsPastLimit() {
		return g.limit
	}
	return g.consumed
}

func (g *fixedGasMeter) ConsumeGas(_ sdk.Gas, _ string) {
}

func (g *fixedGasMeter) IsPastLimit() bool {
	return g.consumed > g.limit
}

func (g *fixedGasMeter) IsOutOfGas() bool {
	return g.consumed >= g.limit
}
