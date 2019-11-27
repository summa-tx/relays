package keeper

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryIsAncestor is a query string tag for IsAncestor
const QueryIsAncestor = "isancestor"

// NewQuerier makes a query routing function
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryIsAncestor:
			return queryIsAncestor(ctx, path[1:], req, keeper)
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

	// parse the first path item as hex.
	// TODO: check that the path is this many items long, error if not
	hexDigest := path[0]
	digestBytes, err := hex.DecodeString(hexDigest)
	if err != nil {
		return []byte{}, types.ErrBadHex(types.DefaultCodespace)
	}
	digestLE := btcspv.NewHash256Digest(digestBytes)

	// parse the second path item as hex.
	// TODO: check that the path is this many items long, error if not
	ancestorDigest := path[1]
	ancestorDigestBytes, err := hex.DecodeString(ancestorDigest)
	if err != nil {
		return []byte{}, types.ErrBadHex(types.DefaultCodespace)
	}
	ancestorDigestLE := btcspv.NewHash256Digest(ancestorDigestBytes)

	// TODO: check that the path is this many items long, error if not
	limit := 25 // TODO: parse from path, use 240 as default if not in path

	// This calls the keeper with the parsed arguments, and gets an answer
	result := keeper.IsAncestor(ctx, ancestorDigestLE, digestLE, limit)

	// Now we format the answer as a response
	response := types.QueryResIsAncestor{digestLE, ancestorDigestLE, result}

	// And we serialize that response as JSON
	res, err := codec.MarshalJSONIndent(keeper.cdc, response)
	if err != nil {
		// TODO: better handling
		panic("could not marshal result to JSON")
	}
	return res, nil
}
