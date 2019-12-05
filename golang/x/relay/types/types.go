package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

// Hash256Digest 32-byte double-sha2 digest
type Hash256Digest = btcspv.Hash256Digest

// Hash160Digest is a 20-byte ripemd160+sha2 hash
type Hash160Digest = btcspv.Hash160Digest

// RawHeader is an 80-byte raw header
type RawHeader = btcspv.RawHeader

// BitcoinHeader is a parsed Bitcoin header
type BitcoinHeader = btcspv.BitcoinHeader

// SPVProof is the base struct for an SPV proof
type SPVProof = btcspv.SPVProof

// Hash256DigestFromHex converts a hex into a Hash256Digest
func Hash256DigestFromHex(hexStr string) (Hash256Digest, sdk.Error) {
	data := hexStr
	if data[:2] == "0x" {
		data = data[2:]
	}

	bytes, decodeErr := hex.DecodeString(data)
	if decodeErr != nil {
		return Hash256Digest{}, ErrBadHex(DefaultCodespace)
	}
	digest, newDigestErr := btcspv.NewHash256Digest(bytes)
	if newDigestErr != nil {
		return Hash256Digest{}, FromBTCSPVError(DefaultCodespace, newDigestErr)
	}
	return digest, nil
}
