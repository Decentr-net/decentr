package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Decentr-net/decentr/x/token"
)

const (
	QueryMinGasPrice     = "min-gas-price"
	QueryIsAccountBanned = "is-banned"
)

// NewQuerier creates a new querier for token clients.
func NewQuerier(keeper Keeper, tokenKeeper token.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryMinGasPrice:
			return queryMinGasPrice(ctx, path[1:], req, keeper)
		case QueryIsAccountBanned:
			return queryIsAccountBanned(ctx, path[1:], req, keeper, tokenKeeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown pdv query endpoint")
		}
	}
}

// nolint: unparam
func queryMinGasPrice(ctx sdk.Context, _ []string, _ abci.RequestQuery, keeper Keeper) ([]byte, error) {
	mgp := keeper.GetParams(ctx).MinGasPrice

	bz, err := keeper.cdc.MarshalBinaryBare(mgp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPanic, err.Error())
	}

	return bz, nil
}

// nolint: unparam
func queryIsAccountBanned(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper Keeper, tokenKeeper token.Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	bz, err := keeper.cdc.MarshalBinaryBare(tokenKeeper.IsBanned(ctx, address))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
