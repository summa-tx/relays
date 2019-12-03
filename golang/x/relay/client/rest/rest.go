package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	s := r.PathPrefix(fmt.Sprintf("/%s", storeName)).Subrouter()

	// r.HandleFunc(fmt.Sprintf("/%s/link", storeName), setLinkHandler(cliCtx)).Methods("POST")

	// @Erin add new query routes below
	//     {} denotes variable parts of the url route
	//     These are our function arguments
	s.HandleFunc("/isancestor/{digest}/{ancestor}/{limit}", isAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/findancestor/{digest}/{offset}", findAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/ismostrecentcommonancestor/{ancestor}/{left}/{right}/", isMostRecentCommonAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/ismostrecentcommonancestor/{ancestor}/{left}/{right}/{limit}", isMostRecentCommonAncestorHandler(cliCtx, storeName)).Methods("GET")
}
