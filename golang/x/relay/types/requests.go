package types

import (
	"encoding/hex"

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
	var h RequestID
	copied := copy(h[:], b)
	if copied != 8 || len(b) != 8 {
		return RequestID{}, ErrBadHexLen(DefaultCodespace, 8, copied)
	}
	return h, nil
}

// RequestIdFromHex converts a hex into a RequestId
func RequestIdFromHex(hexStr string) (RequestID, sdk.Error) {
	data := hexStr
	if data[:2] == "0x" {
		data = data[2:]
	}

	bytes, decodeErr := hex.DecodeString(data)
	if decodeErr != nil {
		return RequestID{}, ErrBadHex(DefaultCodespace)
	}
	id, newIdErr := NewRequestID(bytes)
	if newIdErr != nil {
		return RequestID{}, newIdErr
	}
	return id, nil
}

// // UnmarshalJSON unmarshalls 32 byte digests
// func (id *RequestID) UnmarshalJSON(b []byte) error {
// 	// Have to trim quotation marks off byte array
// 	buf, err := hex.DecodeString(string(id[2:]))
// 	if err != nil {
// 		return err
// 	}
// 	if len(buf) != 8 {
// 		return fmt.Errorf("Expected 8 bytes, got %d bytes", len(buf))
// 	}

// 	copy(id[:], buf)

// 	return nil
// }

// // MarshalJSON marashalls 32 byte digests as 0x-prepended hex
// func (id RequestID) MarshalJSON() ([]byte, error) {
// 	encoded := "\"0x" + hex.EncodeToString(id[:]) + "\""
// 	return []byte(encoded), nil
// }
