package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

// ProofHandler is an interface to which the keeper dispatches valid proofs
type ProofHandler interface {
	HandleValidProof(ctx sdk.Context, filled FilledRequests, requests []ProofRequest)
}

// Hash256Digest 32-byte double-sha2 digest
type Hash256Digest = btcspv.Hash256Digest

// Hash160Digest is a 20-byte ripemd160+sha2 hash
type Hash160Digest = btcspv.Hash160Digest

// RawHeader is an 80-byte raw header
type RawHeader = btcspv.RawHeader

// HexBytes is a type alias to make JSON hex ser/deser easier
type HexBytes = btcspv.HexBytes

// BitcoinHeader is a parsed Bitcoin header
type BitcoinHeader = btcspv.BitcoinHeader

// SPVProof is the base struct for an SPV proof
type SPVProof = btcspv.SPVProof

// Origin is an enum of types denoting requests either from the local chain
// or a remote chain
type Origin int

// Origin possible types
const (
	Local  Origin = 0
	Remote Origin = 1
)

// Hash256DigestFromHex converts a hex into a Hash256Digest
func Hash256DigestFromHex(hexStr string) (Hash256Digest, sdk.Error) {
	data := hexStr
	if data[:2] == "0x" {
		data = data[2:]
	}

	bytes, decodeErr := hex.DecodeString(data)
	if decodeErr != nil {
		return Hash256Digest{}, ErrBadHex(DefaultCodespace, hexStr)
	}
	digest, newDigestErr := btcspv.NewHash256Digest(bytes)
	if newDigestErr != nil {
		return Hash256Digest{}, FromBTCSPVError(DefaultCodespace, newDigestErr)
	}
	return digest, nil
}

// NullHandler does nothing
type NullHandler struct{}

// HandleValidProof handles a valid proof (by doing nothing)
func (n NullHandler) HandleValidProof(ctx sdk.Context, filled FilledRequests, requests []ProofRequest) {
}

// NewNullHandler instantiates a new null handler
func NewNullHandler() NullHandler {
	return NullHandler{}
}
