package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func decodeUint32FromPath(path []string, idx int, defaultLimit uint32) (uint32, sdk.Error) {
	if idx+1 > len(path) {
		return defaultLimit, nil
	}
	// parse int from path[idx], return error if necessary
	num, convErr := strconv.ParseUint(path[idx], 0, 32)
	if convErr != nil {
		return defaultLimit, types.ErrExternal(types.DefaultCodespace, convErr)
	}
	return uint32(num), nil
}

// NewQuerier makes a query routing function
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryIsAncestor:
			return queryIsAncestor(ctx, req, keeper)
		case types.QueryGetRelayGenesis:
			return queryGetRelayGenesis(ctx, req, keeper)
		case types.QueryGetLastReorgLCA:
			return queryGetLastReorgLCA(ctx, req, keeper)
		case types.QueryFindAncestor:
			return queryFindAncestor(ctx, req, keeper)
		case types.QueryHeaviestFromAncestor:
			return queryHeaviestFromAncestor(ctx, req, keeper)
		case types.QueryIsMostRecentCommonAncestor:
			return queryIsMostRecentCommonAncestor(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown relay query endpoint")
		}
	}
}

func queryIsAncestor(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params types.QueryParamsIsAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if unmarshallErr != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", unmarshallErr))
	}

	limit := params.Limit
	if limit == 0 {
		limit = types.DefaultLookupLimit
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result := keeper.IsAncestor(ctx, params.DigestLE, params.ProspectiveAncestor, limit)

	// Now we format the answer as a response
	response := types.QueryResIsAncestor{
		Params: params,
		Res:    result,
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
	result, err := keeper.GetLastReorgLCA(ctx)
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

func queryFindAncestor(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params types.QueryParamsFindAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if unmarshallErr != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", unmarshallErr))
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result, err := keeper.FindAncestor(ctx, params.DigestLE, params.Offset)
	if err != nil {
		return []byte{}, err
	}

	// Now we format the answer as a response
	response := types.QueryResFindAncestor{
		Params: params,
		Res:    result,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}

func queryHeaviestFromAncestor(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params types.QueryParamsHeaviestFromAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if unmarshallErr != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", unmarshallErr))
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result, err := keeper.HeaviestFromAncestor(ctx, params.Ancestor, params.CurrentBest, params.NewBest, params.Limit)
	if err != nil {
		return []byte{}, err
	}

	// Now we format the answer as a response
	response := types.QueryResHeaviestFromAncestor{
		Params: params,
		Res:    result,
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
	} else if len(path) < 3 {
		return []byte{}, types.ErrNotEnoughArguments(types.DefaultCodespace)
	}

	ancestor, ancestorErr := types.Hash256DigestFromHex(path[0])
	if ancestorErr != nil {
		return []byte{}, ancestorErr
	}

	leftDigest, leftErr := types.Hash256DigestFromHex(path[1])
	if leftErr != nil {
		return []byte{}, leftErr
	}

	rightDigest, rightErr := types.Hash256DigestFromHex(path[2])
	if rightErr != nil {
		return []byte{}, rightErr
	}

	limit, limitErr := decodeUint32FromPath(path, 3, 15)
	if limitErr != nil {
		return []byte{}, limitErr
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
