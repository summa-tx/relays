package keeper

import (
	"encoding/json"

	"github.com/summa-tx/relays/golang/x/relay/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (s *KeeperSuite) TestDecodeUint32FromPath() {
	DecodeUintPass := []struct {
		Path         []string
		Idx          int
		DefaultLimit uint32
		Output       uint32
	}{
		{
			[]string{"", "", "12"},
			2,
			15,
			12,
		}, {
			[]string{"", ""},
			2,
			15,
			15,
		},
	}

	DecodeUintFail := []struct {
		Path         []string
		Idx          int
		DefaultLimit uint32
		Err          sdk.CodeType
	}{
		{
			[]string{"", "", "aj"},
			2,
			15,
			types.ExternalError,
		},
	}

	for i := range DecodeUintPass {
		path := DecodeUintPass[i].Path
		index := DecodeUintPass[i].Idx
		limit := DecodeUintPass[i].DefaultLimit
		num, err := decodeUint32FromPath(path, index, limit)
		s.SDKNil(err)
		s.Equal(num, DecodeUintPass[i].Output)
	}
	for i := range DecodeUintFail {
		path := DecodeUintFail[i].Path
		index := DecodeUintFail[i].Idx
		limit := DecodeUintFail[i].DefaultLimit
		_, err := decodeUint32FromPath(path, index, limit)
		s.Equal(DecodeUintFail[i].Err, err.Code())
	}
}

func (s *KeeperSuite) TestNewQuerier() {
	querier := NewQuerier(s.Keeper)

	// Set up neccessary params with a bad path
	path := []string{"badpath"}

	req := abci.RequestQuery{
		Path: "custom/relay/badpath",
		Data: []byte{},
	}

	// Test that NewQuerier errors when given a bad path
	_, err := querier(s.Context, path, req)
	s.Equal(sdk.CodeType(6), err.Code())
}

func (s *KeeperSuite) TestQueryIsAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor
	querier := NewQuerier(s.Keeper)

	s.Keeper.ingestHeader(s.Context, anchor)
	err := s.Keeper.IngestHeaderChain(s.Context, headers)
	s.SDKNil(err)

	params := types.QueryParamsIsAncestor{
		DigestLE:            headers[4].Hash,
		ProspectiveAncestor: headers[1].Hash,
		Limit:               15,
	}
	marshalledParams, marshalErr := json.Marshal(params)
	s.Nil(marshalErr)

	path := []string{"isancestor"}

	req := abci.RequestQuery{
		Path: "custom/relay/isancestor",
		Data: marshalledParams,
	}

	res, err := querier(s.Context, path, req)
	s.SDKNil(err)

	var result types.QueryResIsAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(true, result.Res)

	// If Limit is 0, it will use default limit
	params = types.QueryParamsIsAncestor{
		DigestLE:            headers[4].Hash,
		ProspectiveAncestor: headers[1].Hash,
		Limit:               0,
	}
	marshalledParams, marshalErr = json.Marshal(params)
	s.Nil(marshalErr)
	req = abci.RequestQuery{
		Path: "custom/relay/isancestor",
		Data: marshalledParams,
	}

	res, err = querier(s.Context, path, req)
	s.SDKNil(err)

	unmarshallErr = types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(true, result.Res)

	// Test unmarshall error
	req = abci.RequestQuery{
		Path: "custom/relay/isancestor",
		Data: []byte{1, 1, 1, 1},
	}

	_, err = querier(s.Context, path, req)
	s.Equal(sdk.CodeType(1), err.Code())
}

func (s *KeeperSuite) TestQueryGetRelayGenesis() {
	genesis := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].Anchor
	epochStart := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].PrevEpochStart
	// Make a new querier
	querier := NewQuerier(s.Keeper)

	path := []string{"getrelaygenesis"}

	req := abci.RequestQuery{
		Path: "custom/relay/getrelaygenesis",
		Data: []byte{},
	}

	// Test that GetRelayGenesis errors if RelayGenesis is not found
	_, err := querier(s.Context, path, req)
	s.Equal(sdk.CodeType(types.BadHash256Digest), err.Code())

	// Set Genesis state
	err = s.Keeper.SetGenesisState(s.Context, genesis, epochStart)
	s.SDKNil(err)

	// Use querier handler to get RelayGenesis
	res, err := querier(s.Context, path, req)
	s.SDKNil(err)

	// Unmarshall the result and test
	var result types.QueryResGetRelayGenesis

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(genesis.Hash, result.Res)
}

func (s *KeeperSuite) TestQueryGetLastReorgLCA() {
	genesis := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].Anchor
	epochStart := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].PrevEpochStart
	querier := NewQuerier(s.Keeper)

	path := []string{"getlastreorglca"}

	req := abci.RequestQuery{
		Path: "custom/relay/getlastreorglca",
		Data: []byte{},
	}

	// Test that it errors if it doesn't find LastReorgLCA
	_, err := querier(s.Context, path, req)
	s.Equal(sdk.CodeType(types.BadHash256Digest), err.Code())

	setStateErr := s.Keeper.SetGenesisState(s.Context, genesis, epochStart)
	s.SDKNil(setStateErr)

	res, getLCAErr := querier(s.Context, path, req)
	s.SDKNil(getLCAErr)

	var result types.QueryResGetLastReorgLCA

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(genesis.Hash, result.Res)
}

func (s *KeeperSuite) TestQueryFindAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor
	querier := NewQuerier(s.Keeper)

	params := types.QueryParamsFindAncestor{
		DigestLE: headers[4].Hash,
		Offset:   2,
	}
	marshalledParams, marshalErr := json.Marshal(params)
	s.Nil(marshalErr)

	path := []string{"findancestor"}
	req := abci.RequestQuery{
		Path: "custom/relay/findancestor",
		Data: marshalledParams,
	}

	// Test that it errors if ancestor is not found
	_, findAncestorErr := querier(s.Context, path, req)
	s.Equal(sdk.CodeType(types.UnknownBlock), findAncestorErr.Code())

	// initialize data
	s.Keeper.ingestHeader(s.Context, anchor)
	ingestErr := s.Keeper.IngestHeaderChain(s.Context, headers)
	s.SDKNil(ingestErr)

	// test that it retrieves the correct ancestor
	res, err := querier(s.Context, path, req)
	s.SDKNil(err)

	var result types.QueryResFindAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	// s.SDKNil(unmarshallErr)
	s.Nil(unmarshallErr)
	s.Equal(headers[2].Hash, result.Res)

	// Test unmarshall error
	req = abci.RequestQuery{
		Path: "custom/relay/findancestor",
		Data: []byte{1, 1, 1, 1},
	}

	_, err = querier(s.Context, path, req)
	s.Equal(sdk.CodeType(1), err.Code())
}

func (s *KeeperSuite) TestQueryHeaviestFromAncestor() {
	tv := s.Fixtures.ChainTestCases.HeaviestFromAncestor
	headers := tv.Headers[0:8]
	headersWithMain := tv.Headers[0:9]
	querier := NewQuerier(s.Keeper)

	var headersWithOrphan []types.BitcoinHeader
	headersWithOrphan = append(headersWithOrphan, headers...)
	headersWithOrphan = append(headersWithOrphan, tv.Orphan)

	s.Keeper.ingestHeader(s.Context, tv.Genesis)
	err := s.Keeper.IngestHeaderChain(s.Context, headersWithMain)
	s.SDKNil(err)
	err = s.Keeper.IngestHeaderChain(s.Context, headersWithOrphan)
	s.SDKNil(err)

	params := types.QueryParamsHeaviestFromAncestor{
		Ancestor:    headers[3].Hash,
		CurrentBest: headers[5].Hash,
		NewBest:     headers[4].Hash,
		Limit:       20,
	}
	marshalledParams, marshalErr := json.Marshal(params)
	s.Nil(marshalErr)

	path := []string{"heaviestfromancestor"}

	req := abci.RequestQuery{
		Path: "custom/relay/heaviestfromancestor",
		Data: marshalledParams,
	}

	res, err := querier(s.Context, path, req)
	s.SDKNil(err)

	var result types.QueryResHeaviestFromAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(headers[5].Hash, result.Res)

	// Test that it errors if HeaviestFromAncestorErrors
	params = types.QueryParamsHeaviestFromAncestor{
		Ancestor:    tv.Headers[10].Hash,
		CurrentBest: headers[3].Hash,
		NewBest:     headers[4].Hash,
		Limit:       20,
	}
	marshalledParams, marshalErr = json.Marshal(params)
	s.Nil(marshalErr)

	req = abci.RequestQuery{
		Path: "custom/relay/heaviestfromancestor",
		Data: marshalledParams,
	}

	_, err = querier(s.Context, path, req)
	s.Equal(sdk.CodeType(types.UnknownBlock), err.Code())

	// Test that default limit is used if limit is set to zero
	params = types.QueryParamsHeaviestFromAncestor{
		Ancestor:    headers[3].Hash,
		CurrentBest: headers[5].Hash,
		NewBest:     headers[4].Hash,
		Limit:       0,
	}
	marshalledParams, marshalErr = json.Marshal(params)
	s.Nil(marshalErr)

	req = abci.RequestQuery{
		Path: "custom/relay/heaviestfromancestor",
		Data: marshalledParams,
	}

	res, err = querier(s.Context, path, req)
	s.SDKNil(err)

	unmarshallErr = types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(headers[5].Hash, result.Res)

	// Test unmarshall error
	req = abci.RequestQuery{
		Path: "custom/relay/heaviestfromancestor",
		Data: []byte{1, 1, 1, 1},
	}

	_, err = querier(s.Context, path, req)
	s.Equal(sdk.CodeType(1), err.Code())
}

func (s *KeeperSuite) TestQueryIsMostRecentCommonAncestor() {
	tv := s.Fixtures.ChainTestCases.IsMostRecentCA
	pre := tv.PreRetargetChain
	post := tv.PostRetargetChain
	querier := NewQuerier(s.Keeper)

	var postWithOrphan []types.BitcoinHeader
	postWithOrphan = append(postWithOrphan, post[:len(post)-2]...)
	postWithOrphan = append(postWithOrphan, tv.Orphan)

	err := s.Keeper.SetGenesisState(s.Context, tv.Genesis, tv.OldPeriodStart)
	s.SDKNil(err)

	err = s.Keeper.IngestHeaderChain(s.Context, pre)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.Hash, post)
	s.SDKNil(err)
	err = s.Keeper.IngestDifficultyChange(s.Context, tv.OldPeriodStart.Hash, postWithOrphan)
	s.SDKNil(err)

	params := types.QueryParamsIsMostRecentCommonAncestor{
		Ancestor: post[2].Hash,
		Left:     post[3].Hash,
		Right:    post[2].Hash,
		Limit:    5,
	}
	marshalledParams, marshalErr := json.Marshal(params)
	s.Nil(marshalErr)

	path := []string{"ismostrecentcommonancestor"}

	req := abci.RequestQuery{
		Path: "custom/relay/ismostrecentcommonancestor",
		Data: marshalledParams,
	}

	res, err := querier(s.Context, path, req)
	s.SDKNil(err)

	var result types.QueryResIsMostRecentCommonAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(true, result.Res)

	// Test that it looks up the default limit if limit is set to zero
	params = types.QueryParamsIsMostRecentCommonAncestor{
		Ancestor: post[2].Hash,
		Left:     post[3].Hash,
		Right:    post[2].Hash,
		Limit:    0,
	}
	marshalledParams, marshalErr = json.Marshal(params)
	s.Nil(marshalErr)

	req = abci.RequestQuery{
		Path: "custom/relay/ismostrecentcommonancestor",
		Data: marshalledParams,
	}

	res, err = querier(s.Context, path, req)
	s.SDKNil(err)

	unmarshallErr = types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(true, result.Res)

	// Test unmarshall error
	req = abci.RequestQuery{
		Path: "custom/relay/ismostrecentcommonancestor",
		Data: []byte{1, 1, 1, 1},
	}

	_, err = querier(s.Context, path, req)
	s.Equal(sdk.CodeType(1), err.Code())
}

func (s *KeeperSuite) TestQueryGetRequest() {
	querier := NewQuerier(s.Keeper)

	path := []string{"getrequest"}

	// bad req
	req := abci.RequestQuery{
		Path: "custom/relay/getrequest",
		Data: []byte{0},
	}

	// Errors if it cannot unmarshal req data
	_, err := querier(s.Context, path, req)
	s.Equal(sdk.CodeType(1), err.Code())

	// marshal params and set req
	params := types.QueryParamsGetRequest{
		ID: types.RequestID{},
	}
	marshalledParams, marshalErr := json.Marshal(params)
	s.Nil(marshalErr)

	req = abci.RequestQuery{
		Path: "custom/relay/getrequest",
		Data: marshalledParams,
	}

	// Errors if request is not found
	_, err = querier(s.Context, path, req)
	s.Equal(sdk.CodeType(types.UnknownRequest), err.Code())

	// Set Request
	err = s.Keeper.setRequest(s.Context, []byte{0}, []byte{0}, 0, 0, types.Local, nil)
	s.SDKNil(err)

	// Use querier handler to get request
	res, err := querier(s.Context, path, req)
	s.SDKNil(err)

	// Unmarshall the result and test
	var result types.QueryResGetRequest

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(s.Fixtures.RequestTestCases.EmptyRequest, result.Res)
}
