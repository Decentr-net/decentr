package app

import (
	appante "github.com/Decentr-net/decentr/app/ante"
	"github.com/Decentr-net/decentr/x/community"
	"github.com/Decentr-net/decentr/x/operations"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

func NewAnteHandler(ak auth.AccountKeeper, sk supply.Keeper,
	pk operations.Keeper, ck community.Keeper) sdk.AnteHandler {

	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(),
		appante.NewFixedGasTxDecorator(pk, ck),
		operations.NewMinGasPriceDecorator(pk),
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
