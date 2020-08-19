package pdv

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	cerberusapi "github.com/Decentr-net/cerberus/pkg/api"
)

// NewHandler creates an sdk.Handler for all the pdv type messages
func NewHandler(cerberus cerberusapi.Cerberus, keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreatePDV:
			return handleMsgCreatePDV(ctx, cerberus, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreatePDV(ctx sdk.Context, cerberus cerberusapi.Cerberus, keeper Keeper, msg MsgCreatePDV) (*sdk.Result, error) {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Address)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	exist, err := cerberus.DoesPDVExist(ctx.Context(), msg.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPanic, fmt.Sprintf("cerberus call failed: %s", err.Error()))
	}

	if !exist {
		return nil, errors.New("pdv does not exist")
	}

	keeper.SetPDV(ctx, msg.Address, PDV{Owner: msg.Owner, Address: msg.Address, Type: msg.DataType})

	// ToDo: set tokens count

	return &sdk.Result{}, nil
}
