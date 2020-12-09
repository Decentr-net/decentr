package rest

import (
	"fmt"
	"net/http"

	"github.com/Decentr-net/decentr/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers token-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/balance/{address}", storeName), queryBalanceHandler(cliCtx, storeName)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/stats/{address}", storeName), queryStatsHandler(cliCtx)).Methods(http.MethodGet)
}

func queryBalanceHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bech32Addr := mux.Vars(r)["address"]

		owner, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/balance/%s", storeName, owner), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func queryStatsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bech32Addr := mux.Vars(r)["address"]

		owner, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/stats/%s", types.QuerierRoute, owner), nil)
		if err != nil {
			if err, ok := err.(*sdkerrors.Error); ok {
				if err.Is(sdkerrors.ErrInvalidRequest) {
					rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
					return
				}
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}
