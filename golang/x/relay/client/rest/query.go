package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/summa-tx/relays/golang/x/relay/types"
)

// handler function for isAncestor queries. parses arguments from url string, and passes them through
// as a QueryParamsIsAncestor struct
func isAncestorHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// mux.Vars holds the variable elements of the URL from rest.go
		vars := mux.Vars(r)

		digestLE, sdkErr := types.Hash256DigestFromHex(vars["digest"])
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}
		ancestor, sdkErr := types.Hash256DigestFromHex(vars["ancestor"])
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}

		var limit uint32
		if val, ok := vars["limit"]; ok {
			lim, err := strconv.ParseUint(val, 0, 32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
				return
			}
			limit = uint32(lim)
		}

		params := types.QueryParamsIsAncestor{
			DigestLE:            digestLE,
			ProspectiveAncestor: ancestor,
			Limit:               limit,
		}

		queryData, err := json.Marshal(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// run the query. the routeString is passed as strings to our querier switch/case in `keeper/querier.go`
		res, _, err := cliCtx.QueryWithData("custom/relay/isancestor", queryData)

		// below this is boilerplate
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// handler function for getRelayGenesis queries
func getRelayGenesisHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData("custom/relay/getrelaygenesis", nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// handler function for getLastReorgLCA queries
func getLastReorgLCAHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData("custom/relay/getlastreorglca", nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// handler function for findAncestor queries. parses arguments from url string, and passes them through
// as a QueryParamsFindAncestor struct
func findAncestorHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// mux.Vars holds the variable elements of the URL from rest.go
		vars := mux.Vars(r)

		digestLE, sdkErr := types.Hash256DigestFromHex(vars["digest"])
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}

		off, err := strconv.ParseUint(vars["offset"], 0, 32)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}
		offset := uint32(off)

		params := types.QueryParamsFindAncestor{
			DigestLE: digestLE,
			Offset:   offset,
		}

		queryData, err := json.Marshal(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// run the query. the routeString is passed as strings to our querier switch/case in `keeper/querier.go`
		res, _, err := cliCtx.QueryWithData("custom/relay/findancestor", queryData)

		// below this is boilerplate
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// handler function for IsMostRecentCommonAncestor queries. parses arguments from url string, and passes them through
// as a QueryParamsIsMostRecentCommonAncestor struct
func isMostRecentCommonAncestorHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// mux.Vars holds the variable elements of the URL from rest.go
		vars := mux.Vars(r)

		ancestor, sdkErr := types.Hash256DigestFromHex(vars["ancestor"])
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}

		left, sdkErr := types.Hash256DigestFromHex(vars["left"])
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}

		right, sdkErr := types.Hash256DigestFromHex(vars["right"])
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}

		var limit uint32
		if val, ok := vars["limit"]; ok {
			lim, err := strconv.ParseUint(val, 0, 32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
				return
			}
			limit = uint32(lim)
		}

		params := types.QueryParamsIsMostRecentCommonAncestor{
			Ancestor: ancestor,
			Left:     left,
			Right:    right,
			Limit:    limit,
		}

		queryData, err := json.Marshal(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// run the query. the routeString is passed as strings to our querier switch/case in `keeper/querier.go`
		res, _, err := cliCtx.QueryWithData("custom/relay/ismostrecentcommonancestor", queryData)

		// below this is boilerplate
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// handler function for heaviestFromAncestor queries. parses arguments from url string, and passes them
// through as a QueryParamsHeaviestFromAncestor struct
func heaviestFromAncestorHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// mux.Vars holds the variable elements of the URL from rest.go
		vars := mux.Vars(r)

		ancestor, sdkErr := types.Hash256DigestFromHex(vars["ancestor"])
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}

		currentBest, sdkErr := types.Hash256DigestFromHex(vars["currentBest"])
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}

		newBest, sdkErr := types.Hash256DigestFromHex(vars["newBest"])
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
			return
		}

		var limit uint32
		if val, ok := vars["limit"]; ok {
			lim, err := strconv.ParseUint(val, 0, 32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, sdkErr.Error())
				return
			}
			limit = uint32(lim)
		}

		params := types.QueryParamsHeaviestFromAncestor{
			Ancestor:    ancestor,
			CurrentBest: currentBest,
			NewBest:     newBest,
			Limit:       limit,
		}

		queryData, err := json.Marshal(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// run the query. the routeString is passed as strings to our querier switch/case in `keeper/querier.go`
		res, _, err := cliCtx.QueryWithData("custom/relay/heaviestfromancestor", queryData)

		// below this is boilerplate
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// handler function for getRequest queries. parses arguments from url string, and passes them through
// as a QueryParamsGetReqiest struct
func getRequestHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseUint(vars["id"], 0, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.QueryParamsGetRequest{
			ID: id,
		}

		queryData, err := json.Marshal(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/relay/getRequest", queryData)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
