package rest

import (
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/Decentr-net/decentr/x/pdv/types"
)

// RegisterRoutes registers pdv-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/{address}/owner", storeName), queryOwnerHandler(cliCtx)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/{address}/show", storeName), queryShowHandler(cliCtx)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/{owner}/list", storeName), queryListHandler(cliCtx)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s", storeName), createPDVHandler(cliCtx)).Methods(http.MethodPost)
}

type createPDVReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Address string       `json:"address"`
}

func queryOwnerHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramType := mux.Vars(r)["address"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/owner/%s", types.QuerierRoute, paramType), nil)
		if err != nil {
			if err, ok := err.(*sdkerrors.Error); ok {
				if err.Is(types.ErrNotFound) {
					rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
					return
				}
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryShowHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramType := mux.Vars(r)["address"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/show/%s", types.QuerierRoute, paramType), nil)
		if err != nil {
			if err, ok := err.(*sdkerrors.Error); ok {
				if err.Is(types.ErrNotFound) {
					rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
					return
				}
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryListHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramType := mux.Vars(r)["owner"]
		q := r.URL.Query()

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/list/%s/%s/%s", types.QuerierRoute, paramType, q.Get("page"), q.Get("limit")), nil)
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

		rest.PostProcessResponse(w, cliCtx, res)
	}
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

		msg := types.NewMsgCreatePDV(req.Address, types.PDVTypeCookie, owner)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
