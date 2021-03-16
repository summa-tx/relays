package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
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

	// QueryGetBestDigest is a query string tag for GetBestDigest
	QueryGetBestDigest = "getbestdigest"

	// QueryFindAncestor is a query string tag for FindAncestor
	QueryFindAncestor = "findancestor"

	// QueryHeaviestFromAncestor is a query string tag for HeaviestFromAncestor
	QueryHeaviestFromAncestor = "heaviestfromancestor"

	// QueryIsMostRecentCommonAncestor is a query string tag for IsMostRecentCommonAncestor
	QueryIsMostRecentCommonAncestor = "ismostrecentcommonancestor"

	// QueryGetRequest is a query string tag for getRequest
	QueryGetRequest = "getrequest"

	// QueryCheckRequests is a query string tag for checkRequests
	QueryCheckRequests = "checkrequests"

	// QueryCheckProof is a query string tag for checkProof
	QueryCheckProof = "checkproof"
)

// QueryParamsIsAncestor represents the parameters for an IsAncestor query
type QueryParamsIsAncestor struct {
	DigestLE            btcspv.Hash256Digest
	ProspectiveAncestor btcspv.Hash256Digest
	Limit               uint32
}

// QueryResIsAncestor is the response to a IsAncestor query
type QueryResIsAncestor struct {
	Params QueryParamsIsAncestor
	Res    bool
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
	Res btcspv.Hash256Digest
}

// String formats a QueryResGetRelayGenesis struct
func (r QueryResGetRelayGenesis) String() string {
	digest := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf("%s\n", digest)
}

// QueryResGetLastReorgLCA is the response struct for queryGetLastReorgLCA
type QueryResGetLastReorgLCA struct {
	Res btcspv.Hash256Digest
}

// String formats a QueryResGetLastReorgLCA struct
func (r QueryResGetLastReorgLCA) String() string {
	digest := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf("%s\n", digest)
}

// QueryResGetBestDigest is the response struct for queryGetBestDigest
type QueryResGetBestDigest struct {
	Res btcspv.Hash256Digest
}

// String formats a QueryResGetBestDigest struct
func (r QueryResGetBestDigest) String() string {
	digest := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf("%s\n", digest)
}

// QueryParamsFindAncestor represents the parameters for a FindAncestor query
type QueryParamsFindAncestor struct {
	DigestLE btcspv.Hash256Digest
	Offset   uint32
}

// QueryResFindAncestor is the response struct for queryFindAncestor
type QueryResFindAncestor struct {
	Params QueryParamsFindAncestor
	Res    btcspv.Hash256Digest
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
	Ancestor    btcspv.Hash256Digest
	CurrentBest btcspv.Hash256Digest
	NewBest     btcspv.Hash256Digest
	Limit       uint32
}

// QueryResHeaviestFromAncestor is the response struct for queryHeaviestFromAncestor
type QueryResHeaviestFromAncestor struct {
	Params QueryParamsHeaviestFromAncestor
	Res    btcspv.Hash256Digest
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
	Ancestor btcspv.Hash256Digest
	Left     btcspv.Hash256Digest
	Right    btcspv.Hash256Digest
	Limit    uint32
}

// QueryResIsMostRecentCommonAncestor is the response struct for queryIsMostRecentCommonAncestor
type QueryResIsMostRecentCommonAncestor struct {
	Params QueryParamsIsMostRecentCommonAncestor
	Res    bool
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
	ID RequestID
}

// QueryResGetRequest is the response struct for queryGetRequest
type QueryResGetRequest struct {
	Params QueryParamsGetRequest
	Res    ProofRequest
}

// String formats a QueryResIsMostRecentCommonAncestor struct
func (r QueryResGetRequest) String() string {
	spends := "0x" + hex.EncodeToString(r.Res.Spends[:])
	pays := "0x" + hex.EncodeToString(r.Res.Pays[:])
	return fmt.Sprintf(
		"ID: %d, Spends: %s, Pays: %s, Value: %d, Active: %t, Confirmations: %d",
		r.Params.ID, spends, pays, r.Res.PaysValue, r.Res.ActiveState, r.Res.NumConfs)
}

// QueryParamsCheckRequests is the response struct for queryCheckRequests
type QueryParamsCheckRequests struct {
	Filled FilledRequests
}

// QueryResCheckRequests is the response struct for queryCheckRequests
type QueryResCheckRequests struct {
	Params       QueryParamsCheckRequests
	Valid        bool
	ErrorMessage string
}

// String formats a QueryResCheckRequests struct
func (r QueryResCheckRequests) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}

// QueryParamsCheckProof is the response struct for queryCheckProof
type QueryParamsCheckProof struct {
	Proof SPVProof
}

// QueryResCheckProof is the response struct for queryCheckProof
type QueryResCheckProof struct {
	Params       QueryParamsCheckProof
	Valid        bool
	ErrorMessage string
}

// String formats a QueryResCheckProof struct
func (r QueryResCheckProof) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}
