package keeper

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/pdv/types"
)

// query endpoints supported by the pdv Querier
const (
	QueryOwner        = "owner"
	QueryShow         = "show"
	QueryList         = "list"
	QueryStats        = "stats"
	QueryCerberusAddr = "cerberus-addr"
)

const isoDateFormat = "2006-01-02"

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
		case QueryStats:
			return queryStats(ctx, path[1:], req, keeper)
		case QueryCerberusAddr:
			return queryCerberusAddr(keeper)
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
		from  *time.Time
		limit = uint(20)
	)

	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if path[1] != "" {
		v, err := time.Parse(time.RFC3339, path[1])
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid page")
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

	p, err := keeper.stats.ListPDV(owner, from, limit)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPanic, err.Error())
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, p)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
// queryStats returns map[time.Time]sdk.Int. The statistics is daily, every key is truncated by 24 hours.
func queryStats(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	s, err := keeper.stats.GetStats(owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPanic, err.Error())
	}

	hs := make(map[string]sdk.Int, len(s))
	for k, v := range s {
		hs[k.Format(isoDateFormat)] = v
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, hs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
func queryCerberusAddr(keeper Keeper) ([]byte, error) {
	return []byte(viper.GetString(types.FlagCerberusAddr)), nil
}
