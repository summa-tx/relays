package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (suite *KeeperSuite) TestDecodeUint32FromPath() {
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
		suite.Nil(err)
		suite.Equal(num, DecodeUintPass[i].Output)
	}
	for i := range DecodeUintFail {
		path := DecodeUintFail[i].Path
		index := DecodeUintFail[i].Idx
		limit := DecodeUintFail[i].DefaultLimit
		_, err := decodeUint32FromPath(path, index, limit)
		suite.Equal(DecodeUintFail[i].Err, err.Code())
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
	s.Equal(err.Code(), sdk.CodeType(6))
}

func (s *KeeperSuite) TestQueryIsAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor
	querier := NewQuerier(s.Keeper)

	s.Keeper.ingestHeader(s.Context, anchor)
	s.Keeper.IngestHeaderChain(s.Context, headers)

	params := types.QueryParamsIsAncestor{
		DigestLE:            headers[4].HashLE,
		ProspectiveAncestor: headers[1].HashLE,
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
	s.Nil(err)

	var result types.QueryResIsAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(result.Res, true)

	// If Limit is 0, it will use default limit
	params = types.QueryParamsIsAncestor{
		DigestLE:            headers[4].HashLE,
		ProspectiveAncestor: headers[1].HashLE,
		Limit:               0,
	}
	marshalledParams, marshalErr = json.Marshal(params)
	s.Nil(marshalErr)
	req = abci.RequestQuery{
		Path: "custom/relay/isancestor",
		Data: marshalledParams,
	}

	res, err = querier(s.Context, path, req)
	s.Nil(err)

	unmarshallErr = types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(result.Res, true)
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
	s.Equal(err.Code(), sdk.CodeType(105))

	// Set Genesis state
	err = s.Keeper.SetGenesisState(s.Context, genesis, epochStart)
	s.Nil(err)

	// Use querier handler to get RelayGenesis
	res, err := querier(s.Context, path, req)
	s.Nil(err)

	// Unmarshall the result and test
	var result types.QueryResGetRelayGenesis

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(result.Res, genesis.HashLE)
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
	s.Equal(err.Code(), sdk.CodeType(105))

	setStateErr := s.Keeper.SetGenesisState(s.Context, genesis, epochStart)
	s.Nil(setStateErr)

	res, getLCAErr := querier(s.Context, path, req)
	s.Nil(getLCAErr)

	var result types.QueryResGetLastReorgLCA

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(result.Res, genesis.HashLE)
}

func (s *KeeperSuite) TestQueryFindAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor
	querier := NewQuerier(s.Keeper)

	params := types.QueryParamsFindAncestor{
		DigestLE: headers[4].HashLE,
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
	s.Equal(findAncestorErr.Code(), sdk.CodeType(103))

	// initialize data
	s.Keeper.ingestHeader(s.Context, anchor)
	s.Keeper.IngestHeaderChain(s.Context, headers)

	// test that it retrieves the correct ancestor
	res, err := querier(s.Context, path, req)
	s.Nil(err)

	var result types.QueryResFindAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(result.Res, headers[2].HashLE)
}
