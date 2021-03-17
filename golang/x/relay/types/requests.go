package types

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

// RequestID is an 8 byte id used to store requests
type RequestID [8]byte

// ProofRequest is info about a proof request
type ProofRequest struct {
	Spends      btcspv.Hash256Digest
	Pays        btcspv.Hash256Digest
	PaysValue   uint64
	ActiveState bool
	NumConfs    uint8
	Origin      Origin
	Action      HexBytes
}

// NewRequestID instantiates a RequestID from a byte slice
func NewRequestID(b []byte) (RequestID, sdk.Error) {
	if len(b) != 8 {
		return RequestID{}, ErrBadHexLen(DefaultCodespace, 8, len(b))
	}
	var h RequestID
	copied := copy(h[:], b)
	if copied != 8 {
		return RequestID{}, ErrBadHexLen(DefaultCodespace, 8, copied)
	}
	return h, nil
}

// RequestIDFromString converts a hex string or integer string into a RequestID
func RequestIDFromString(s string) (RequestID, error) {
	var idBytes []byte
	var err error

	if s[:2] == "0x" {
		idBytes, err = hex.DecodeString(s[2:])
		if err != nil {
			return RequestID{}, ErrBadHex(DefaultCodespace, s)
		}
	} else {
		id, parseErr := strconv.ParseUint(s, 10, 64)
		if parseErr != nil {
			return RequestID{}, parseErr
		}

		// convert to bytes
		binary.BigEndian.PutUint64(idBytes, id)
	}

	requestID, newIDErr := NewRequestID(idBytes)
	if newIDErr != nil {
		return RequestID{}, newIDErr
	}
	return requestID, err
}

// UnmarshalJSON unmarshalls 8 byte requestID
func (r *RequestID) UnmarshalJSON(b []byte) error {
	// Have to trim quotation marks off byte array
	buf, err := hex.DecodeString(btcspv.Strip0xPrefix(string(b[1 : len(b)-1])))
	if err != nil {
		return err
	}

	if len(buf) != 8 {
		return fmt.Errorf("Expected 8 bytes, got %d bytes", len(buf))
	}

	copy(r[:], buf)

	return nil
}

// MarshalJSON marashalls 8 byte RequestID as 0x-prepended hex
func (r RequestID) MarshalJSON() ([]byte, error) {
	encoded := "\"0x" + hex.EncodeToString(r[:]) + "\""
	return []byte(encoded), nil
}
