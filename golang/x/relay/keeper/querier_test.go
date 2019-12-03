package keeper

import (
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
		Err   string
	}{
		{
			"jjjjjj",
			"ERROR:\nCodespace: relay\nCode: 106\nMessage: \"Bad hex string in query or msg\"\n",
		}, {
			"ffffff",
			"ERROR:\nCodespace: relay\nCode: 107\nMessage: \"Expected 32 bytes in a Hash256Digest, got 3\"\n",
		},
	}

	for i := range Hash256FromHexPass {
		digest, err := hash256DigestFromHex(Hash256FromHexPass[i].Input)
		suite.Nil(err)
		suite.Equal(digest, Hash256FromHexPass[i].Output)
	}
	for i := range Hash256FromHexFail {
		_, err := hash256DigestFromHex(Hash256FromHexFail[i].Input)
		suite.EqualError(err, Hash256FromHexFail[i].Err)
	}
}

func (suite *KeeperSuite) TestDecodeUint32FromPath() {
	// (path []string, idx int, defaultLimit uint32) (uint32, sdk.Error)
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
		Err          string
	}{
		{
			[]string{"", "", "aj"},
			2,
			15,
			"ERROR:\nCodespace: relay\nCode: 601\nMessage: \"strconv.ParseUint: parsing \\\"aj\\\": invalid syntax\"\n",
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
		suite.EqualError(err, DecodeUintFail[i].Err)
	}
}
