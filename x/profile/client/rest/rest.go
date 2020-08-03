package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers profile-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/private/{address}", storeName), queryPrivateHandler(cliCtx, storeName)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/public/{address}", storeName), queryPublicHandler(cliCtx, storeName)).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/%s/private/{address}", storeName), setPrivateHandler(cliCtx)).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/%s/public/{address}", storeName), setPublicHandler(cliCtx)).Methods(http.MethodPost)
}
