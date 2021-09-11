package operations

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/community"
	"github.com/Decentr-net/decentr/x/token"
)

// NewHandler creates an sdk.Handler for all the pdv type messages
func NewHandler(keeper Keeper, tokensKeeper token.Keeper, communityKeeper community.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {

		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgDistributeRewards:
			return handleMsgDistributeRewards(ctx, keeper, tokensKeeper, msg)
		case MsgResetAccount:
			return handleMsgResetAccount(ctx, keeper, tokensKeeper, communityKeeper, msg)
		case MsgBanAccount:
			return handleMsgBanAccount(ctx, keeper, tokensKeeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgDistributeRewards(ctx sdk.Context, keeper Keeper, tokensKeeper token.Keeper, msg MsgDistributeRewards) (*sdk.Result, error) {
	owners := keeper.GetParams(ctx).Supervisors

	for _, v := range owners {
		addr, _ := sdk.AccAddressFromBech32(v)
		if msg.Owner.Equals(addr) && !addr.Empty() {
			for _, reward := range msg.Rewards {
				tokensKeeper.AddTokens(ctx, reward.Receiver, sdk.NewIntFromUint64(reward.Reward))
			}
			return &sdk.Result{}, nil
		}
	}

	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Owner is not a Cerberus owner")
}

func handleMsgResetAccount(ctx sdk.Context, keeper Keeper, tokensKeeper token.Keeper, communityKeeper community.Keeper, msg MsgResetAccount) (*sdk.Result, error) {
	reset := func(resetBy sdk.AccAddress) {
		tokensKeeper.ResetAccount(ctx, msg.Owner)
		communityKeeper.ResetAccount(ctx, msg.Owner)
		ctx.Logger().Info(fmt.Sprintf("account %s reset by %s", msg.Owner, resetBy))
	}

	for _, v := range keeper.GetParams(ctx).Supervisors {
		addr, _ := sdk.AccAddressFromBech32(v)
		if msg.Owner.Equals(addr) && !addr.Empty() {
			reset(addr)
			return &sdk.Result{}, nil
		}
	}

	if !msg.Owner.Equals(msg.AccountOwner) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized,
			fmt.Sprintf("%s can not delete %s", msg.Owner, msg.AccountOwner))
	}

	reset(msg.AccountOwner)
	return &sdk.Result{}, nil
}

func handleMsgBanAccount(ctx sdk.Context, keeper Keeper, tokenKeeper token.Keeper, msg MsgBanAccount) (*sdk.Result, error) {
	found := false
	for _, v := range keeper.GetParams(ctx).Supervisors {
		addr, _ := sdk.AccAddressFromBech32(v)
		if msg.Owner.Equals(addr) && !addr.Empty() {
			found = true
			break
		}
	}

	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized,
			fmt.Sprintf("%s can not ban %s", msg.Owner, msg.Address))
	}

	tokenKeeper.SetBan(ctx, msg.Address, msg.Ban)

	if msg.Ban {
		ctx.Logger().Info(fmt.Sprintf("account %s banned by %s", msg.Address, msg.Owner))
	} else {
		ctx.Logger().Info(fmt.Sprintf("account %s unbanned by %s", msg.Address, msg.Owner))
	}

	return &sdk.Result{}, nil
}
