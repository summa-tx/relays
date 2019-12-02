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

type QueryResGetRelayGenesis struct {
	Digest Hash256Digest `json:"digest"`
}

func (r QueryResGetRelayGenesis) String() string {
	digest := "0x" + hex.EncodeToString(r.Digest[:])
	return fmt.Sprintf("%s\n", digest)
}

type QueryResGetLastReorgCA struct {
	Digest Hash256Digest `json:"digest"`
}

func (r QueryResGetLastReorgCA) String() string {
	digest := "0x" + hex.EncodeToString(r.Digest[:])
	return fmt.Sprintf("%s\n", digest)
}

type QueryResFindAncestor struct {
	Digest Hash256Digest `json:"digest"`
}

func (r QueryResFindAncestor) String() string {
	digest := "0x" + hex.EncodeToString(r.Digest[:])
	return fmt.Sprintf("%s\n", digest)
}

type QueryResHeaviestFromAncestor struct {
	Digest Hash256Digest `json:"digest"`
}

func (r QueryResHeaviestFromAncestor) String() string {
	digest := "0x" + hex.EncodeToString(r.Digest[:])
	return fmt.Sprintf("%s\n", digest)
}
