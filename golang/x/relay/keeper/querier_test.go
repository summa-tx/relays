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

// func (s *KeeperSuite) TestNewQuerier() {
// 	// testCases := s.Fixtures.HeaderTestCases.ValidateChain
// 	// handler := NewHandler(s.Keeper)

// 	// newMsg := types.NewMsgIngestHeaderChain(getAccAddress(), testCases[0].Headers)

// 	// res := handler(s.Context, newMsg)
// 	// s.Equal(res.Code, sdk.CodeType(103))

// 	// s.Keeper.ingestHeader(s.Context, testCases[0].Anchor)
// 	// res = handler(s.Context, newMsg)
// 	// s.Equal(res.Events[0].Type, "extension")
// }

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
}

func (s *KeeperSuite) TestQueryGetRelayGenesis() {
	genesis := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].Anchor
	epochStart := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].PrevEpochStart
	querier := NewQuerier(s.Keeper)

	err := s.Keeper.SetGenesisState(s.Context, genesis, epochStart)
	s.Nil(err)

	// gen, err := s.Keeper.GetRelayGenesis(s.Context)
	// s.Nil(err)
	// s.Equal(genesis.HashLE, gen)

	path := []string{"getrelaygenesis"}

	req := abci.RequestQuery{
		Path: "custom/relay/getrelaygenesis",
		Data: []byte{},
	}

	res, err := querier(s.Context, path, req)
	s.Nil(err)

	var result types.QueryResGetRelayGenesis

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(result.Res, genesis.HashLE)
}

func (s *KeeperSuite) TestQueryGetLastReorgLCA() {
	genesis := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].Anchor
	epochStart := s.Fixtures.HeaderTestCases.ValidateDiffChange[0].PrevEpochStart
	querier := NewQuerier(s.Keeper)

	err := s.Keeper.SetGenesisState(s.Context, genesis, epochStart)
	s.Nil(err)

	path := []string{"getlastreorglca"}

	req := abci.RequestQuery{
		Path: "custom/relay/getlastreorglca",
		Data: []byte{},
	}

	res, err := querier(s.Context, path, req)
	s.Nil(err)

	var result types.QueryResGetLastReorgLCA

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(result.Res, genesis.HashLE)
}

func (s *KeeperSuite) TestQueryFindAncestor() {
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	anchor := s.Fixtures.HeaderTestCases.ValidateChain[0].Anchor
	querier := NewQuerier(s.Keeper)

	s.Keeper.ingestHeader(s.Context, anchor)
	s.Keeper.IngestHeaderChain(s.Context, headers)

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

	res, err := querier(s.Context, path, req)
	s.Nil(err)

	var result types.QueryResFindAncestor

	unmarshallErr := types.ModuleCdc.UnmarshalJSON(res, &result)
	s.Nil(unmarshallErr)
	s.Equal(result.Res, headers[2].HashLE)
}
