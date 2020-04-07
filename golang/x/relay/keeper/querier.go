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
	num, convErr := strconv.ParseUint(path[idx], 10, 32)
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
		case types.QueryGetBestDigest:
			return queryGetBestDigest(ctx, req, keeper)
		case types.QueryFindAncestor:
			return queryFindAncestor(ctx, req, keeper)
		case types.QueryHeaviestFromAncestor:
			return queryHeaviestFromAncestor(ctx, req, keeper)
		case types.QueryIsMostRecentCommonAncestor:
			return queryIsMostRecentCommonAncestor(ctx, req, keeper)
		case types.QueryGetRequest:
			return queryGetRequest(ctx, req, keeper)
		case types.QueryCheckRequests:
			return queryCheckRequests(ctx, req, keeper)
		case types.QueryCheckProof:
			return queryCheckProof(ctx, req, keeper)
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

func queryGetBestDigest(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	// This calls the keeper and gets an answer
	result, err := keeper.GetBestKnownDigest(ctx)
	if err != nil {
		return []byte{}, err
	}

	// Now we format the answer as a response
	response := types.QueryResGetBestDigest{
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

	limit := params.Limit
	if limit == 0 {
		limit = types.DefaultLookupLimit
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result, err := keeper.HeaviestFromAncestor(ctx, params.Ancestor, params.CurrentBest, params.NewBest, limit)
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

func queryIsMostRecentCommonAncestor(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params types.QueryParamsIsMostRecentCommonAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if unmarshallErr != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", unmarshallErr))
	}

	limit := params.Limit
	if limit == 0 {
		limit = types.DefaultLookupLimit
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result := keeper.IsMostRecentCommonAncestor(ctx, params.Ancestor, params.Left, params.Right, limit)

	// Now we format the answer as a response
	response := types.QueryResIsMostRecentCommonAncestor{
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

func queryGetRequest(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params types.QueryParamsGetRequest

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if unmarshallErr != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", unmarshallErr))
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	result, resErr := keeper.getRequest(ctx, params.ID)
	if resErr != nil {
		return []byte{}, resErr
	}

	// Now we format the answer as a response
	response := types.QueryResGetRequest{
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

func queryCheckRequests(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params types.QueryParamsCheckRequests
	var errMsg string
	valid := true

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if unmarshallErr != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", unmarshallErr))
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	resErr := keeper.checkRequestsFilled(ctx, params.Filled)
	if resErr != nil {
		valid = false
		errMsg = resErr.Error()
	}

	// Now we format the answer as a response
	response := types.QueryResCheckRequests{
		Params:       params,
		Valid:        valid,
		ErrorMessage: errMsg,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}

func queryCheckProof(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	var params types.QueryParamsCheckProof
	var errMsg string
	valid := true

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if unmarshallErr != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", unmarshallErr))
	}

	// This calls the keeper with the parsed arguments, and gets an answer
	resErr := keeper.validateProof(ctx, params.Proof)
	if resErr != nil {
		valid = false
		errMsg = resErr.Error()
	}

	// Now we format the answer as a response
	response := types.QueryResCheckProof{
		Params:       params,
		Valid:        valid,
		ErrorMessage: errMsg,
	}

	// And we serialize that response as JSON
	res, marshalErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if marshalErr != nil {
		return []byte{}, types.ErrMarshalJSON(types.DefaultCodespace)
	}
	return res, nil
}
