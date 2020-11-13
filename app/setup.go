package app

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/Decentr-net/decentr/x/community"
	"github.com/Decentr-net/decentr/x/pdv"
)

func NewAnteHandler(ak auth.AccountKeeper, sk supply.Keeper) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewGasExcludingSetUpContextDecorator(pdv.MsgCreatePDV{}, community.MsgCreatePost{},
			community.MsgDeletePost{}, community.MsgSetLike{}),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, sk),
		ante.NewSigGasConsumeDecorator(ak, auth.DefaultSigVerificationGasConsumer),
		ante.NewSigVerificationDecorator(ak),
		ante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
	)
}

type GasExcludingSetUpContextDecorator struct {
	Exclude []reflect.Type
}

func NewGasExcludingSetUpContextDecorator(exclude ...interface{}) GasExcludingSetUpContextDecorator {
	e := make([]reflect.Type, len(exclude))

	for i, v := range exclude {
		e[i] = reflect.TypeOf(v)
	}

	return GasExcludingSetUpContextDecorator{
		Exclude: e,
	}
}

func (sud GasExcludingSetUpContextDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	freeFromGas := true

check:
	for _, v := range tx.GetMsgs() {
		for _, e := range sud.Exclude {
			if reflect.TypeOf(v).String() != e.String() {
				freeFromGas = false
				break check
			}
		}
	}

	if !freeFromGas {
		return ante.NewSetUpContextDecorator().AnteHandle(ctx, tx, simulate, next)
	}

	return next(ctx.WithGasMeter(freeGasMeter{}).WithMinGasPrices(nil), tx, simulate)
}

type freeGasMeter struct {
}

func (g freeGasMeter) GasConsumed() sdk.Gas {
	return 0
}

func (g freeGasMeter) GasConsumedToLimit() sdk.Gas {
	return 0
}

func (g freeGasMeter) Limit() sdk.Gas {
	return 0
}

func (g freeGasMeter) ConsumeGas(_ sdk.Gas, _ string) {
}

func (g freeGasMeter) IsPastLimit() bool {
	return false
}

func (g freeGasMeter) IsOutOfGas() bool {
	return false
}
