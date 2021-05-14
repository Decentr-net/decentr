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
			return pk.GetFixedGasParams(ctx).ResetAccount
		},
		operations.MsgDistributeRewards{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return pk.GetFixedGasParams(ctx).DistributeRewards
		},
		community.MsgCreatePost{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetFixedGasParams(ctx).CreatePost
		},
		community.MsgDeletePost{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetFixedGasParams(ctx).DeletePost
		},
		community.MsgSetLike{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetFixedGasParams(ctx).SetLike
		},
		community.MsgFollow{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetFixedGasParams(ctx).Follow
		},
		community.MsgUnfollow{}.Type(): func(ctx sdk.Context) sdk.Gas {
			return ck.GetFixedGasParams(ctx).Unfollow
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
			consumed := fixedGas(ctx)
			return ctx.WithGasMeter(NewFixedGasMeter(consumed, limit)), nil
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
