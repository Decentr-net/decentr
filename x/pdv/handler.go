package pdv

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/token"
	"github.com/Decentr-net/decentr/x/utils"
)

// NewHandler creates an sdk.Handler for all the pdv type messages
func NewHandler(keeper Keeper, tokensKeeper token.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {

		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgDistributeRewards:
			return handleMsgDistributeRewards(ctx, keeper, tokensKeeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgDistributeRewards(ctx sdk.Context, keeper Keeper, tokensKeeper token.Keeper, msg MsgDistributeRewards) (*sdk.Result, error) {
	owners := keeper.GetCerberusOwners(ctx)

	for _, v := range owners {
		addr, _ := sdk.AccAddressFromBech32(v)
		if msg.Owner.Equals(addr) && !addr.Empty() {
			for _, reward := range msg.Rewards {
				tokensKeeper.AddTokens(ctx, reward.Receiver, sdk.NewIntFromUint64(reward.Reward), utils.GetHash(msg))
			}
			return &sdk.Result{}, nil
		}
	}

	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Owner is not a Cerberus owner")
}
