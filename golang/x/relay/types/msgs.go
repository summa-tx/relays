package types

import (
	"bytes"
	"encoding/hex"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

// RouterKey is a name for the router
const RouterKey = ModuleName // this was defined in your key.go file

const HEIGHT_INTERVAL = 4

var bestKnownDigest []byte
var lastReorgCommonAncestor []byte

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

// TODO: Add digest type that is 32 bytes long
// AddHeaders adds headers to storage after validating.  We check
// integrity and consistency of the header chain.
func AddHeaders(anchor []byte, headers []byte, internal bool) (bool, error) {
	var height sdk.Uint
	var header []byte
	var currentDigest []byte
	var previousDigest []byte = btcspv.Hash256(anchor)

	target := btcspv.ExtractTarget(headers[0:80])
	// TODO: write findHeight in queries
	anchorHeight := FindHeight(previousDigest) /* NB: errors if unknown */
	extractedTarget := btcspv.ExtractTarget(anchor)
	if !(internal || extractedTarget == target) {
		return false, errors.New("Header array length must be divisible by 80")
	}

	/*
		NB:
		1. check that the header has sufficient work
		2. check that headers are in a coherent chain (no retargets, hash links good)
		3. Store the block connection
		4. Store the height
	*/
	for i := 0; i < len(headers)/80; i++ {
		start := i * 80
		header := headers[start:(start + 80)]
		height := anchorHeight.Add(sdk.NewUint(uint64(i + 1)))
		currentDigest = btcspv.Hash256(header)
		/*
			NB:
			if the block is already authenticated, we don't need to a work check
			Or write anything to state. This saves gas
		*/
		if bytes.Equal(currentDigest, bytes.Repeat([]byte{0}, 32)) {
			if btcspv.BytesToBigUint(btcspv.ReverseEndianness(currentDigest)).LTE(target) {
				return false, errors.New("Header work is insufficient")
			}
			currentDigest = previousDigest
			if height%HEIGHT_INTERVAL == 0 {
				/*
					NB: We store the height only every 4th header to save gas
				*/
				currentDigest = height
			}
		}
		/* NB: we do still need to make chain level checks tho */
		if !(btcspv.ExtractTarget(header) == target) {
			return false, errors.New("Target changed unexpectedly")
		}
		if !(btcspv.ValidateHeaderPrevHash(header, previousDigest)) {
			return false, errors.New("Headers do not form a consistent chain")
		}

		previousDigest = currentDigest
	}
	// TODO: How to do this in go...
	// emit Extension(
	// btcspv.Hash256(anchor),
	// currentDigest);
	return true, nil
}

// Adds headers to storage, performs additional validation of retarget.
func AddHeadersWithRetarget(oldPeriodStartHeader []byte, oldPeriodEndHeader []byte, headers []byte) (bool, error) {
	/* NB: requires that both blocks are known */
	startHeight := FindHeight(btcspv.Hash256(oldPeriodStartHeader))
	endHeight := FindHeight(btcspv.Hash256(oldPeriodEndHeader))

	/* NB: retargets should happen at 2016 block intervals */
	if endHeight%2016 != 2015 {
		return false, errors.New("Must provide the last header of the closing difficulty period")
	} else if endHeight != startHeight.Add(sdk.NewUint(2015)) {
		return false, errors.New("Must provide exactly 1 difficulty period")
	} else if btcspv.ExtractDifficulty(oldPeriodStartHeader) != btcspv.ExtractDifficulty(oldPeriodEndHeader) {
		return false, errors.New("Period header difficulties do not match")
	}

	/* NB: This comparison looks weird because header nBits encoding truncates targets */

	newPeriodStart := headers[0:80]
	actualTarget := btcspv.ExtractTarget(newPeriodStart)
	expectedTarget := btcspv.RetargetAlgorithm(
		btcspv.ExtractTarget(oldPeriodStartHeader),
		btcspv.ExtractTimestamp(oldPeriodStartHeader),
		btcspv.ExtractTimestamp(oldPeriodEndHeader))
	// TODO: Fix next line &
	if actualTarget&expectedTarget != actualTarget {
		return false, errors.New("Invalid retarget provided")
	}

	// Pass all but the first through to be added
	return AddHeaders(oldPeriodEndHeader, headers, true)

	return true, nil
}

// TODO: write isMostRecentAncestor and heaviestFromAncestor
// Gives a starting point for the relay. We don't check this AT ALL really. Don't use relays with bad genesis
func MarkNewHeaviest(ancestor []byte, currentBest []byte, newBest []byte, limit sdk.Uint) (bool, error) {
	newBestDigest := btcspv.Hash256(newBest)
	currentBestDigest := btcspv.Hash256(currentBest)
	// TODO: Where is bestKnownDigest coming from?
	if !bytes.Equal(currentBestDigest, bestKnownDigest) {
		return false, errors.New("Passed in best is not best known")
	} else if bytes.Equal(newBestDigest, bytes.Repeat([]byte{0}, 32)) {
		return false, errors.New("New best is unknown")
	} else if !IsMostRecentAncestor(ancestor, bestKnownDigest, newBestDigest, limit) {
		return false, errors.New("Ancestor must be heaviest common ancestor")
	} else if !bytes.Equal(
		HeaviestFromAncestor(ancestor, currentBest, newBest),
		newBestDigest) {
		return false, errors.New("New best hash does not have more work than previous")
	}
	bestKnownDigest = newBestDigest
	lastReorgCommonAncestor = ancestor
	// TODO: Figure out how to do this in go
	// emit Reorg(
	// 	currentBestDigest,
	// 	newBestDigest,
	// 	ancestor);

	return true, nil
}
