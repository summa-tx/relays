package types

import (
	"bytes"
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

// QueryResGetParent is the payload for a GetParent Query
type QueryResGetParent struct {
	Digest string `json:"digest"`
	Parent string `json:"parent"`
}

// implement fmt.Stringer
func (r QueryResGetParent) String() string {
	return strings.Join([]string{r.Digest, r.Parent}, "\n")
}

var relayGenesis []byte

// Checks if a digest is an ancestor of the current one
// Limit the amount of lookups (and thus gas usage) with limit
func IsMostRecentAncestor(ancestor []byte, left []byte, right []byte, limit sdk.Uint) bool {
	if bytes.Equal(ancestor, left) && bytes.Equal(ancestor, right) {
		return true
	}

	leftCurrent := left
	rightCurrent := right
	leftPrev := left
	rightPrev := right

	for i := sdk.NewUint(0); i.LT(limit); i.Add(sdk.NewUint(1)) {
		if bytes.Equal(leftPrev, ancestor) {
			leftCurrent = leftPrev // cheap
			// TODO: leftPrev = previousBlock[leftPrev]
			leftPrev = rightPrev // expensive
		}
		if bytes.Equal(rightPrev, ancestor) {
			rightCurrent = rightPrev // cheap
			// TODO: rightPrev = previousBlock[rightPrev]
			rightPrev = rightPrev // expensive
		}
	}

	if bytes.Equal(leftCurrent, rightCurrent) {
		return false
	} /* NB: If the same, they're a nearer ancestor */
	if !bytes.Equal(leftPrev, rightPrev) {
		return false
	} /* NB: Both must be ancestor */
	return true
}

// Getter for relayGenesis.
// This is an initialization parameter.
func GetRelayGenesis() []byte {
	return relayGenesis
}

// Getter for relayGenesis.
// This is updated only by calling MarkNewHeaviest
func getLastReorgCommonAncestor() []byte {
	return lastReorgCommonAncestor
}

// Finds the height of a header by its digest
// Will fail if the header is unknown
func FindHeight(digest []byte) sdk.Uint {
	height := sdk.NewUint(0)
	current := digest
	for i := sdk.NewUint(0); i.LT(sdk.NewUint(HEIGHT_INTERVAL + 1)); i = i.Add(sdk.NewUint(1)) {
		// TODO: height = blockHeight[current]
		height = current
		if height.IsZero() {
			// TODO: current = previousBlock[current]
			current = current
		} else {
			return height.Add(i)
		}
	}
	// TODO: What does this do?
	// revert("unknown block")
}

// Finds an ancestor for a block by its digest
// Will fail if the header is unknown
func FindAncestor(digest []byte, offset sdk.Uint) ([]byte, error) {
	current := digest
	for i := sdk.NewUint(0); i.LT(offset); i.Add(sdk.NewUint(1)) {
		// current = previousBlock[current]
	}
	if bytes.Equal(current, bytes.Repeat([]byte{0}, 32)) {
		return []byte{0}, errors.New("Unknown ancestor")
	}
	return current
}

// Checks if a digest is an ancestor of the current one
// Limit the amount of lookups (and thus gas usage) with limit
func IsAncestor(ancestor []byte, descendant []byte, limit sdk.Uint) bool {
	current := descendant
	/* NB: 200 gas/read, so gas is capped at ~200 * limit */
	for i := sdk.NewUint(0); i.LT(limit); i.Add(sdk.NewUint(1)) {
		if bytes.Equal(current, ancestor) {
			return true
		}
		// current = previousBlock[current]
	}
	return false
}

// Decides which header is heaviest from the ancestor
// Does not support reorgs above 2017 blocks (:
func HeaviestFromAncestor(ancestor []byte, left []byte, right []byte) ([]byte, error) {
	ancestorHeight := FindHeight(ancestor)
	leftHeight := FindHeight(btcspv.Hash256(left))
	rightHeight := FindHeight(btcspv.Hash256(right))

	if leftHeight.LT(ancestorHeight) && rightHeight.LT(ancestorHeight) {
		return []byte{0}, errors.New("A descendant height is below the ancestor height")
	}

	/* NB: we can shortcut if one block is in a new difficulty window and the other isn't */
	length := sdk.NewUint(2016)
	nextPeriodStartHeight := ancestorHeight.Add(length).Sub(ancestorHeight % length)
	leftInPeriod := leftHeight.LT(nextPeriodStartHeight)
	rightInPeriod := rightHeight.LT(nextPeriodStartHeight)

	/*
		NB:
		1. Left is in a new window, right is in the old window. Left is heavier
		2. Right is in a new window, left is in the old window. Right is heavier
		3. Both are in the same window, choose the higher one
		4. They're in different new windows. Choose the heavier one
	*/
	if !leftInPeriod && rightInPeriod {
		return btcspv.Hash256(left), nil
	} else if leftInPeriod && !rightInPeriod {
		return btcspv.Hash256(right), nil
	} else if leftInPeriod && rightInPeriod {
		if leftHeight.GTE(rightHeight) {
			return btcspv.Hash256(left), nil
		} else {
			return btcspv.Hash256(right), nil
		}
	} else {
		if (leftHeight % length).Mul(btcspv.ExtractDifficulty(left)).LT((rightHeight % length).Mul(btcspv.ExtractDifficulty(right))) {
			return btcspv.Hash256(right), nil
		} else {
			return btcspv.Hash256(left), nil
		}
	}
}
