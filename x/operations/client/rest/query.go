package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/Decentr-net/decentr/x/operations/keeper"
	"github.com/Decentr-net/decentr/x/operations/types"
)

func queryMinGasPriceHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryMinGasPrice)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var mgp sdk.DecCoin
		cliCtx.Codec.MustUnmarshalBinaryBare(res, &mgp)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, mgp)
	}
}

func queryIsAccountBannedHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["account"]

		accountOwner, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, keeper.QueryIsAccountBanned, accountOwner)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bz, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var ban bool
		cliCtx.Codec.MustUnmarshalBinaryBare(bz, &ban)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, ban)
	}
}
