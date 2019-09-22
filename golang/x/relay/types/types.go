package types

// Link is a link in the chain
type Link struct {
	Digest string `json:"digest"`
	Parent string `json:"parent"`
}
