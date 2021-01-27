package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// query endpoints supported by the pdv Querier
const (
	QueryOwner        = "owner"
	QueryShow         = "show"
	QueryList         = "list"
	QueryCerberusAddr = "cerberus-addr"
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
		case QueryCerberusAddr:
			return queryCerberusAddr(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown pdv query endpoint")
		}
	}
}

// nolint: unparam
func queryOwner(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	pdv := keeper.GetPDV(ctx, path[0])
	return []byte(pdv.Owner.String()), nil
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
	var (
		from  *uint64
		limit = uint(20)
	)

	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if path[1] != "" {
		v, err := strconv.ParseUint(path[1], 10, 64)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid from")
		}
		from = &v
	}

	if path[2] != "" {
		v, err := strconv.Atoi(path[2])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid page")
		}
		limit = uint(v)
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, keeper.ListPDV(ctx, owner, from, limit))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryCerberusAddr(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	return []byte(keeper.GetCerberusAddr(ctx)), nil
}
