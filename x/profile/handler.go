package profile

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the profile type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgSetPrivate:
			return handleMsgSetPrivate(ctx, keeper, msg)
		case MsgSetPublic:
			return handleMsgSetPublic(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgSetPrivate(ctx sdk.Context, keeper Keeper, msg MsgSetPrivate) (*sdk.Result, error) {
	profile := keeper.GetProfile(ctx, msg.Owner)
	if !msg.Owner.Equals(profile.Owner) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	profile.Private = msg.Private
	keeper.SetProfile(ctx, msg.Owner, profile)
	return &sdk.Result{}, nil
}

func handleMsgSetPublic(ctx sdk.Context, keeper Keeper, msg MsgSetPublic) (*sdk.Result, error) {
	profile := keeper.GetProfile(ctx, msg.Owner)
	if !msg.Owner.Equals(profile.Owner) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	profile.Public = msg.Public
	keeper.SetProfile(ctx, msg.Owner, profile)
	return &sdk.Result{}, nil
}
