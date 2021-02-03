package rest

import (
	"errors"
	"fmt"
	"net/http"

	cerberusapi "github.com/Decentr-net/cerberus/pkg/api"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/Decentr-net/decentr/x/pdv/types"
)

// RegisterRoutes registers pdv-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s", storeName), createPDVHandler(cliCtx)).Methods(http.MethodPost)

	r.HandleFunc(fmt.Sprintf("/%s/cerberus-addr", storeName), cerberusAddrHandler(cliCtx)).Methods(http.MethodGet)
}

type createPDVReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	ID      uint64       `json:"id"`
}

func createPDVHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createPDVReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		caddr, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/cerberus-addr", types.QuerierRoute), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get cerberus addr: %w", err).Error())
			return
		}

		if _, err := cerberusapi.NewClient(string(caddr), secp256k1.PrivKeySecp256k1{}).GetPDVMeta(r.Context(), req.BaseReq.From, req.ID); err != nil {
			if errors.Is(err, cerberusapi.ErrNotFound) {
				rest.WriteErrorResponse(w, http.StatusNotFound, "pdv does not exist")
				return
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to check pdv existence: %w", err).Error())
			return
		}

		msg := types.NewMsgCreatePDV(owner, req.ID)

		utils.WriteGenerateStdTxResponse(w, cliCtx.WithHeight(height), req.BaseReq, []sdk.Msg{msg})
	}
}

func cerberusAddrHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/cerberus-addr", types.QuerierRoute), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), struct {
			Address string `json:"address"`
		}{
			Address: string(res),
		})
	}
}
