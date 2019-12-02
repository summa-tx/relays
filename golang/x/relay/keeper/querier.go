package keeper

import (
	"encoding/hex"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryIsAncestor is a query string tag for IsAncestor
const QueryIsAncestor = "isancestor"

// QueryGetRelayGenesis is a query string tag for GetRelayGenesis
const QueryGetRelayGenesis = "getrelaygenesis"

// QueryGetLastReorgCA is a query string tag for GetLastReorgCA
const QueryGetLastReorgCA = "getlastreorgca"

// QueryFindAncestor is a query string tag for FindAncestor
const QueryFindAncestor = "findancestor"

// QueryHeaviestFromAncestor is a query string tag for HeaviestFromAncestor
const QueryHeaviestFromAncestor = "heaviestfromancestor"

// NewQuerier makes a query routing function
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryIsAncestor:
			return queryIsAncestor(ctx, path[1:], req, keeper)
		case QueryGetRelayGenesis:
			return queryGetRelayGenesis(ctx, req, keeper)
		case QueryGetLastReorgCA:
			return queryGetLastReorgCA(ctx, req, keeper)
		case QueryFindAncestor:
			return queryFindAncestor(ctx, path[1:], req, keeper)
		case QueryHeaviestFromAncestor:
			return queryHeaviestFromAncestor(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown relay query endpoint")
		}
	}
}

func queryIsAncestor(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// Take some URL path, parse the arguments out of it, pass it to the keeper,
	// and format the result as a QueryResIsAncestor

	// The path I expect here looks like this:
	// /getancestor/abcd1234.../second_digest/limit

	// getAncestor is removed by the handler switch/case block above

	// check that the path is this many items long, error if not
	if len(path) > 3 {
		return []byte{}, types.ErrTooManyArguments(types.DefaultCodespace)
	} else if len(path) < 3 {
		return []byte{}, types.ErrNotEnoughArguments(types.DefaultCodespace)
	}

	// parse the first path item as hex.
	hexDigest := path[0]
	digestBytes, decodeErr := hex.DecodeString(hexDigest)
	if decodeErr != nil {
		return []byte{}, types.ErrBadHex(types.DefaultCodespace)
	}
	digestLE, newDigestErr := btcspv.NewHash256Digest(digestBytes)
	if newDigestErr != nil {
		return []byte{}, types.FromBTCSPVError(types.DefaultCodespace, newDigestErr)
	}

	// parse the second path item as hex.
	ancestorDigest := path[1]
	ancestorDigestBytes, decodeErr := hex.DecodeString(ancestorDigest)
	if err != nil {
		return []byte{}, types.ErrBadHex(types.DefaultCodespace)
	}
	ancestorDigestLE, newDigestErr := btcspv.NewHash256Digest(ancestorDigestBytes)
	if newDigestErr != nil {
		return []byte{}, types.FromBTCSPVError(types.DefaultCodespace, newDigestErr)
	}

	limit, convErr := strconv.Atoi(path[2]) // TODO: parse from path, use 240 as default if not in path
	if convErr != nil {
		return []byte{}, types.FromBTCSPVError(types.DefaultCodespace, convErr)
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result := keeper.IsAncestor(ctx, ancestorDigestLE, digestLE, uint32(limit))

	// Now we format the answer as a response
	response := types.QueryResIsAncestor{
		Digest:              digestLE,
		ProspectiveAncestor: ancestorDigestLE,
		IsAncestor:          result,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}

func queryGetRelayGenesis(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// This calls the keeper and gets an answer
	result, err := keeper.GetRelayGenesis(ctx)

	// Now we format the answer as a response
	response := types.QueryResGetRelayGenesis{
		Digest: result,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}

func queryGetLastReorgCA(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// This calls the keeper and gets an answer
	result, err := keeper.GetRelayGenesis(ctx)

	// Now we format the answer as a response
	response := types.QueryResGetLastReorgCA{
		Digest: result,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}

func queryFindAncestor(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// check that the path is this many items long, error if not
	if len(path) > 2 {
		return []byte{}, types.ErrTooManyArguments(types.DefaultCodespace)
	} else if len(path) < 2 {
		return []byte{}, types.ErrNotEnoughArguments(types.DefaultCodespace)
	}

	hexDigest := path[0]
	digestBytes, decodeErr := hex.DecodeString(hexDigest)
	if decodeErr != nil {
		return []byte{}, types.ErrBadHex(types.DefaultCodespace)
	}
	digestLE, newDigestErr := btcspv.NewHash256Digest(digestBytes)
	if newDigestErr != nil {
		return []byte{}, types.FromBTCSPVError(types.DefaultCodespace, newDigestErr)
	}

	offset, convErr := strconv.Atoi(path[1])
	if convErr != nil {
		return []byte{}, types.FromBTCSPVError(types.DefaultCodespace, convErr)
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result, err := keeper.FindAncestor(ctx, digestLE, uint32(offset))

	// Now we format the answer as a response
	response := types.QueryResFindAncestor{
		Digest: result,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}

func queryHeaviestFromAncestor(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// check that the path is this many items long, error if not
	if len(path) > 4 {
		return []byte{}, types.ErrTooManyArguments(types.DefaultCodespace)
	} else if len(path) < 4 {
		return []byte{}, types.ErrNotEnoughArguments(types.DefaultCodespace)
	}

	ancestorHex := path[0]
	ancestorBytes, decodeErr := hex.DecodeString(ancestorHex)
	if decodeErr != nil {
		return []byte{}, types.ErrBadHex(types.DefaultCodespace)
	}
	ancestor, newDigestErr := btcspv.NewHash256Digest(ancestorBytes)
	if newDigestErr != nil {
		return []byte{}, types.FromBTCSPVError(types.DefaultCodespace, newDigestErr)
	}

	// parse the second path item as hex.
	currentBest := path[1]
	currentBestBytes, decodeErr := hex.DecodeString(currentBest)
	if err != nil {
		return []byte{}, types.ErrBadHex(types.DefaultCodespace)
	}
	currentBestDigest, newDigestErr := btcspv.NewHash256Digest(currentBestBytes)
	if newDigestErr != nil {
		return []byte{}, types.FromBTCSPVError(types.DefaultCodespace, newDigestErr)
	}

	newBest := path[2]
	newBestBytes, decodeErr := hex.DecodeString(newBest)
	if err != nil {
		return []byte{}, types.ErrBadHex(types.DefaultCodespace)
	}
	newBestDigest, newDigestErr := btcspv.NewHash256Digest(newBestBytes)
	if newDigestErr != nil {
		return []byte{}, types.FromBTCSPVError(types.DefaultCodespace, newDigestErr)
	}

	limit, convErr := strconv.Atoi(path[3])
	if convErr != nil {
		return []byte{}, types.FromBTCSPVError(types.DefaultCodespace, convErr)
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result, err := keeper.HeaviestFromAncestor(ctx, ancestor, currentBestDigest, newBestDigest, uint32(limit))

	// Now we format the answer as a response
	response := types.QueryResHeaviestFromAncestor{
		Digest: result,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}
