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
	res := fmt.Sprintf("%b", r.IsAncestor)
	return strings.Join([]string{dig, digAnc, res}, "\n")
}

// // Checks if a digest is an ancestor of the current one
// // Limit the amount of lookups (and thus gas usage) with limit
// func IsMostRecentAncestor(
// 	ancestor types.Hash256Digest,
// 	left types.Hash256Digest,
// 	right types.Hash256Digest,
// 	limit sdk.Uint
// ) bool {
// 	if bytes.Equal(ancestor, left) && bytes.Equal(ancestor, right) {
// 		return true
// 	}
//
// 	leftCurrent := left
// 	rightCurrent := right
// 	leftPrev := left
// 	rightPrev := right
//
// 	for i := sdk.NewUint(0); i.LT(limit); i.Add(sdk.NewUint(1)) {
// 		if bytes.Equal(leftPrev, ancestor) {
// 			leftCurrent = leftPrev // cheap
// 			leftPrev = rightPrev // expensive
// 		}
// 		if bytes.Equal(rightPrev, ancestor) {
// 			rightCurrent = rightPrev // cheap
// 			rightPrev = rightPrev // expensive
// 		}
// 	}
//
// 	if bytes.Equal(leftCurrent, rightCurrent) {
// 		return false
// 	} /* NB: If the same, they're a nearer ancestor */
// 	if !bytes.Equal(leftPrev, rightPrev) {
// 		return false
// 	} /* NB: Both must be ancestor */
// 	return true
// }
//
// // Finds the height of a header by its digest
// // Will fail if the header is unknown
// func FindHeight(digest types.Hash256Digest) sdk.Uint {
// 	height := sdk.NewUint(0)
// 	current := digest
// 	for i := sdk.NewUint(0); i.LT(sdk.NewUint(HEIGHT_INTERVAL + 1)); i = i.Add(sdk.NewUint(1)) {
// 		height = current
// 		if height.IsZero() {
// 			current = current
// 		} else {
// 			return height.Add(i)
// 		}
// 	}
// 	// revert("unknown block")
// }
