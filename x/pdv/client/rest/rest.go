package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	cerberusapi "github.com/Decentr-net/cerberus/pkg/api"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/tendermint/tendermint/crypto/secp256k1"

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

	r.HandleFunc(fmt.Sprintf("/%s/cerberus-addr", storeName), cerberusAddrHandler(cliCtx)).Methods(http.MethodGet)
}

type createPDVReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Address string       `json:"address"`
}

func queryOwnerHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramAddress := mux.Vars(r)["address"]

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/owner/%s", types.QuerierRoute, paramAddress), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), string(res))
	}
}

func queryShowHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramAddress := mux.Vars(r)["address"]

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/show/%s", types.QuerierRoute, paramAddress), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func queryListHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramOwner := mux.Vars(r)["owner"]
		q := r.URL.Query()

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/list/%s/%s/%s", types.QuerierRoute, paramOwner, q.Get("from"), q.Get("limit")), nil)
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

		if hex.EncodeToString(owner) != strings.Split(req.Address, "-")[0] { // Checks if the the msg sender is the same as the current owner
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid owner")
			return
		}

		caddr, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/cerberus-addr", types.QuerierRoute), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get cerberus addr: %w", err).Error())
			return
		}

		exists, err := cerberusapi.NewClient(string(caddr), secp256k1.PrivKeySecp256k1{}).DoesPDVExist(r.Context(), req.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to check pdv existence: %w", err).Error())
			return
		}

		if !exists {
			rest.WriteErrorResponse(w, http.StatusNotFound, "pdv does not exist")
			return
		}

		msg := types.NewMsgCreatePDV(uint64(time.Now().Unix()), req.Address, types.PDVTypeCookie, owner)

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
