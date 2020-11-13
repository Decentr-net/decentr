package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers community-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s", storeName), createPostHandler(cliCtx)).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/owner/{postOwner}/{postUUID}/like", storeName), likePostHandler(cliCtx)).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/{owner}/{uuid}/delete", storeName), deletePostHandler(cliCtx)).Methods(http.MethodPost)
}
