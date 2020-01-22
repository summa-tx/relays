package types

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RequestID is an 8 byte id used to store requests
type RequestID [8]byte

// ProofRequest is info about a proof request
type ProofRequest struct {
	Spends      Hash256Digest `json:"spends"`
	Pays        Hash256Digest `json:"pays"`
	PaysValue   uint64        `json:"paysValue"`
	ActiveState bool          `json:"activeState"`
	NumConfs    uint8         `json:"numConfs"`
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
			return RequestID{}, ErrBadHex(DefaultCodespace)
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
