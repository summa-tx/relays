package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	restName = "relay"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/link", storeName), setLinkHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/parent/{%s}", storeName, restName), getParentHandler(cliCtx, storeName)).Methods("GET")
}
