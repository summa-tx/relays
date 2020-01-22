package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	s := r.PathPrefix(fmt.Sprintf("/%s", storeName)).Subrouter()

	// add new tx msg routes here
	s.HandleFunc("/ingestheaderchain", ingestHeaderChainHandler(cliCtx)).Methods("POST")
	s.HandleFunc("/ingestdiffchange", ingestDifficultyChangeHandler(cliCtx)).Methods("POST")
	s.HandleFunc("/marknewheaviest", markNewHeaviestHandler(cliCtx)).Methods("POST")
	s.HandleFunc("/newrequest", newRequestHandler(cliCtx)).Methods("POST")

	// add new query routes below
	// {} denotes variable parts of the url route
	// These are our function arguments
	s.HandleFunc("/isancestor/{digest}/{ancestor}/", isAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/isancestor/{digest}/{ancestor}/{limit}", isAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/getrelaygenesis", getRelayGenesisHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/getlastreorglca", getLastReorgLCAHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/findancestor/{digest}/{offset}", findAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/ismostrecentcommonancestor/{ancestor}/{left}/{right}/", isMostRecentCommonAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/ismostrecentcommonancestor/{ancestor}/{left}/{right}/{limit}", isMostRecentCommonAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/heaviestfromancestor/{ancestor}/{currentbest}/{newbest}/", heaviestFromAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/heaviestfromancestor/{ancestor}/{currentbest}/{newbest}/{limit}", heaviestFromAncestorHandler(cliCtx, storeName)).Methods("GET")
	s.HandleFunc("/getrequest/{id}", getRequestHandler(cliCtx, storeName)).Methods("GET")
}
