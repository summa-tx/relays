package types

import (
	"bytes"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	return true
}

// function _isMostRecentAncestor(
// 	bytes32 _ancestor,
// 	bytes32 _left,
// 	bytes32 _right,
// 	uint256 _limit
// ) internal view returns (bool) {
// 	/* NB: sure why not */
// 	if (_ancestor == _left && _ancestor == _right) {
// 			return true;
// 	}

// 	bytes32 _leftCurrent = _left;
// 	bytes32 _rightCurrent = _right;
// 	bytes32 _leftPrev = _left;
// 	bytes32 _rightPrev = _right;

// 	for(uint256 i = 0; i < _limit; i = i.add(1)) {
// 			if (_leftPrev != _ancestor) {
// 					_leftCurrent = _leftPrev;  // cheap
// 					_leftPrev = previousBlock[_leftPrev];  // expensive
// 			}
// 			if (_rightPrev != _ancestor) {
// 					_rightCurrent = _rightPrev;  // cheap
// 					_rightPrev = previousBlock[_rightPrev];  // expensive
// 			}
// 	}
// 	if (_leftCurrent == _rightCurrent) {return false;} /* NB: If the same, they're a nearer ancestor */
// 	if (_leftPrev != _rightPrev) {return false;} /* NB: Both must be ancestor */
// 	return true;
// }

// Getter for relayGenesis.
// This is an initialization parameter.
func GetRelayGenesis() []byte {
	return relayGenesis
}

// function getRelayGenesis() public view returns (bytes32) {
// 	return relayGenesis;
// }

// Getter for relayGenesis.
// This is updated only by calling MarkNewHeaviest
func getLastReorgCommonAncestor() []byte {
	return lastReorgCommonAncestor
}

// Finds the height of a header by its digest
// Will fail if the header is unknown
func FindHeight(digest []byte) sdk.Uint {

}

// function _findHeight(bytes32 _digest) internal view returns (uint256) {
// 	uint256 _height = 0;
// 	bytes32 _current = _digest;
// 	for (uint256 i = 0; i < HEIGHT_INTERVAL + 1; i = i.add(1)) {
// 			_height = blockHeight[_current];
// 			if (_height == 0) {
// 					_current = previousBlock[_current];
// 			} else {
// 					return _height.add(i);
// 			}
// 	}
// 	revert("Unknown block");
// }

// Finds an ancestor for a block by its digest
// Will fail if the header is unknown
func FindAncestor(digest []byte, offset sdk.Uint) []byte {

}

// function _findAncestor(bytes32 _digest, uint256 _offset) internal view returns (bytes32) {
// 	bytes32 _current = _digest;
// 	for (uint256 i = 0; i < _offset; i = i.add(1)) {
// 			_current = previousBlock[_current];
// 	}
// 	require(_current != bytes32(0), "Unknown ancestor");
// 	return _current;
// }

// Checks if a digest is an ancestor of the current one
// Limit the amount of lookups (and thus gas usage) with limit
func IsAncestor(ancestor []byte, descendant []byte, limit sdk.Uint) bool {
	return false
}

// function _isAncestor(bytes32 _ancestor, bytes32 _descendant, uint256 _limit) internal view returns (bool) {
// 	bytes32 _current = _descendant;
// 	/* NB: 200 gas/read, so gas is capped at ~200 * limit */
// 	for (uint256 i = 0; i < _limit; i = i.add(1)) {
// 			if (_current == _ancestor) {
// 					return true;
// 			}
// 			_current = previousBlock[_current];
// 	}
// 	return false;
// }

// Decides which header is heaviest from the ancestor
// Does not support reorgs above 2017 blocks (:
func HeaviestFromAncestor(ancestor []byte, left []byte, right []byte) []byte {

}

// function _heaviestFromAncestor(
// 	bytes32 _ancestor,
// 	bytes memory _left,
// 	bytes memory _right
// ) internal view returns (bytes32) {
// 	uint256 _ancestorHeight = _findHeight(_ancestor);
// 	uint256 _leftHeight = _findHeight(_left.hash256());
// 	uint256 _rightHeight = _findHeight(_right.hash256());

// 	require(
// 			_leftHeight >= _ancestorHeight && _rightHeight >= _ancestorHeight,
// 			"A descendant height is below the ancestor height");

// 	/* NB: we can shortcut if one block is in a new difficulty window and the other isn't */
// 	uint256 _nextPeriodStartHeight = _ancestorHeight.add(2016).sub(_ancestorHeight % 2016);
// 	bool _leftInPeriod = _leftHeight < _nextPeriodStartHeight;
// 	bool _rightInPeriod = _rightHeight < _nextPeriodStartHeight;

// 	/*
// 	NB:
// 	1. Left is in a new window, right is in the old window. Left is heavier
// 	2. Right is in a new window, left is in the old window. Right is heavier
// 	3. Both are in the same window, choose the higher one
// 	4. They're in different new windows. Choose the heavier one
// 	*/
// 	if (!_leftInPeriod && _rightInPeriod) {return _left.hash256();}
// 	if (_leftInPeriod && !_rightInPeriod) {return _right.hash256();}
// 	if (_leftInPeriod && _rightInPeriod) {
// 			return _leftHeight >= _rightHeight ? _left.hash256() : _right.hash256();
// 	} else {  // if (!_leftInPeriod && !_rightInPeriod) {
// 			if (((_leftHeight % 2016).mul(_left.extractDifficulty())) <
// 					(_rightHeight % 2016).mul(_right.extractDifficulty())) {
// 					return _right.hash256();
// 			} else {
// 					return _left.hash256();
// 			}
// 	}
// }
