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
	r.HandleFunc(fmt.Sprintf("/%s/post/{postOwner}/{postUUID}", storeName), getPostHandler(cliCtx)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/posts/{owner}", storeName), queryListUserPostsHandler(cliCtx)).Methods(http.MethodGet)

	r.HandleFunc(fmt.Sprintf("/%s/posts", storeName), createPostHandler(cliCtx)).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/posts/{postOwner}/{postUUID}/like", storeName), likePostHandler(cliCtx)).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/posts/{postOwner}/{postUUID}/delete", storeName), deletePostHandler(cliCtx)).Methods(http.MethodPost)

	r.HandleFunc(fmt.Sprintf("/%s/moderators", storeName), moderatorsHandler(cliCtx)).Methods(http.MethodGet)

	r.HandleFunc(fmt.Sprintf("/%s/followers/follow/{whom}", storeName), createFollowHandler(cliCtx)).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/followers/unfollow/{whom}", storeName), createUnfollowHandler(cliCtx)).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/followers/{owner}/followees", storeName), queryFolloweesHandler(cliCtx)).Methods(http.MethodGet)
}

func getPostHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramOwner := mux.Vars(r)["postOwner"]
		paramUUID := mux.Vars(r)["postUUID"]

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/post/%s/%s", types.QuerierRoute, paramOwner, paramUUID), nil)
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

func queryFolloweesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paramOwner := mux.Vars(r)["owner"]

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/followees/%s", types.QuerierRoute, paramOwner), nil)
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

func moderatorsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/moderators", types.QuerierRoute), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}
