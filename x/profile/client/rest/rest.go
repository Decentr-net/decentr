package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers profile-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/private/{address}", storeName), queryPrivate(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/public/{address}", storeName), queryPublic(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/private/{address}", storeName), setPrivateHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/public/{address}", storeName), setPublicHandler(cliCtx)).Methods("POST")
}
