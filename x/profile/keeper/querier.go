package keeper

import (
	"encoding/base64"

	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// query endpoints supported by the profile Querier
const (
	QueryPublic  = "public"
	QueryPrivate = "private"
)

// NewQuerier creates a new querier for profile clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryPublic:
			return queryPublic(ctx, path[1:], req, keeper)
		case QueryPrivate:
			return queryPrivate(ctx, path[1:], req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown profile query endpoint")
		}
	}
}

// nolint: unparam
func queryPublic(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	profile := keeper.GetProfile(ctx, owner)

	res, err := codec.MarshalJSONIndent(keeper.cdc, profile.Public)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
func queryPrivate(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	profile := keeper.GetProfile(ctx, owner)
	return []byte(base64.StdEncoding.EncodeToString(profile.Private)), nil
}
