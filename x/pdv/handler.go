package pdv

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	cerberusapi "github.com/Decentr-net/cerberus/pkg/api"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// NewHandler creates an sdk.Handler for all the pdv type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		cerberus := getCerberus(ctx, keeper)

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

var (
	prevCerberusAddress string
	cerberus            cerberusapi.Cerberus
)

func getCerberus(ctx sdk.Context, keeper Keeper) cerberusapi.Cerberus {
	addr := keeper.GetCerberusAddr(ctx)
	if prevCerberusAddress != addr {
		// address changed, create a new Cerberus
		cerberus = cerberusapi.NewClient(addr, secp256k1.PrivKeySecp256k1{})
		prevCerberusAddress = addr
	}
	return cerberus
}

func handleMsgCreatePDV(ctx sdk.Context, cerberus cerberusapi.Cerberus, keeper Keeper, msg MsgCreatePDV) (*sdk.Result, error) {
	if hex.EncodeToString(msg.Owner) != strings.Split(msg.Address, "-")[0] { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}

	exist, err := cerberus.DoesPDVExist(ctx.Context(), msg.Address)
	if err != nil {
		return nil, fmt.Errorf("cerberus call failed: %w", err)
	}

	if !exist {
		return nil, errors.New("pdv does not exist")
	}

	if keeper.IsHashPresent(ctx, msg.Address) {
		return &sdk.Result{}, nil
	}

	keeper.SetPDV(ctx, msg.Address, PDV{Timestamp: msg.Timestamp, Owner: msg.Owner, Address: msg.Address, Type: msg.DataType})

	return &sdk.Result{}, nil
}
