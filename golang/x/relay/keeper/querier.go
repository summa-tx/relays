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

// QueryGetLastReorgLCA is a query string tag for GetLastReorgLCA
const QueryGetLastReorgLCA = "getlastreorglca"

// QueryFindAncestor is a query string tag for FindAncestor
const QueryFindAncestor = "findancestor"

// QueryHeaviestFromAncestor is a query string tag for HeaviestFromAncestor
const QueryHeaviestFromAncestor = "heaviestfromancestor"

// QueryIsMostRecentCommonAncestor is a query string tag for IsMostRecentCommonAncestor
const QueryIsMostRecentCommonAncestor = "ismostrecentcommonancestor"

// hash256DigestFromHex converts a hex into a Hash256Digest
func hash256DigestFromHex(hexStr string) (types.Hash256Digest, sdk.Error) {
	bytes, decodeErr := hex.DecodeString(hexStr)
	if decodeErr != nil {
		return types.Hash256Digest{}, types.ErrBadHex(types.DefaultCodespace)
	}
	digest, newDigestErr := btcspv.NewHash256Digest(bytes)
	if newDigestErr != nil {
		return types.Hash256Digest{}, types.FromBTCSPVError(types.DefaultCodespace, newDigestErr)
	}
	return digest, nil
}

// NewQuerier makes a query routing function
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryIsAncestor:
			return queryIsAncestor(ctx, path[1:], req, keeper)
		case QueryGetRelayGenesis:
			return queryGetRelayGenesis(ctx, req, keeper)
		case QueryGetLastReorgLCA:
			return queryGetLastReorgLCA(ctx, req, keeper)
		case QueryFindAncestor:
			return queryFindAncestor(ctx, path[1:], req, keeper)
		case QueryHeaviestFromAncestor:
			return queryHeaviestFromAncestor(ctx, path[1:], req, keeper)
		case QueryIsMostRecentCommonAncestor:
			return queryIsMostRecentCommonAncestor(ctx, path[1:], req, keeper)
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
	} else if len(path) < 2 {
		return []byte{}, types.ErrNotEnoughArguments(types.DefaultCodespace)
	}

	// parse the first path item as hex.
	digestLE, digestErr := hash256DigestFromHex(path[0])
	if digestErr != nil {
		return []byte{}, digestErr
	}

	// parse the second path item as hex.
	ancestor, ancestorErr := hash256DigestFromHex(path[1])
	if ancestorErr != nil {
		return []byte{}, ancestorErr
	}

	var limit uint32
	if len(path) == 2 {
		limit = 15
	} else {
		num, convErr := strconv.ParseUint(path[2], 0, 10)
		if convErr != nil {
			return []byte{}, types.ErrExternal(types.DefaultCodespace, convErr)
		}
		limit = uint32(num)
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result := keeper.IsAncestor(ctx, ancestor, digestLE, uint32(limit))

	// Now we format the answer as a response
	response := types.QueryResIsAncestor{
		Digest:              digestLE,
		ProspectiveAncestor: ancestor,
		Res:                 result,
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
	if err != nil {
		return []byte{}, err
	}

	// Now we format the answer as a response
	response := types.QueryResGetRelayGenesis{
		Res: result,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}

func queryGetLastReorgLCA(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// This calls the keeper and gets an answer
	result, err := keeper.GetRelayGenesis(ctx)
	if err != nil {
		return []byte{}, err
	}

	// Now we format the answer as a response
	response := types.QueryResGetLastReorgLCA{
		Res: result,
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

	digestLE, digestErr := hash256DigestFromHex(path[0])
	if digestErr != nil {
		return []byte{}, digestErr
	}

	offset, convErr := strconv.ParseUint(path[1], 0, 10)
	if convErr != nil {
		return []byte{}, types.ErrExternal(types.DefaultCodespace, convErr)
	}
	newOffset := uint32(offset)

	// This calls the keeper with the parsed arguments, and gets an answer
	result, err := keeper.FindAncestor(ctx, digestLE, newOffset)
	if err != nil {
		return []byte{}, err
	}

	// Now we format the answer as a response
	response := types.QueryResFindAncestor{
		DigestLE: digestLE,
		Offset:   newOffset,
		Res:      result,
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
	} else if len(path) < 3 {
		return []byte{}, types.ErrNotEnoughArguments(types.DefaultCodespace)
	}

	ancestor, ancestorErr := hash256DigestFromHex(path[0])
	if ancestorErr != nil {
		return []byte{}, ancestorErr
	}

	// parse the second path item as hex.
	currentBestDigest, currentBestErr := hash256DigestFromHex(path[1])
	if currentBestErr != nil {
		return []byte{}, currentBestErr
	}

	newBestDigest, newBestErr := hash256DigestFromHex(path[2])
	if newBestErr != nil {
		return []byte{}, newBestErr
	}

	var limit uint32
	if len(path) == 3 {
		limit = 15
	} else {
		num, convErr := strconv.ParseUint(path[3], 0, 10)
		if convErr != nil {
			return []byte{}, types.ErrExternal(types.DefaultCodespace, convErr)
		}
		limit = uint32(num)
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result, err := keeper.HeaviestFromAncestor(ctx, ancestor, currentBestDigest, newBestDigest, limit)
	if err != nil {
		return []byte{}, err
	}

	// Now we format the answer as a response
	response := types.QueryResHeaviestFromAncestor{
		Ancestor:    ancestor,
		CurrentBest: currentBestDigest,
		NewBest:     newBestDigest,
		Limit:       limit,
		Res:         result,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}

func queryIsMostRecentCommonAncestor(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// check that the path is this many items long, error if not
	if len(path) > 4 {
		return []byte{}, types.ErrTooManyArguments(types.DefaultCodespace)
	} else if len(path) < 4 {
		return []byte{}, types.ErrNotEnoughArguments(types.DefaultCodespace)
	}

	ancestor, ancestorErr := hash256DigestFromHex(path[0])
	if ancestorErr != nil {
		return []byte{}, ancestorErr
	}

	leftDigest, leftErr := hash256DigestFromHex(path[1])
	if leftErr != nil {
		return []byte{}, leftErr
	}

	rightDigest, rightErr := hash256DigestFromHex(path[2])
	if rightErr != nil {
		return []byte{}, rightErr
	}

	var limit uint32
	if len(path) == 3 {
		limit = 15
	} else {
		num, convErr := strconv.ParseUint(path[3], 0, 10)
		if convErr != nil {
			return []byte{}, types.ErrExternal(types.DefaultCodespace, convErr)
		}
		limit = uint32(num)
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result := keeper.IsMostRecentCommonAncestor(ctx, ancestor, leftDigest, rightDigest, limit)

	// Now we format the answer as a response
	response := types.QueryResIsMostRecentCommonAncestor{
		Ancestor: ancestor,
		Left:     leftDigest,
		Right:    rightDigest,
		Limit:    limit,
		Res:      result,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}
