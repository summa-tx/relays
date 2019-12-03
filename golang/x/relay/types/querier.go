package types

import (
	"encoding/hex"
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

// QueryResFindAncestor is the response struct for queryFindAncestor
type QueryResFindAncestor struct {
	DigestLE Hash256Digest `json:"digestLE"`
	Offset   uint32        `json:"offset"`
	Res      Hash256Digest `json:"result"`
}

// String formats a QueryResFindAncestor struct
func (r QueryResFindAncestor) String() string {
	dig := "0x" + hex.EncodeToString(r.DigestLE[:])
	offset := r.Offset
	res := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf(
		"Digest LE: %s, Offset: %d, Result: %s",
		dig, offset, res)
}

// QueryResHeaviestFromAncestor is the response struct for queryHeaviestFromAncestor
type QueryResHeaviestFromAncestor struct {
	Ancestor    Hash256Digest `json:"ancestor"`
	CurrentBest Hash256Digest `json:"currentBest"`
	NewBest     Hash256Digest `json:"newBest"`
	Limit       uint32        `json:"limit"`
	Res         Hash256Digest `json:"result"`
}

// String formats a QueryResHeaviestFromAncestor struct
func (r QueryResHeaviestFromAncestor) String() string {
	anc := "0x" + hex.EncodeToString(r.Ancestor[:])
	curBest := "0x" + hex.EncodeToString(r.CurrentBest[:])
	newBest := "0x" + hex.EncodeToString(r.NewBest[:])
	res := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf(
		"Ancestor: %s, Current Best: %s, New Best: %s, Limit: %d, Result: %s",
		anc, curBest, newBest, r.Limit, res)
}

// QueryResIsMostRecentCommonAncestor is the response struct for queryIsMostRecentCommonAncestor
type QueryResIsMostRecentCommonAncestor struct {
	Ancestor Hash256Digest `json:"ancestor"`
	Left     Hash256Digest `json:"left"`
	Right    Hash256Digest `json:"right"`
	Limit    uint32        `json:"limit"`
	Res      bool          `json:"result"`
}

// String formats a QueryResIsMostRecentCommonAncestor struct
func (r QueryResIsMostRecentCommonAncestor) String() string {
	anc := "0x" + hex.EncodeToString(r.Ancestor[:])
	left := "0x" + hex.EncodeToString(r.Left[:])
	right := "0x" + hex.EncodeToString(r.Right[:])
	return fmt.Sprintf(
		"Ancestor: %s, Left: %s, Right: %s, Limit: %d, Result: %t",
		anc, left, right, r.Limit, r.Res)
}
