package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const (
	// DefaultLookupLimit is the default limit for lookup requests
	DefaultLookupLimit = 18

	// QueryIsAncestor is a query string tag for IsAncestor
	QueryIsAncestor = "isancestor"

	// QueryGetRelayGenesis is a query string tag for GetRelayGenesis
	QueryGetRelayGenesis = "getrelaygenesis"

	// QueryGetLastReorgLCA is a query string tag for GetLastReorgLCA
	QueryGetLastReorgLCA = "getlastreorglca"

	// QueryFindAncestor is a query string tag for FindAncestor
	QueryFindAncestor = "findancestor"

	// QueryHeaviestFromAncestor is a query string tag for HeaviestFromAncestor
	QueryHeaviestFromAncestor = "heaviestfromancestor"

	// QueryIsMostRecentCommonAncestor is a query string tag for IsMostRecentCommonAncestor
	QueryIsMostRecentCommonAncestor = "ismostrecentcommonancestor"

	// QueryGetRequest is a query string tag for getRequest
	QueryGetRequest = "getrequest"

	// QueryCheckRequest is a query string tag for checkRequests
	QueryCheckRequests = "checkrequests"
)

// QueryParamsIsAncestor represents the parameters for an IsAncestor query
type QueryParamsIsAncestor struct {
	DigestLE            Hash256Digest `json:"digest"`
	ProspectiveAncestor Hash256Digest `json:"prospectiveAncestor"`
	Limit               uint32        `json:"limit"`
}

// QueryResIsAncestor is the response to a IsAncestor query
type QueryResIsAncestor struct {
	Params QueryParamsIsAncestor `json:"params"`
	Res    bool                  `json:"result"`
}

// String formats a QueryResIsAncestor struct
func (r QueryResIsAncestor) String() string {
	dig := "0x" + hex.EncodeToString(r.Params.DigestLE[:])
	digAnc := "0x" + hex.EncodeToString(r.Params.ProspectiveAncestor[:])
	return fmt.Sprintf(
		"Digest: %s, Ancestor: %s, Limit: %d, Result: %t",
		dig, digAnc, r.Params.Limit, r.Res)
}

// QueryResGetRelayGenesis is the response struct for queryGetRelayGenesis
type QueryResGetRelayGenesis struct {
	Res Hash256Digest `json:"result"`
}

// String formats a QueryResGetRelayGenesis struct
func (r QueryResGetRelayGenesis) String() string {
	digest := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf("%s\n", digest)
}

// QueryResGetLastReorgLCA is the response struct for queryGetLastReorgLCA
type QueryResGetLastReorgLCA struct {
	Res Hash256Digest `json:"result"`
}

// String formats a QueryResGetLastReorgLCA struct
func (r QueryResGetLastReorgLCA) String() string {
	digest := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf("%s\n", digest)
}

// QueryParamsFindAncestor represents the parameters for a FindAncestor query
type QueryParamsFindAncestor struct {
	DigestLE Hash256Digest `json:"digestLE"`
	Offset   uint32        `json:"offset"`
}

// QueryResFindAncestor is the response struct for queryFindAncestor
type QueryResFindAncestor struct {
	Params QueryParamsFindAncestor `json:"params"`
	Res    Hash256Digest           `json:"result"`
}

// String formats a QueryResFindAncestor struct
func (r QueryResFindAncestor) String() string {
	dig := "0x" + hex.EncodeToString(r.Params.DigestLE[:])
	offset := r.Params.Offset
	res := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf(
		"Digest LE: %s, Offset: %d, Result: %s",
		dig, offset, res)
}

// QueryParamsHeaviestFromAncestor is the params struct for queryHeaviestFromAncestor
type QueryParamsHeaviestFromAncestor struct {
	Ancestor    Hash256Digest `json:"ancestor"`
	CurrentBest Hash256Digest `json:"currentBest"`
	NewBest     Hash256Digest `json:"newBest"`
	Limit       uint32        `json:"limit"`
}

// QueryResHeaviestFromAncestor is the response struct for queryHeaviestFromAncestor
type QueryResHeaviestFromAncestor struct {
	Params QueryParamsHeaviestFromAncestor `json:"params"`
	Res    Hash256Digest                   `json:"result"`
}

// String formats a QueryResHeaviestFromAncestor struct
func (r QueryResHeaviestFromAncestor) String() string {
	anc := "0x" + hex.EncodeToString(r.Params.Ancestor[:])
	curBest := "0x" + hex.EncodeToString(r.Params.CurrentBest[:])
	newBest := "0x" + hex.EncodeToString(r.Params.NewBest[:])
	res := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf(
		"Ancestor: %s, Current Best: %s, New Best: %s, Limit: %d, Result: %s",
		anc, curBest, newBest, r.Params.Limit, res)
}

// QueryParamsIsMostRecentCommonAncestor is the params struct for queryIsMostRecentCommonAncestor
type QueryParamsIsMostRecentCommonAncestor struct {
	Ancestor Hash256Digest `json:"ancestor"`
	Left     Hash256Digest `json:"left"`
	Right    Hash256Digest `json:"right"`
	Limit    uint32        `json:"limit"`
}

// QueryResIsMostRecentCommonAncestor is the response struct for queryIsMostRecentCommonAncestor
type QueryResIsMostRecentCommonAncestor struct {
	Params QueryParamsIsMostRecentCommonAncestor `json:"params"`
	Res    bool                                  `json:"result"`
}

// String formats a QueryResIsMostRecentCommonAncestor struct
func (r QueryResIsMostRecentCommonAncestor) String() string {
	anc := "0x" + hex.EncodeToString(r.Params.Ancestor[:])
	left := "0x" + hex.EncodeToString(r.Params.Left[:])
	right := "0x" + hex.EncodeToString(r.Params.Right[:])
	return fmt.Sprintf(
		"Ancestor: %s, Left: %s, Right: %s, Limit: %d, Result: %t",
		anc, left, right, r.Params.Limit, r.Res)
}

// QueryParamsGetRequest is the response struct for queryGetRequest
type QueryParamsGetRequest struct {
	ID RequestID `json:"id"`
}

// QueryResGetRequest is the response struct for queryGetRequest
type QueryResGetRequest struct {
	Params QueryParamsGetRequest `json:"params"`
	Res    ProofRequest          `json:"result"`
}

// String formats a QueryResIsMostRecentCommonAncestor struct
func (r QueryResGetRequest) String() string {
	spends := "0x" + hex.EncodeToString(r.Res.Spends[:])
	pays := "0x" + hex.EncodeToString(r.Res.Pays[:])
	return fmt.Sprintf(
		"ID: %d, Spends: %s, Pays: %s, Value: %d, Active: %t, Confirmations: %d",
		r.Params.ID, spends, pays, r.Res.PaysValue, r.Res.ActiveState, r.Res.NumConfs)
}

type QueryParamsCheckRequests struct {
	Filled FilledRequests `json:"filledRequests"`
}

type QueryResCheckRequests struct {
	Params       QueryParamsCheckRequests `json:"params"`
	Valid        bool                     `json:"valid"`
	ErrorMessage string                   `json:"errorMessage"`
}

func (r QueryResCheckRequests) String() string {
	json, _ := json.Marshal(r)
	return fmt.Sprint(string(json))
}
