package types

import "fmt"

// Hash256Digest 32-byte double-sha2 digest
type Hash256Digest [32]byte

// NewHash256Digest instantiates a Hash256Digest from a byte slice
func NewHash256Digest(b []byte) (Hash256Digest, error) {
	var h Hash256Digest
	copied := copy(h[:], b)
	if copied != 32 {
		return h, fmt.Errorf("Expected 32 bytes in a Hash256 digest, got %d", copied)
	}
	return h, nil
}

// Link is a link in the chain
type Link struct {
	Digest string `json:"digest"`
	Parent string `json:"parent"`
}
