package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

func TestHash256DigestFromHex(t *testing.T) {
	Hash256FromHexPass := []struct {
		Input  string
		Output btcspv.Hash256Digest
	}{
		{
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			btcspv.Hash256Digest{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			btcspv.Hash256Digest{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	}

	Hash256FromHexFail := []struct {
		Input string
		Err   sdk.CodeType
	}{
		{
			"jjjjjj",
			BadHex,
		}, {
			"ffffff",
			BitcoinSPV,
		},
	}

	for i := range Hash256FromHexPass {
		digest, err := Hash256DigestFromHex(Hash256FromHexPass[i].Input)
		assert.Nil(t, err)
		assert.Equal(t, digest, Hash256FromHexPass[i].Output)
	}
	for i := range Hash256FromHexFail {
		_, err := Hash256DigestFromHex(Hash256FromHexFail[i].Input)
		assert.Equal(t, Hash256FromHexFail[i].Err, err.Code())
	}
}
