package decentr

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the decentr type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreatePDV:
			return handleMsgCreatePDV(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreatePDV(ctx sdk.Context, keeper Keeper, msg MsgCreatePDV) (*sdk.Result, error) {
	// TODO: call Cerberus
	hash := ""

	if !msg.Owner.Equals(keeper.GetOwner(ctx, hash)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	keeper.SetPDV(ctx, hash, PDV{Owner: msg.Owner, Hash: hash})
	return &sdk.Result{}, nil
}
