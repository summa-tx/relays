package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is a name for the router
const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetLink defines a SetLink message
type MsgSetLink struct {
	Address string `json:"signer"`
	Header  string `json:"header"`
}

// NewMsgSetLink is a constructor function for MsgSetLink
func NewMsgSetLink(address, header string, owner sdk.AccAddress) MsgSetLink {
	return MsgSetLink{
		address,
		header,
	}
}

// GetSigners defines whose signature is required
func (msg MsgSetLink) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return []sdk.AccAddress{}
	}
	return []sdk.AccAddress{addr}
}

// Type should return the action
func (msg MsgSetLink) Type() string { return "set_link" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetLink) ValidateBasic() sdk.Error {
	b, err := hex.DecodeString(msg.Header)
	if err != nil || len(b) != 80 {
		return ErrBadHeaderLength(DefaultCodespace)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetLink) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route should return the name of the module
func (msg MsgSetLink) Route() string { return RouterKey }

// TODO: Write AddHeaders
// AddHeaders adds headers to storage after validating.  We check
// integrity and consistency of the header chain.
func AddHeaders(anchor []byte, headers []byte, internal bool) (bool, error) {

}

// function _addHeaders(bytes memory _anchor, bytes memory _headers, bool _internal) internal returns (bool) {
// 	uint256 _height;
// 	bytes memory _header;
// 	bytes32 _currentDigest;
// 	bytes32 _previousDigest = _anchor.hash256();

// 	uint256 _target = _headers.slice(0, 80).extractTarget();
// 	uint256 _anchorHeight = _findHeight(_previousDigest);  /* NB: errors if unknown */

// 	require(
// 			_internal || _anchor.extractTarget() == _target,
// 			"Unexpected retarget on external call");
// 	require(_headers.length % 80 == 0, "Header array length must be divisible by 80");

// 	/*
// 	NB:
// 	1. check that the header has sufficient work
// 	2. check that headers are in a coherent chain (no retargets, hash links good)
// 	3. Store the block connection
// 	4. Store the height
// 	*/
// 	for (uint256 i = 0; i < _headers.length / 80; i = i.add(1)) {
// 			_header = _headers.slice(i.mul(80), 80);
// 			_height = _anchorHeight.add(i + 1);
// 			_currentDigest = _header.hash256();

// 			/*
// 			NB:
// 			if the block is already authenticated, we don't need to a work check
// 			Or write anything to state. This saves gas
// 			*/
// 			if (previousBlock[_currentDigest] == bytes32(0)) {
// 					require(
// 							abi.encodePacked(_currentDigest).reverseEndianness().bytesToUint() <= _target,
// 							"Header work is insufficient");
// 					previousBlock[_currentDigest] = _previousDigest;
// 					if (_height % HEIGHT_INTERVAL == 0) {
// 							/*
// 							NB: We store the height only every 4th header to save gas
// 							*/
// 							blockHeight[_currentDigest] = _height;
// 					}
// 			}

// 			/* NB: we do still need to make chain level checks tho */
// 			require(_header.extractTarget() == _target, "Target changed unexpectedly");
// 			require(_header.validateHeaderPrevHash(_previousDigest), "Headers do not form a consistent chain");

// 			_previousDigest = _currentDigest;
// 	}

// 	emit Extension(
// 			_anchor.hash256(),
// 			_currentDigest);
// 	return true;
// }

// TODO: write AddHeadersWithRetarget
// Adds headers to storage, performs additional validation of retarget.
func AddHeadersWithRetarget(oldPeriodStartHeader []byte, oldPeriodEndHeader []byte, headers []byte) (bool, error) {

}

// function _addHeadersWithRetarget(
// 	bytes memory _oldPeriodStartHeader,
// 	bytes memory _oldPeriodEndHeader,
// 	bytes memory _headers
// ) internal returns (bool) {
// 	/* NB: requires that both blocks are known */
// 	uint256 _startHeight = _findHeight(_oldPeriodStartHeader.hash256());
// 	uint256 _endHeight = _findHeight(_oldPeriodEndHeader.hash256());

// 	/* NB: retargets should happen at 2016 block intervals */
// 	require(
// 			_endHeight % 2016 == 2015,
// 			"Must provide the last header of the closing difficulty period");
// 	require(
// 			_endHeight == _startHeight.add(2015),
// 			"Must provide exactly 1 difficulty period");
// 	require(
// 			_oldPeriodStartHeader.extractDifficulty() == _oldPeriodEndHeader.extractDifficulty(),
// 			"Period header difficulties do not match");

// 	/* NB: This comparison looks weird because header nBits encoding truncates targets */
// 	bytes memory _newPeriodStart = _headers.slice(0, 80);
// 	uint256 _actualTarget = _newPeriodStart.extractTarget();
// 	uint256 _expectedTarget = BTCUtils.retargetAlgorithm(
// 			_oldPeriodStartHeader.extractTarget(),
// 			_oldPeriodStartHeader.extractTimestamp(),
// 			_oldPeriodEndHeader.extractTimestamp()
// 	);
// 	require(
// 			(_actualTarget & _expectedTarget) == _actualTarget,
// 			"Invalid retarget provided");

// 	// Pass all but the first through to be added
// 	return _addHeaders(_oldPeriodEndHeader, _headers, true);
// }

// TODO: write MarkNewHeaviest
// Gives a starting point for the relay. We don't check this AT ALL really. Don't use relays with bad genesis
func MarkNewHeaviest(ancestor []byte, currentBest []byte, newBest []byte, limit uint256) (bool, error) {

}

// function _markNewHeaviest(
// 	bytes32 _ancestor,
// 	bytes memory _currentBest,
// 	bytes memory _newBest,
// 	uint256 _limit
// ) internal returns (bool) {
// 	bytes32 _newBestDigest = _newBest.hash256();
// 	bytes32 _currentBestDigest = _currentBest.hash256();
// 	require(_currentBestDigest == bestKnownDigest, "Passed in best is not best known");
// 	require(
// 			previousBlock[_newBestDigest] != bytes32(0),
// 			"New best is unknown");
// 	require(
// 			_isMostRecentAncestor(_ancestor, bestKnownDigest, _newBestDigest, _limit),
// 			"Ancestor must be heaviest common ancestor");
// 	require(
// 			_heaviestFromAncestor(_ancestor, _currentBest, _newBest) == _newBestDigest,
// 			"New best hash does not have more work than previous");

// 	bestKnownDigest = _newBestDigest;
// 	lastReorgCommonAncestor = _ancestor;
// 	emit Reorg(
// 			_currentBestDigest,
// 			_newBestDigest,
// 			_ancestor);
// 	return true;
// }
