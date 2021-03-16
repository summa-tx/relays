package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/summa-tx/bitcoin-spv/golang/btcspv"

	"github.com/summa-tx/relays/proto"
)

var _ proto.QueryServer;

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

func (query *QueryParamsIsAncestor) FromProto(q *proto.QueryParamsIsAncestor) (error) {
	digest, err := bufToH256(q.DigestLE)
	if err != nil {
		return err
	}

	ancestor, err := bufToH256(q.ProspectiveAncestor)
	if err != nil {
		return err
	}

	query.DigestLE = digest
	query.ProspectiveAncestor = ancestor
	query.Limit = q.Limit

	return nil
}

// QueryResIsAncestor is the response to a IsAncestor query
type QueryResIsAncestor struct {
	Params QueryParamsIsAncestor
	Res    bool
}

func (query *QueryResIsAncestor) FromProto(q *proto.QueryResIsAncestor) (error) {
	params, err := q.Params.FromProto()
	if err != nil {
		return err
	}

	query.Params = params
	query.res = res

	return nil
}

// QueryResGetRelayGenesis is the response struct for queryGetRelayGenesis
type QueryResGetRelayGenesis struct {
	Res btcspv.Hash256Digest
}

func (query *QueryResGetRelayGenesis) FromProto(q *proto.QueryResGetRelayGenesis) (error) {
	res, err := bufToH256(q.Res)
	if err != nil {
		return err
	}

	query.Res = res

	return nil
}

// QueryResGetLastReorgLCA is the response struct for queryGetLastReorgLCA
type QueryResGetLastReorgLCA struct {
	Res btcspv.Hash256Digest
}

func (query *QueryResGetLastReorgLCA) FromProto(q *proto.QueryResGetLastReorgLCA) (error) {
	res, err := bufToH256(q.Res)
	if err != nil {
		return err
	}

	query.Res = res

	return nil
}

// QueryResGetBestDigest is the response struct for queryGetBestDigest
type QueryResGetBestDigest struct {
	Res btcspv.Hash256Digest
}

func (query *QueryResGetBestDigest) FromProto(q *proto.QueryResGetBestDigest) (error) {
	res, err := bufToH256(q.Res)
	if err != nil {
		return err
	}

	query.Res = res

	return nil
}

// QueryParamsFindAncestor represents the parameters for a FindAncestor query
type QueryParamsFindAncestor struct {
	DigestLE btcspv.Hash256Digest
	Offset   uint32
}

func (query *QueryParamsFindAncestor) FromProto(q *proto.QueryParamsFindAncestor) (error) {
	digest, err := bufToH256(q.DigestLE)
	if err != nil {
		return err
	}

	query.DigestLE = digest
	query.Offset = q.Offset

	return nil
}

// QueryResFindAncestor is the response struct for queryFindAncestor
type QueryResFindAncestor struct {
	Params QueryParamsFindAncestor
	Res    btcspv.Hash256Digest
}

func (query *QueryResFindAncestor) FromProto(q *proto.QueryResFindAncestor) (error) {
	params, err := q.Params.FromProto()
	if err != nil {
		return err
	}

	res, err := bufToH256(q.Res)
	if err != nil {
		return err
	}

	query.Params = params
	query.res = res

	return nil
}

// QueryParamsHeaviestFromAncestor is the params struct for queryHeaviestFromAncestor
type QueryParamsHeaviestFromAncestor struct {
	Ancestor    btcspv.Hash256Digest
	CurrentBest btcspv.Hash256Digest
	NewBest     btcspv.Hash256Digest
	Limit       uint32
}

func (query *QueryParamsHeaviestFromAncestor) FromProto(q *proto.QueryParamsHeaviestFromAncestor) (error) {
	ancestor, err := bufToH256(q.Ancestor)
	if err != nil {
		return err
	}

	currentBest, err := bufToH256(q.CurrentBest)
	if err != nil {
		return err
	}

	newBest, err := bufToH256(NewBest)
	if err != nil {
		return err
	}

	query.Ancestor = ancestor
	query.CurrentBest = currentBest
	query.NewBest = newBest
	query.Limit = q.Limit

	return nil
}

// QueryResHeaviestFromAncestor is the response struct for queryHeaviestFromAncestor
type QueryResHeaviestFromAncestor struct {
	Params QueryParamsHeaviestFromAncestor
	Res    btcspv.Hash256Digest
}

func (query *QueryResHeaviestFromAncestor) FromProto(q *proto.QueryResHeaviestFromAncestor) (error) {
	params, err := q.Params.FromProto()
	if err != nil {
		return err
	}

	res, err := bufToH256(q.Res)
	if err != nil {
		return err
	}

	query.Params = params
	query.res = res

	return nil
}

// QueryParamsIsMostRecentCommonAncestor is the params struct for queryIsMostRecentCommonAncestor
type QueryParamsIsMostRecentCommonAncestor struct {
	Ancestor btcspv.Hash256Digest
	Left     btcspv.Hash256Digest
	Right    btcspv.Hash256Digest
	Limit    uint32
}

func (query *QueryParamsIsMostRecentCommonAncestor) FromProto(q *proto.QueryParamsIsMostRecentCommonAncestor) {
	ancestor, err := bufToH256(q.Ancestor)
	if err != nil {
		return err
	}

	left, err := bufToH256(q.Left)
	if err != nil {
		return err
	}

	right, err := bufToH256(q.Right)
	if err != nil {
		return err
	}

	query.Ancestor = ancestor
	query.Left = left
	query.Right = right
	query.Limit = q.limit

	return nil
}

// QueryResIsMostRecentCommonAncestor is the response struct for queryIsMostRecentCommonAncestor
type QueryResIsMostRecentCommonAncestor struct {
	Params QueryParamsIsMostRecentCommonAncestor
	Res    bool
}

func (query *QueryResIsMostRecentCommonAncestor) FromProto(q *proto.QueryResIsMostRecentCommonAncestor) (error) {
	params, err := q.Params.FromProto()
	if err != nil {
		return err
	}

	query.Params = params
	query.res = q.Res

	return nil
}

// QueryParamsGetRequest is the response struct for queryGetRequest
type QueryParamsGetRequest struct {
	ID RequestID
}

func (query *QueryParamsGetRequest) FromProto(q *proto.QueryParamsGetRequest) (error) {
	id, err := bufToRequestID(q.ID)
	if err != nil {
		return err
	}

	query.ID = id

	return nil
}

// QueryResGetRequest is the response struct for queryGetRequest
type QueryResGetRequest struct {
	Params QueryParamsGetRequest
	Res    ProofRequest
}

func (query *QueryResGetRequest) FromProto(q *proto.QueryResGetRequest) (error) {
	params, err := q.Params.FromProto()
	if err != nil {
		return err
	}

	res, err := proofRequestFromProto(q.Res)
	if err != nil {
		return err
	}

	query.Params = params
	query.res = res

	return nil
}

// QueryParamsCheckRequests is the response struct for queryCheckRequests
type QueryParamsCheckRequests struct {
	Filled FilledRequests
}

func (query *QueryParamsCheckRequests) FromProto(q *proto.QueryParamsCheckRequests) (error) {
	filledRequests, err := filledRequestsFromProto(q)
	if err != nil {
		return err
	}

	query.Filled = filledRequests

	return nil
}

// QueryResCheckRequests is the response struct for queryCheckRequests
type QueryResCheckRequests struct {
	Params       QueryParamsCheckRequests
	Valid        bool
	ErrorMessage string
}

func (query *QueryResCheckRequests) FromProto(q *proto.QueryResCheckRequests) (error) {
	params, err := q.Params.FromProto()
	if err != nil {
		return err
	}

	query.Params = params
	query.Valid = q.Valid
	query.ErrorMessage = q.ErrorMessage

	return nil
}

// QueryParamsCheckProof is the response struct for queryCheckProof
type QueryParamsCheckProof struct {
	Proof SPVProof
}

func (query *QueryParamsCheckProof) FromProto(q *proto.QueryParamsCheckProof) (error) {
	spvProof, err :=spvProofFromProto(q)
	if err != nil {
		return err
	}

	query.Proof = spvProof

	return nil
}

// QueryResCheckProof is the response struct for queryCheckProof
type QueryResCheckProof struct {
	Params       QueryParamsCheckProof
	Valid        bool
	ErrorMessage string
}

func (query *QueryResCheckProof) FromProto(q *proto.QueryResCheckProof) (error) {
	params, err := q.Params.FromProto()
	if err != nil {
		return err
	}

	query.Params = params
	query.Valid = q.Valid
	query.ErrorMessage = q.ErrorMessage

	return nil
}
