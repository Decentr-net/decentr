package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/pdv/types"
)

// query endpoints supported by the pdv Querier
const (
	QueryOwner = "owner"
	QueryShow  = "show"
	QueryList  = "list"
)

// NewQuerier creates a new querier for pdv clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryOwner:
			return queryOwner(ctx, path[1:], req, keeper)
		case QueryShow:
			return queryShow(ctx, path[1:], req, keeper)
		case QueryList:
			return queryList(ctx, path[1:], req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown pdv query endpoint")
		}
	}
}

// nolint: unparam
func queryOwner(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	pdv := keeper.GetPDV(ctx, path[0])
	if pdv.Owner.Empty() {
		return nil, types.ErrNotFound
	}

	return pdv.Owner, nil
}

// nolint: unparam
func queryShow(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	pdv := keeper.GetPDV(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, pdv)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
func queryList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	page := uint(0)
	limit := uint(20)

	if path[1] != "" {
		v, err := strconv.Atoi(path[1])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid page")
		}
		page = uint(v)
	}

	if path[2] == "" {
		v, err := strconv.Atoi(path[2])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid limit")
		}
		limit = uint(v)
	}

	i := keeper.GetPDVsIteratorPaginated(ctx, path[0], page, limit)

	m := make(map[string]types.PDV)
	for ; i.Valid(); i.Next() {
		var pdv types.PDV

		k, v := i.Key(), i.Value()
		keeper.cdc.MustUnmarshalBinaryBare(v, &pdv)

		m[string(k)] = pdv
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, m)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
