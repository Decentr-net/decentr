package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/Decentr-net/decentr/x/token/keeper"
	"github.com/Decentr-net/decentr/x/token/types"
)

// RegisterRoutes registers token-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/balance/{address}", storeName), queryBalanceHandler(cliCtx, storeName)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/balance/{address}/history", storeName), queryHistoryHandler(cliCtx, storeName)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/pool", storeName), queryPoolHandler(cliCtx, storeName)).Methods(http.MethodGet)
}

func queryBalanceHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bech32Addr := mux.Vars(r)["address"]

		owner, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		bz, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/balance/%s", storeName, owner), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var out keeper.Balance
		cliCtx.Codec.MustUnmarshalBinaryBare(bz, &out)

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), out)
	}
}

func queryPoolHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bz, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pool", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var out keeper.Pool
		cliCtx.Codec.MustUnmarshalBinaryBare(bz, &out)

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), out)
	}
}

func queryHistoryHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bech32Addr := mux.Vars(r)["address"]

		owner, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		bz, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/history/%s", storeName, owner), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var out []types.RewardDistribution
		cliCtx.Codec.MustUnmarshalBinaryBare(bz, &out)

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), out)
	}
}
