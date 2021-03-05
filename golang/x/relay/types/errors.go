package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// DefaultCodespace is the default code space
	DefaultCodespace string = ModuleName

	// 100-block -- shared errors

	// BadHeaderLength means Header array length not divisible by 80
	BadHeaderLength uint32 = 101
	// BadHeaderLengthMessage is the corresponding message
	BadHeaderLengthMessage = "Header array length must be divisble by 80 but header labeled %s with header %x has length %d"

	// BadHeight occurs when a proposed descendant is below a proposed ancestor
	BadHeight uint32 = 102
	// BadHeightMessage is the corresponding message
	BadHeightMessage = "Block labeled %s with digest %x is below the ancestor height"

	// HeightMismatch occurs when blocks do not have have consecutive height increments
	HeightMismatch uint32 = 103
	// HeightMismatchMessage is the corresponding message
	HeightMismatchMessage = "Height mismatch between blocks %x and %x"

	// UnknownBlock is the error code for unknown blocks
	UnknownBlock uint32 = 104
	// UnknownBlockMessage is the corresponding message
	UnknownBlockMessage = "Unknown block labeled %s with digest %x"

	// BadHash256Digest occurs when a wrong-length hash256 digest is found
	BadHash256Digest uint32 = 105
	// BadHash256DigestMessage is the corresponding message
	BadHash256DigestMessage = "Digest %s had wrong length"

	// BadHex occurs when a hex argument couldn't be deserialized
	BadHex uint32 = 106
	// BadHexMessage is the corresponding message
	BadHexMessage = "Bad hex string in query or msg: %s"

	// BadHexLen occurs when a hex argument is the wrong length
	BadHexLen uint32 = 107
	// BadHexLenMessage is the corresponding message
	BadHexLenMessage = "Expected %d bytes in a RequestID, got %d"

	// BitcoinSPV is the code for errors bubbled up from Bitcoin SPV
	BitcoinSPV uint32 = 108
	// BitcoinSPVMessage is the corresponding message

	// AlreadyInit is the code for a second attempt to init the relay
	AlreadyInit uint32 = 109
	// AlreadyInitMessage is the corresponding message
	AlreadyInitMessage = "Relay has already set genesis state"

	// BadOffset occurs when chain has traversed to block with no link
	BadOffset uint32 = 110
	// BadOffsetMessage is the corresponding message
	BadOffsetMessage = "Reached bottom of relay chain: block with digest %x has no link"

	// 200-block -- AddHeaders

	// UnexpectedRetarget indicates a retarget was seen during AddHeaders loop
	UnexpectedRetarget uint32 = 201
	// UnexpectedRetargetMessage is the corresponding message
	UnexpectedRetargetMessage = "Target changed unexpectedly at block %x"

	// 300-block AddHeadersWithRetarget

	// WrongEnd means the end block is at the wrong height
	WrongEnd uint32 = 301
	// WrongEndMessage is the corresponding message
	WrongEndMessage = "Must provide the last header of the closing difficulty period"

	// WrongStart means the start block is at the wrong height
	WrongStart uint32 = 302
	// WrongStartMessage is the corresponding message
	WrongStartMessage = "Must provide exactly 1 difficulty period"

	// PeriodMismatch means the start and end block do not have the same difficulty
	PeriodMismatch uint32 = 303
	// PeriodMismatchMessage is the corresponding message
	PeriodMismatchMessage = "Period header difficulties do not match"

	// BadRetarget means the provided blocks did not create the expected retarget
	BadRetarget uint32 = 304
	// BadRetargetMessage is the corresponding message
	BadRetargetMessage = "Invalid retarget provided"

	// 400-block -- MarkNewBestHeight

	// LimitTooHigh indicates that the requested limit is >2016
	LimitTooHigh uint32 = 402
	// LimitTooHighMessage is the corresponding message
	LimitTooHighMessage = "Requested lookup limit must be 2016 or lower. Got %d"

	// NotBestKnown means a block should have been the best known, but wasn't
	NotBestKnown uint32 = 403
	// NotBestKnownMessage is the corresponding message
	NotBestKnownMessage = "Provided digest %x is not current best known, expecting block with hash %x"

	// NotHeaviestAncestor means a later common ancestor was found
	NotHeaviestAncestor uint32 = 404
	// NotHeaviestAncestorMessage is the corresponding message
	NotHeaviestAncestorMessage = "Ancestor %x is not heaviest common ancestor"

	// NotHeavier means the proposed new best is not heavier than the current best
	NotHeavier uint32 = 405
	// NotHeavierMessage is the corresponding message
	NotHeavierMessage = "New best received %x does not have more work than previous best %x"

	// 500-block Queries

	// MarshalJSON means there was an error marshalling a query result to json
	MarshalJSON uint32 = 503
	// MarshalJSONMessage is the corresponding message
	MarshalJSONMessage = "Could not marshal result to JSON"

	// 600-block Proof Requests

	// UnknownRequest means the request was not found
	UnknownRequest uint32 = 601
	// UnknownRequestMessage is the corresponding message
	UnknownRequestMessage = "Request not found"

	// SpendsLength means the spend value is not 36 bytes
	SpendsLength uint32 = 602
	// SpendsLengthMessage is the corresponding message
	SpendsLengthMessage = "Spends value is not 36 bytes"

	// PaysLength means the pays value is greater than 50 bytes
	PaysLength uint32 = 603
	// PaysLengthMessage is the corresponding message
	PaysLengthMessage = "Pays value is greater than 50 bytes"

	// InvalidVin means the vin is not valid
	InvalidVin uint32 = 604
	// InvalidVinMessage is the corresponding message
	InvalidVinMessage = "Vin is not valid"

	// InvalidVout means the vout is not valid
	InvalidVout uint32 = 605
	// InvalidVoutMessage is the corresponding message
	InvalidVoutMessage = "Vout is not valid"

	// ClosedRequest means the request is not active
	ClosedRequest uint32 = 606
	// ClosedRequestMessage is the corresponding message
	ClosedRequestMessage = "Request is not active"

	// RequestPays means the output does not match the pays request
	RequestPays uint32 = 607
	// RequestPaysMessage is the corresponding message
	RequestPaysMessage = "Output does not match pays for requestID %d"

	// RequestValue means the pays value and value of the output does not match
	RequestValue uint32 = 608
	// RequestValueMessage is the corresponding message
	RequestValueMessage = "Output value does not match pays value for requestID %d"

	// RequestSpends means the request spends does not match the input
	RequestSpends uint32 = 609
	// RequestSpendsMessage is the corresponding message
	RequestSpendsMessage = "Input does not match spends for requestID %d"

	// NotAncestor means the LCA is not an ancestor of the SPV Proof header
	NotAncestor uint32 = 610
	// NotAncestorMessage is the corresponding message
	NotAncestorMessage = "LCA %x not ancestor of proof header"

	// NotEnoughConfs means the proof does not have enough confirmations
	NotEnoughConfs uint32 = 611
	// NotEnoughConfsMessage is the corresponding message
	NotEnoughConfsMessage = "Not enough confirmations for requestID %d"

	// ActionLength means the pays value is greater than 50 bytes
	ActionLength uint32 = 612
	// ActionLengthMessage is the corresponding message
	ActionLengthMessage = "Action value is greater than 500 bytes"

	// 700-block External

	// ExternalError is an error from a dependency
	ExternalError uint32 = 701
)

// ErrBadHeaderLength throws an error
func ErrBadHeaderLength(codespace string, label string, digest RawHeader, length int) *sdkerrors.Error {
	return sdkerrors.Register(codespace, BadHeaderLength, fmt.Sprintf(BadHeaderLengthMessage, label, digest, length))
}

// ErrBadHeight throws an error
func ErrBadHeight(codespace string, label string, digest Hash256Digest) *sdkerrors.Error {
	return sdkerrors.Register(codespace, BadHeight, fmt.Sprintf(BadHeightMessage, label, digest))
}

// ErrHeightMismatch throws an error
func ErrHeightMismatch(codespace string, prevDigest, digest Hash256Digest) *sdkerrors.Error {
	return sdkerrors.Register(codespace, HeightMismatch, fmt.Sprintf(HeightMismatchMessage, prevDigest, digest))
}

// ErrUnknownBlock throws an error
func ErrUnknownBlock(codespace string, label string, digest Hash256Digest) *sdkerrors.Error {
	return sdkerrors.Register(codespace, UnknownBlock, fmt.Sprintf(UnknownBlockMessage, label, digest))
}

// ErrUnexpectedRetarget throws an error
func ErrUnexpectedRetarget(codespace string, rawHeader RawHeader) *sdkerrors.Error {
	return sdkerrors.Register(codespace, UnexpectedRetarget, fmt.Sprintf(UnexpectedRetargetMessage, rawHeader))
}

// ErrWrongEnd throws an error
func ErrWrongEnd(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, WrongEnd, WrongEndMessage)
}

// ErrWrongStart throws an error
func ErrWrongStart(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, WrongStart, WrongStartMessage)
}

// ErrPeriodMismatch throws an error
func ErrPeriodMismatch(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, PeriodMismatch, PeriodMismatchMessage)
}

// ErrBadRetarget throws an error
func ErrBadRetarget(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, BadRetarget, BadRetargetMessage)
}

// ErrLimitTooHigh throws an error
func ErrLimitTooHigh(codespace string, limit string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, LimitTooHigh, fmt.Sprintf(LimitTooHighMessage, limit))
}

// ErrNotBestKnown throws an error
func ErrNotBestKnown(codespace string, invalidBest, expectedBest Hash256Digest) *sdkerrors.Error {
	return sdkerrors.Register(codespace, NotBestKnown, fmt.Sprintf(NotBestKnownMessage, invalidBest, expectedBest))
}

// ErrNotHeaviestAncestor throws an error
func ErrNotHeaviestAncestor(codespace string, ancestor Hash256Digest) *sdkerrors.Error {
	return sdkerrors.Register(codespace, NotHeaviestAncestor, fmt.Sprintf(NotHeaviestAncestorMessage, ancestor))
}

// ErrNotHeavier throws an error
func ErrNotHeavier(codespace string, newBest, prevBest Hash256Digest) *sdkerrors.Error {
	return sdkerrors.Register(codespace, NotHeavier, fmt.Sprintf(NotHeavierMessage, newBest, prevBest))
}

// ErrBadHash256Digest throws an error
func ErrBadHash256Digest(codespace string, invalidDigest string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, BadHash256Digest, fmt.Sprintf(BadHash256DigestMessage, invalidDigest))
}

// ErrBadHex throws an error
func ErrBadHex(codespace string, invalidHex string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, BadHex, fmt.Sprintf(BadHexMessage, invalidHex))
}

// ErrBadHexLen throws an error
func ErrBadHexLen(codespace string, expected, actual int) *sdkerrors.Error {
	return sdkerrors.Register(codespace, BadHexLen, fmt.Sprintf(BadHexLenMessage, expected, actual))
}

// ErrAlreadyInit throws an error
func ErrAlreadyInit(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, AlreadyInit, AlreadyInitMessage)
}

// ErrBadOffset throws an error
func ErrBadOffset(codespace string, digest Hash256Digest) *sdkerrors.Error {
	return sdkerrors.Register(codespace, BadHexLen, fmt.Sprintf(BadOffsetMessage, digest))
}

// ErrMarshalJSON throws an error
func ErrMarshalJSON(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, MarshalJSON, MarshalJSONMessage)
}

// FromBTCSPVError converts a btcutils error into an sdk error
func FromBTCSPVError(codespace string, err error) *sdkerrors.Error {
	return sdkerrors.Register(codespace, BitcoinSPV, err.Error())
}

// ErrUnknownRequest throws an error
func ErrUnknownRequest(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, UnknownRequest, UnknownRequestMessage)
}

// ErrSpendsLength throws an error
func ErrSpendsLength(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, SpendsLength, SpendsLengthMessage)
}

// ErrPaysLength throws an error
func ErrPaysLength(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, PaysLength, PaysLengthMessage)
}

// ErrActionLength throws an error
func ErrActionLength(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, ActionLength, ActionLengthMessage)
}

// ErrInvalidVin throws an error
func ErrInvalidVin(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, InvalidVin, InvalidVinMessage)
}

// ErrInvalidVout throws an error
func ErrInvalidVout(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, InvalidVout, InvalidVoutMessage)
}

// ErrClosedRequest throws an error
func ErrClosedRequest(codespace string) *sdkerrors.Error {
	return sdkerrors.Register(codespace, ClosedRequest, ClosedRequestMessage)
}

// ErrRequestPays throws an error
func ErrRequestPays(codespace string, requestID RequestID) *sdkerrors.Error {
	return sdkerrors.Register(codespace, RequestPays, fmt.Sprintf(RequestPaysMessage, requestID))
}

// ErrRequestValue throws an error
func ErrRequestValue(codespace string, requestID RequestID) *sdkerrors.Error {
	return sdkerrors.Register(codespace, RequestValue, fmt.Sprintf(RequestValueMessage, requestID))
}

// ErrRequestSpends throws an error
func ErrRequestSpends(codespace string, requestID RequestID) *sdkerrors.Error {
	return sdkerrors.Register(codespace, RequestSpends, fmt.Sprintf(RequestSpendsMessage, requestID))
}

// ErrNotAncestor throws an error
func ErrNotAncestor(codespace string, lca Hash256Digest) *sdkerrors.Error {
	return sdkerrors.Register(codespace, NotAncestor, fmt.Sprintf(NotAncestorMessage, lca))
}

// ErrNotEnoughConfs throws an error
func ErrNotEnoughConfs(codespace string, requestID RequestID) *sdkerrors.Error {
	return sdkerrors.Register(codespace, NotEnoughConfs, fmt.Sprintf(NotEnoughConfsMessage, requestID))
}

// ErrExternal converts any external error into an sdk error
func ErrExternal(codespace string, err error) *sdkerrors.Error {
	return sdkerrors.Register(codespace, ExternalError, err.Error())
}
