package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// query endpoints supported by the pdv Querier
const (
	QueryOwner = "owner"
	QueryShow  = "show"
)

// NewQuerier creates a new querier for pdv clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryOwner:
			return queryOwner(ctx, path[1:], req, keeper)
		case QueryShow:
			return queryShow(ctx, path[1:], req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown pdv query endpoint")
		}
	}
}

// nolint: unparam
func queryOwner(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	pdv := keeper.GetPDV(ctx, path[0])
	return pdv.Owner, nil
}

// nolint: unparam
func queryShow(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	pdv := keeper.GetPDV(ctx, path[0])

	//TODO: get list from keeper

	res, err := codec.MarshalJSONIndent(keeper.cdc, pdv)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
