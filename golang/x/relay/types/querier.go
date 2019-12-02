package types

import (
	"encoding/hex"
	"fmt"
)

type QueryResIsAncestor struct {
	Digest              Hash256Digest `json:"digest"`
	ProspectiveAncestor Hash256Digest `json:"prospectiveAncestor"`
	Res                 bool          `json:"res"`
}

func (r QueryResIsAncestor) String() string {
	dig := "0x" + hex.EncodeToString(r.Digest[:])
	digAnc := "0x" + hex.EncodeToString(r.ProspectiveAncestor[:])
	// res := fmt.Sprintf("%t\n", r.Res)
	// return strings.Join([]string{dig, digAnc, res}, "\n")
	return fmt.Sprintf(
		"Digest: %s, Ancestor: %s, Result: %t",
		dig, digAnc, r.Res)
}

// TODO: update these structs to include arguments
type QueryResGetRelayGenesis struct {
	Res Hash256Digest `json:"res"`
}

func (r QueryResGetRelayGenesis) String() string {
	digest := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf("%s\n", digest)
}

type QueryResGetLastReorgCA struct {
	Res Hash256Digest `json:"res"`
}

func (r QueryResGetLastReorgCA) String() string {
	digest := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf("%s\n", digest)
}

type QueryResFindAncestor struct {
	DigestLE Hash256Digest `json:"digestLE"`
	Offset   uint32        `json:"offset"`
	Res      Hash256Digest `json:"res"`
}

func (r QueryResFindAncestor) String() string {
	dig := "0x" + hex.EncodeToString(r.DigestLE[:])
	offset := r.Offset
	res := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf(
		"Digest LE: %s, Offset: %d, Result: %s",
		dig, offset, res)
}

type QueryResHeaviestFromAncestor struct {
	Ancestor    Hash256Digest `json:"ancestor"`
	CurrentBest Hash256Digest `json:"currentBest"`
	NewBest     Hash256Digest `json:"newBest"`
	Limit       uint32        `json:"limit"`
	Res         Hash256Digest `json:"res"`
}

func (r QueryResHeaviestFromAncestor) String() string {
	anc := "0x" + hex.EncodeToString(r.Ancestor[:])
	curBest := "0x" + hex.EncodeToString(r.CurrentBest[:])
	newBest := "0x" + hex.EncodeToString(r.NewBest[:])
	limit := r.Limit
	res := "0x" + hex.EncodeToString(r.Res[:])
	return fmt.Sprintf(
		"Ancestor: %s, Current Best: %s, New BestL: %s, Limit: %d, Result: %s",
		anc, curBest, newBest, limit, res)
}
