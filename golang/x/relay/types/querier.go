package types

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type QueryResIsAncestor struct {
	Digest              Hash256Digest `json:"digest"`
	ProspectiveAncestor Hash256Digest `json:"prospectiveAncestor"`
	IsAncestor          bool          `json:isAncestor`
}

func (r QueryResIsAncestor) String() string {
	dig := "0x" + hex.EncodeToString(r.Digest[:])
	digAnc := "0x" + hex.EncodeToString(r.ProspectiveAncestor[:])
	res := fmt.Sprintf("%b", r.IsAncestor)
	return strings.Join([]string{dig, digAnc, res}, "\n")
}

//
// var relayGenesis []byte
//
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
// 			// TODO: leftPrev = previousBlock[leftPrev]
// 			leftPrev = rightPrev // expensive
// 		}
// 		if bytes.Equal(rightPrev, ancestor) {
// 			rightCurrent = rightPrev // cheap
// 			// TODO: rightPrev = previousBlock[rightPrev]
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
// // TODO: delete?  There is already a GetRelayGenesis in querier.go
// // Getter for relayGenesis.
// // This is an initialization parameter.
// func GetRelayGenesis() types.Hash256Digest {
// 	return relayGenesis
// }
//
// // Getter for relayGenesis.
// // This is updated only by calling MarkNewHeaviest
// func getLastReorgCommonAncestor() types.Hash256Digest {
// 	return lastReorgCommonAncestor
// }
//
// // Finds the height of a header by its digest
// // Will fail if the header is unknown
// func FindHeight(digest types.Hash256Digest) sdk.Uint {
// 	height := sdk.NewUint(0)
// 	current := digest
// 	for i := sdk.NewUint(0); i.LT(sdk.NewUint(HEIGHT_INTERVAL + 1)); i = i.Add(sdk.NewUint(1)) {
// 		// TODO: height = blockHeight[current]
// 		height = current
// 		if height.IsZero() {
// 			// TODO: current = previousBlock[current]
// 			current = current
// 		} else {
// 			return height.Add(i)
// 		}
// 	}
// 	// TODO: What does this do?
// 	// revert("unknown block")
// }
//
// // TODO: delete?  Already in links.go
// // Finds an ancestor for a block by its digest
// // Will fail if the header is unknown
// func FindAncestor(digest types.Hash256Digest, offset sdk.Uint) ([]byte, error) {
// 	current := digest
// 	for i := sdk.NewUint(0); i.LT(offset); i.Add(sdk.NewUint(1)) {
// 		// current = previousBlock[current]
// 	}
// 	if bytes.Equal(current, bytes.Repeat([]byte{0}, 32)) {
// 		return []byte{0}, errors.New("Unknown ancestor")
// 	}
// 	return current
// }
//
// // TODO: Delete? Already in links.go
// // Checks if a digest is an ancestor of the current one
// // Limit the amount of lookups (and thus gas usage) with limit
// func IsAncestor(
// 	ancestor types.Hash256Digest,
// 	descendant types.Hash256Digest,
// 	limit sdk.Uint
// ) bool {
// 	current := descendant
// 	/* NB: 200 gas/read, so gas is capped at ~200 * limit */
// 	for i := sdk.NewUint(0); i.LT(limit); i.Add(sdk.NewUint(1)) {
// 		if bytes.Equal(current, ancestor) {
// 			return true
// 		}
// 		// current = previousBlock[current]
// 	}
// 	return false
// }
//
// // TODO: Delete, already in chain.go
// // Decides which header is heaviest from the ancestor
// // Does not support reorgs above 2017 blocks (:
// func HeaviestFromAncestor(
// 	ancestor types.Hash256Digest,
// 	left types.Hash256Digest,
// 	right types.Hash256Digest
// ) (types.Hash256Digest, error) {
// 	ancestorHeight := FindHeight(ancestor)
// 	leftHeight := FindHeight(btcspv.Hash256(left))
// 	rightHeight := FindHeight(btcspv.Hash256(right))
//
// 	if leftHeight.LT(ancestorHeight) && rightHeight.LT(ancestorHeight) {
// 		return []byte{0}, errors.New("A descendant height is below the ancestor height")
// 	}
//
// 	/* NB: we can shortcut if one block is in a new difficulty window and the other isn't */
// 	length := sdk.NewUint(2016)
// 	nextPeriodStartHeight := ancestorHeight.Add(length).Sub(ancestorHeight % length)
// 	leftInPeriod := leftHeight.LT(nextPeriodStartHeight)
// 	rightInPeriod := rightHeight.LT(nextPeriodStartHeight)
//
// 	/*
// 		NB:
// 		1. Left is in a new window, right is in the old window. Left is heavier
// 		2. Right is in a new window, left is in the old window. Right is heavier
// 		3. Both are in the same window, choose the higher one
// 		4. They're in different new windows. Choose the heavier one
// 	*/
// 	if !leftInPeriod && rightInPeriod {
// 		return btcspv.Hash256(left), nil
// 	} else if leftInPeriod && !rightInPeriod {
// 		return btcspv.Hash256(right), nil
// 	} else if leftInPeriod && rightInPeriod {
// 		if leftHeight.GTE(rightHeight) {
// 			return btcspv.Hash256(left), nil
// 		} else {
// 			return btcspv.Hash256(right), nil
// 		}
// 	} else {
// 		if (leftHeight % length).Mul(btcspv.ExtractDifficulty(left)).LT((rightHeight % length).Mul(btcspv.ExtractDifficulty(right))) {
// 			return btcspv.Hash256(right), nil
// 		} else {
// 			return btcspv.Hash256(left), nil
// 		}
// 	}
// }
