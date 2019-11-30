package types

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type QueryResIsAncestor struct {
	Digest              Hash256Digest `json:"digest"`
	ProspectiveAncestor Hash256Digest `json:"prospectiveAncestor"`
	IsAncestor          bool          `json:"isAncestor"`
}

func (r QueryResIsAncestor) String() string {
	dig := "0x" + hex.EncodeToString(r.Digest[:])
	digAnc := "0x" + hex.EncodeToString(r.ProspectiveAncestor[:])
	res := fmt.Sprintf("%t\n", r.IsAncestor)
	return strings.Join([]string{dig, digAnc, res}, "\n")
}
