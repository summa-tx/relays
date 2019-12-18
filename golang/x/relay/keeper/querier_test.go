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
