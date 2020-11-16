package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/Decentr-net/decentr/x/community/types"
)

// RegisterRoutes registers community-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s", storeName), createPostHandler(cliCtx)).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/owner/{postOwner}/{postUUID}/like", storeName), likePostHandler(cliCtx)).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/{owner}/{uuid}/delete", storeName), deletePostHandler(cliCtx)).Methods(http.MethodPost)

	r.HandleFunc(fmt.Sprintf("/%s", storeName), queryListPostsHandler(cliCtx)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/popular", storeName), queryListPopularPostsHandler(cliCtx)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/{owner}/posts", storeName), queryListUserPostsHandler(cliCtx)).Methods(http.MethodGet)
}

func queryListUserPostsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramOwner := mux.Vars(r)["owner"]
		q := r.URL.Query()

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/user/%s/%s/%s", types.QuerierRoute, paramOwner, q.Get("from"), q.Get("limit")), nil)
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

func queryListPopularPostsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/popular/%s/%s/%s/%s", types.QuerierRoute, q.Get("fromOwner"), q.Get("fromUUID"), q.Get("category"), q.Get("limit")), nil)
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

func queryListPostsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/posts/%s/%s/%s/%s", types.QuerierRoute, q.Get("fromOwner"), q.Get("fromUUID"), q.Get("category"), q.Get("limit")), nil)
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
