package keeper

import (
	"github.com/Decentr-net/decentr/x/token/types"
	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// query endpoints supported by the token Querier
const (
	QueryBalance = "balance"
)

// NewQuerier creates a new querier for token clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, path[1:], req, keeper)
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
	out := TokenToFloat64(balance)

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

// TokenToFloat64 converts token to its float64 representation
func TokenToFloat64(token sdk.Int) float64 {
	return float64(token.Int64()) / float64(types.Denominator)
}
