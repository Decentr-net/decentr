package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// query endpoints supported by the pdv Querier
const (
	QueryCerberusAddr = "cerberus-addr"
)

// NewQuerier creates a new querier for pdv clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryCerberusAddr:
			return queryCerberusAddr(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown pdv query endpoint")
		}
	}
}

func queryCerberusAddr(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	return []byte(keeper.GetCerberusAddr(ctx)), nil
}
