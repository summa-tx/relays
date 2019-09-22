package types

import "strings"

// QueryResGetParent is the payload for a GetParent Query
type QueryResGetParent struct {
	Digest string `json:"digest"`
	Parent string `json:"parent"`
}

// implement fmt.Stringer
func (r QueryResGetParent) String() string {
	return strings.Join([]string{r.Digest, r.Parent}, "\n")
}
