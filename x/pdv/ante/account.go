package ante

import (
	"fmt"

	"github.com/Decentr-net/decentr/x/token"
	"github.com/Decentr-net/decentr/x/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type CreateAccountDecorator struct {
	tokenKeeper   token.Keeper
	accountKeeper keeper.AccountKeeper
}

func NewCreateAccountDecorator(ak auth.AccountKeeper, tokenKeeper token.Keeper) CreateAccountDecorator {
	return CreateAccountDecorator{
		tokenKeeper: tokenKeeper, accountKeeper: ak}
}

func (cad CreateAccountDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	for _, msg := range tx.GetMsgs() {
		sendTx, ok := msg.(bank.MsgSend)
		if ok {
			acc := cad.accountKeeper.GetAccount(ctx, sendTx.ToAddress)
			if acc != nil && !cad.tokenKeeper.IsInitialBalanceSet(ctx, sendTx.ToAddress) {
				cad.tokenKeeper.SetBalance(ctx, sendTx.ToAddress, utils.InitialTokenBalance())
				ctx.Logger().Info(fmt.Sprintf("account %s created", sendTx.ToAddress))
			}
		}
	}
	return next(ctx, tx, simulate)
}
