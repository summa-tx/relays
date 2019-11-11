package types

// Hash256Digest 32-byte double-sha2 digest
type Hash256Digest = [32]byte

// Link is a link in the chain
type Link struct {
	Digest string `json:"digest"`
	Parent string `json:"parent"`
}
