package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/utils"
)

// query endpoints supported by the token Querier
const (
	QueryBalance = "balance"
	QueryStats   = "stats"
)

const isoDateFormat = "2006-01-02"

// DateValue is date-value stat item
type DateValue struct {
	Date  string  `json:"date"`
	Value float64 `json:"value" amino:"unsafe"`
}

// NewQuerier creates a new querier for token clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, path[1:], req, keeper)
		case QueryStats:
			return queryStats(ctx, path[1:], req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown token query endpoint")
		}
	}
}

// nolint: unparam
func queryBalance(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	balance := keeper.GetBalance(ctx, owner)
	out := utils.TokenToFloat64(balance)

	res, err := codec.MarshalJSONIndent(keeper.cdc, struct {
		Balance float64 `json:"balance" amino:"unsafe"`
	}{
		Balance: out,
	})
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

	s := keeper.GetStats(ctx, owner)

	i := 0
	stats := make([]DateValue, len(s))
	for k, v := range s {
		stats[i] = DateValue{Date: time.Unix(int64(k), 0).Format(isoDateFormat), Value: v}
		i++
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, stats)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
