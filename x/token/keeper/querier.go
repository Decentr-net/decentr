package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Decentr-net/decentr/x/utils"
)

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
func queryBalance(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper Keeper) ([]byte, error) {
	owner, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	var balance sdk.Int
	if !keeper.IsInitialBalanceSet(ctx, owner) {
		// account not found, return default value
		balance = keeper.InitialTokenBalance()
	} else {
		balance = keeper.GetBalance(ctx, owner)
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, struct {
		Balance float64 `json:"balance" amino:"unsafe"`
	}{
		Balance: utils.TokenToFloat64(balance),
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
