package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultCodespace is the default code space
	DefaultCodespace sdk.CodespaceType = ModuleName

	// 100-block -- shared errors

	// BadHeaderLength means Header array length not divisible by 80
	BadHeaderLength sdk.CodeType = 101
	// BadHeaderLengthMessage is the corresponding message
	BadHeaderLengthMessage = "Header array length must be divisble by 80 but header labeled %s with header %x has length %d"

	// UnknownBlock is the error code for unknown blocks
	UnknownBlock sdk.CodeType = 103
	// UnknownBlockMessage is the corresponding message
	UnknownBlockMessage = "Unknown block labeled %s with digest %x"

	// BadHeight occurs when a proposed descendant is below a proposed ancestor
	BadHeight sdk.CodeType = 104
	// BadHeightMessage is the corresponding message
	BadHeightMessage = "Invalid height: %s"

	// BadHash256Digest occurs when a wrong-length hash256 digest is found
	BadHash256Digest sdk.CodeType = 105
	// BadHash256DigestMessage is the corresponding message
	BadHash256DigestMessage = "Digest %s had wrong length"

	// BadHex occurs when a hex argument couldn't be deserialized
	BadHex sdk.CodeType = 106
	// BadHexMessage is the corresponding message
	BadHexMessage = "Bad hex string in query or msg: %s"

	// BadHexLen occurs when a hex argument is the wrong length
	BadHexLen sdk.CodeType = 107
	// BadHexLenMessage is the corresponding message
	BadHexLenMessage = "Expected %d bytes in a RequestID, got %d"

	// BitcoinSPV is the code for errors bubbled up from Bitcoin SPV
	BitcoinSPV sdk.CodeType = 108
	// BitcoinSPVMessage is the corresponding message

	// AlreadyInit is the code for a second attempt to init the relay
	AlreadyInit sdk.CodeType = 109
	// AlreadyInitMessage is the corresponding message
	AlreadyInitMessage = "Relay has already set genesis state"

	// 200-block -- AddHeaders

	// UnexpectedRetarget indicates a retarget was seen during AddHeaders loop
	UnexpectedRetarget sdk.CodeType = 201
	// UnexpectedRetargetMessage is the corresponding message
	UnexpectedRetargetMessage = "Target changed unexpectedly at block %x"

	// 300-block AddHeadersWithRetarget

	// WrongEnd means the end block is at the wrong height
	WrongEnd sdk.CodeType = 301
	// WrongEndMessage is the corresponding message
	WrongEndMessage = "Must provide the last header of the closing difficulty period"

	// WrongStart means the start block is at the wrong height
	WrongStart sdk.CodeType = 302
	// WrongStartMessage is the corresponding message
	WrongStartMessage = "Must provide exactly 1 difficulty period"

	// PeriodMismatch means the start and end block do not have the same difficulty
	PeriodMismatch sdk.CodeType = 303
	// PeriodMismatchMessage is the corresponding message
	PeriodMismatchMessage = "Period header difficulties do not match"

	// BadRetarget means the provided blocks did not create the expected retarget
	BadRetarget sdk.CodeType = 304
	// BadRetargetMessage is the corresponding message
	BadRetargetMessage = "Invalid retarget provided"

	// 400-block -- MarkNewBestHeight

	// NotBestKnown means a block should have been the best known, but wasn't
	NotBestKnown sdk.CodeType = 403
	// NotBestKnownMessage is the corresponding message
	NotBestKnownMessage = "Provided digest %x is not current best known, expecting block with hash %x"

	// NotHeaviestAncestor means a later common ancestor was found
	NotHeaviestAncestor sdk.CodeType = 404
	// NotHeaviestAncestorMessage is the corresponding message
	NotHeaviestAncestorMessage = "Ancestor %x is not heaviest common ancestor"

	// NotHeavier means the proposed new best is not heavier than the current best
	NotHeavier sdk.CodeType = 405
	// NotHeavierMessage is the corresponding message
	NotHeavierMessage = "New best received %x does not have more work than previous best %x"

	// 500-block Queries

	// MarshalJSON means there was an error marshalling a query result to json
	MarshalJSON sdk.CodeType = 503
	// MarshalJSONMessage is the corresponding message
	MarshalJSONMessage = "Could not marshal result to JSON"

	// 600-block Proof Requests

	// UnknownRequest means the request was not found
	UnknownRequest sdk.CodeType = 601
	// UnknownRequestMessage is the corresponding message
	UnknownRequestMessage = "Request not found"

	// SpendsLength means the spend value is not 36 bytes
	SpendsLength sdk.CodeType = 602
	// SpendsLengthMessage is the corresponding message
	SpendsLengthMessage = "Spends value is not 36 bytes"

	// PaysLength means the pays value is greater than 50 bytes
	PaysLength sdk.CodeType = 603
	// PaysLengthMessage is the corresponding message
	PaysLengthMessage = "Pays value is greater than 50 bytes"

	// InvalidVin means the vin is not valid
	InvalidVin sdk.CodeType = 604
	// InvalidVinMessage is the corresponding message
	InvalidVinMessage = "Vin is not valid"

	// InvalidVout means the vout is not valid
	InvalidVout sdk.CodeType = 605
	// InvalidVoutMessage is the corresponding message
	InvalidVoutMessage = "Vout is not valid"

	// ClosedRequest means the request is not active
	ClosedRequest sdk.CodeType = 606
	// ClosedRequestMessage is the corresponding message
	ClosedRequestMessage = "Request is not active"

	// RequestPays means the output does not match the pays request
	RequestPays sdk.CodeType = 607
	// RequestPaysMessage is the corresponding message
	RequestPaysMessage = "Output not match pays for requestID %d"

	// RequestValue means the pays value and value of the output does not match
	RequestValue sdk.CodeType = 608
	// RequestValueMessage is the corresponding message
	RequestValueMessage = "Output value not match pays value for requestID %d"

	// RequestSpends means the request spends does not match the input
	RequestSpends sdk.CodeType = 609
	// RequestSpendsMessage is the corresponding message
	RequestSpendsMessage = "Input does not match spends for requestID %d"

	// NotAncestor means the LCA is not an ancestor of the SPV Proof header
	NotAncestor sdk.CodeType = 610
	// NotAncestorMessage is the corresponding message
	NotAncestorMessage = "LCA %x not ancestor of proof header"

	// NotEnoughConfs means the proof does not have enough confirmations
	NotEnoughConfs sdk.CodeType = 611
	// NotEnoughConfsMessage is the corresponding message
	NotEnoughConfsMessage = "Not enough confirmations for requestID %d"

	// ActionLength means the pays value is greater than 50 bytes
	ActionLength sdk.CodeType = 612
	// ActionLengthMessage is the corresponding message
	ActionLengthMessage = "Action value is greater than 500 bytes"

	// 700-block External

	// ExternalError is an error from a dependency
	ExternalError sdk.CodeType = 701
)

// ErrBadHeaderLength throws an error
func ErrBadHeaderLength(codespace sdk.CodespaceType, label string, digest RawHeader, length int) sdk.Error {
	return sdk.NewError(codespace, BadHeaderLength, fmt.Sprint(BadHeaderLengthMessage, label, digest, length))
}

// ErrUnknownBlock throws an error
func ErrUnknownBlock(codespace sdk.CodespaceType, label string, digest Hash256Digest) sdk.Error {
	return sdk.NewError(codespace, UnknownBlock, fmt.Sprint(UnknownBlockMessage, label, digest))
}

// ErrBadHeight throws an error
func ErrBadHeight(codespace sdk.CodespaceType, details string) sdk.Error {
	message := fmt.Sprintf(BadHeightMessage, details)
	return sdk.NewError(codespace, BadHeight, message)
}

// ErrUnexpectedRetarget throws an error
func ErrUnexpectedRetarget(codespace sdk.CodespaceType, rawHeader RawHeader) sdk.Error {
	return sdk.NewError(codespace, UnexpectedRetarget, fmt.Sprint(UnexpectedRetargetMessage, rawHeader))
}

// ErrWrongEnd throws an error
func ErrWrongEnd(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, WrongEnd, WrongEndMessage)
}

// ErrWrongStart throws an error
func ErrWrongStart(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, WrongStart, WrongStartMessage)
}

// ErrPeriodMismatch throws an error
func ErrPeriodMismatch(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, PeriodMismatch, PeriodMismatchMessage)
}

// ErrBadRetarget throws an error
func ErrBadRetarget(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadRetarget, BadRetargetMessage)
}

// ErrNotBestKnown throws an error
func ErrNotBestKnown(codespace sdk.CodespaceType, invalidBest, expectedBest []byte) sdk.Error {
	return sdk.NewError(codespace, NotBestKnown, fmt.Sprintf(NotBestKnownMessage, invalidBest, expectedBest))
}

// ErrNotHeaviestAncestor throws an error
func ErrNotHeaviestAncestor(codespace sdk.CodespaceType, ancestor Hash256Digest) sdk.Error {
	return sdk.NewError(codespace, NotHeaviestAncestor, fmt.Sprintf(NotHeaviestAncestorMessage, ancestor))
}

// ErrNotHeavier throws an error
func ErrNotHeavier(codespace sdk.CodespaceType, newBest, prevBest []byte) sdk.Error {
	return sdk.NewError(codespace, NotHeavier, fmt.Sprintf(NotHeavierMessage, newBest, prevBest))
}

// ErrBadHash256Digest throws an error
func ErrBadHash256Digest(codespace sdk.CodespaceType, invalidDigest string) sdk.Error {
	return sdk.NewError(codespace, BadHash256Digest, fmt.Sprintf(BadHash256DigestMessage, invalidDigest))
}

// ErrBadHex throws an error
func ErrBadHex(codespace sdk.CodespaceType, invalidHex string) sdk.Error {
	return sdk.NewError(codespace, BadHex, fmt.Sprintf(BadHexMessage, invalidHex))
}

// ErrBadHexLen throws an error
func ErrBadHexLen(codespace sdk.CodespaceType, expected, actual int) sdk.Error {
	return sdk.NewError(codespace, BadHexLen, fmt.Sprintf(BadHexLenMessage, expected, actual))
}

// ErrAlreadyInit throws an error
func ErrAlreadyInit(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, AlreadyInit, AlreadyInitMessage)
}

// ErrMarshalJSON throws an error
func ErrMarshalJSON(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, MarshalJSON, MarshalJSONMessage)
}

// FromBTCSPVError converts a btcutils error into an sdk error
func FromBTCSPVError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, BitcoinSPV, err.Error())
}

// ErrUnknownRequest throws an error
func ErrUnknownRequest(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, UnknownRequest, UnknownRequestMessage)
}

// ErrSpendsLength throws an error
func ErrSpendsLength(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, SpendsLength, SpendsLengthMessage)
}

// ErrPaysLength throws an error
func ErrPaysLength(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, PaysLength, PaysLengthMessage)
}

// ErrActionLength throws an error
func ErrActionLength(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, ActionLength, ActionLengthMessage)
}

// ErrInvalidVin throws an error
func ErrInvalidVin(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, InvalidVin, InvalidVinMessage)
}

// ErrInvalidVout throws an error
func ErrInvalidVout(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, InvalidVout, InvalidVoutMessage)
}

// ErrClosedRequest throws an error
func ErrClosedRequest(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, ClosedRequest, ClosedRequestMessage)
}

// ErrRequestPays throws an error
func ErrRequestPays(codespace sdk.CodespaceType, requestID RequestID) sdk.Error {
	return sdk.NewError(codespace, RequestPays, fmt.Sprintf(RequestPaysMessage, requestID))
}

// ErrRequestValue throws an error
func ErrRequestValue(codespace sdk.CodespaceType, requestID RequestID) sdk.Error {
	return sdk.NewError(codespace, RequestValue, fmt.Sprintf(RequestValueMessage, requestID))
}

// ErrRequestSpends throws an error
func ErrRequestSpends(codespace sdk.CodespaceType, requestID RequestID) sdk.Error {
	return sdk.NewError(codespace, RequestSpends, fmt.Sprintf(RequestSpendsMessage, requestID))
}

// ErrNotAncestor throws an error
func ErrNotAncestor(codespace sdk.CodespaceType, lca Hash256Digest) sdk.Error {
	return sdk.NewError(codespace, NotAncestor, fmt.Sprintf(NotAncestorMessage, lca))
}

// ErrNotEnoughConfs throws an error
func ErrNotEnoughConfs(codespace sdk.CodespaceType, requestID RequestID) sdk.Error {
	return sdk.NewError(codespace, NotEnoughConfs, fmt.Sprintf(NotEnoughConfsMessage, requestID))
}

// ErrExternal converts any external error into an sdk error
func ErrExternal(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, ExternalError, err.Error())
}
