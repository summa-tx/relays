package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func (suite *KeeperSuite) TestHash256DigestFromHex() {
	Hash256FromHexPass := []struct {
		Input  string
		Output types.Hash256Digest
	}{
		{
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			types.Hash256Digest{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	}

	Hash256FromHexFail := []struct {
		Input string
		Err   sdk.CodeType
	}{
		{
			"jjjjjj",
			types.BadHex,
		}, {
			"ffffff",
			types.BitcoinSPV,
		},
	}

	for i := range Hash256FromHexPass {
		digest, err := hash256DigestFromHex(Hash256FromHexPass[i].Input)
		suite.Nil(err)
		suite.Equal(digest, Hash256FromHexPass[i].Output)
	}
	for i := range Hash256FromHexFail {
		_, err := hash256DigestFromHex(Hash256FromHexFail[i].Input)
		suite.Equal(Hash256FromHexFail[i].Err, err.Code())
	}
}

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
