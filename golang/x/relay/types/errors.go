package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultCodespace is the default code space
	DefaultCodespace sdk.CodespaceType = ModuleName

	// 100-block -- shared errors

	// UnknownError is an unknown error
	UnknownError sdk.CodeType = 100
	// UnknownErrorMessage is the corresponding message
	UnknownErrorMessage = "Unknown Error"

	// BadHeaderLength means Header array length not divisible by 80
	BadHeaderLength sdk.CodeType = 101
	// BadHeaderLengthMessage is the corresponding message
	BadHeaderLengthMessage = "Header array length must be divisible by 80"

	// InsufficientWork means a block did not have sufficient work
	InsufficientWork sdk.CodeType = 102
	// InsufficientWorkMessage is the corresponding message
	InsufficientWorkMessage = "Header work is insufficient"

	// UnknownBlock is the error code for unknown blocks
	UnknownBlock sdk.CodeType = 103
	// UnknownBlockMessage is the corresponding message
	UnknownBlockMessage = "Unknown block"

	// BadHeight occurs when a proposed descendant is below a proposed ancestor
	BadHeight sdk.CodeType = 104
	// BadHeightMessage is the corresponding message
	BadHeightMessage = "A descendant height is below the ancestor height"

	// BadHash256Digest occurs when a wrong-length hash256 digest is found
	BadHash256Digest sdk.CodeType = 105
	// BadHash256DigestMessage is the corresponding message
	BadHash256DigestMessage = "Digest had wrong length"

	// BadHex occurs when a hex argument couldn't be deserialized
	BadHex sdk.CodeType = 106
	// BadHexMessage is the corresponding message
	BadHexMessage = "Bad hex string in query or msg"

	// BitcoinSPV is the code for errors bubbled up from Bitcoin SPV
	BitcoinSPV sdk.CodeType = 107
	// BitcoinSPVMessage is the corresponding message

	// AlreadyInit is the code for a second attempt to init the relay
	AlreadyInit sdk.CodeType = 108
	// AlreadyInitMessage is the corresponding message
	AlreadyInitMessage = "Relay has already set genesis state"

	// 200-block -- AddHeaders

	// UnexpectedRetarget indicates a retarget was seen during AddHeaders loop
	UnexpectedRetarget sdk.CodeType = 201
	// UnexpectedRetargetMessage is the corresponding message
	UnexpectedRetargetMessage = "Target changed unexpectedly"

	// BadLink indicates a broken link was found in the header array during IngestHeaders loop
	BadLink sdk.CodeType = 202
	// BadLinkMessage is the corresponding message
	BadLinkMessage = "Headers do not form a consistent chain"

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
	NotBestKnownMessage = "Provided digest is not current best known"

	// NotHeaviestAncestor means a later common ancestor was found
	NotHeaviestAncestor sdk.CodeType = 404
	// NotHeaviestAncestorMessage is the corresponding message
	NotHeaviestAncestorMessage = "Ancestor must be heaviest common ancestor"

	// NotHeavier means the proposed new best is not heavier than the current best
	NotHeavier sdk.CodeType = 405
	// NotHeavierMessage is the corresponding message
	NotHeavierMessage = "New best hash does not have more work than previous"

	// 500-block Queries

	// TODO: Delete the next 2 error codes
	// NotEnoughArguments means there are not enough arguments specified in the path of a query
	NotEnoughArguments sdk.CodeType = 501
	// NotEnoughArgumentsMessage is the corresponding message
	NotEnoughArgumentsMessage = "Not enough arguments"

	// TooManyArguments means there are too many arguments specified in the path of a query
	TooManyArguments sdk.CodeType = 502
	// TooManyArgumentsMessage is the corresponding message
	TooManyArgumentsMessage = "Too many arguments"

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
	RequestPaysMessage = "Does not match pays request"

	// RequestValue means the pays value and value of the output does not match
	RequestValue sdk.CodeType = 608
	// RequestValueMessage is the corresponding message
	RequestValueMessage = "Does not match value request"

	// RequestSpends means the request spends does not match the input
	RequestSpends sdk.CodeType = 609
	// RequestSpendsMessage is the corresponding message
	RequestSpendsMessage = "Does not match spends request"

	// 700-block External

	ExternalError sdk.CodeType = 701
)

// ErrUnknownError throws an error
func ErrUnknownError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, UnknownError, UnknownErrorMessage)
}

// ErrBadHeaderLength throws an error
func ErrBadHeaderLength(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadHeaderLength, BadHeaderLengthMessage)
}

// ErrInsufficientWork throws an error
func ErrInsufficientWork(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, InsufficientWork, InsufficientWorkMessage)
}

// ErrUnknownBlock throws an error
func ErrUnknownBlock(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, UnknownBlock, UnknownBlockMessage)
}

// ErrBadHeight throws an error
func ErrBadHeight(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadHeight, BadHeightMessage)
}

// ErrUnexpectedRetarget throws an error
func ErrUnexpectedRetarget(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, UnexpectedRetarget, UnexpectedRetargetMessage)
}

// ErrBadLink throws an error
func ErrBadLink(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadLink, BadLinkMessage)
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
func ErrNotBestKnown(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, NotBestKnown, NotBestKnownMessage)
}

// ErrNotHeaviestAncestor throws an error
func ErrNotHeaviestAncestor(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, NotHeaviestAncestor, NotHeaviestAncestorMessage)
}

// ErrNotHeavier throws an error
func ErrNotHeavier(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, NotHeavier, NotHeavierMessage)
}

// ErrBadHash256Digest throws an error
func ErrBadHash256Digest(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadHash256Digest, BadHash256DigestMessage)
}

// ErrBadHex throws an error
func ErrBadHex(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadHex, BadHexMessage)
}

// ErrAlreadyInit throws an error
func ErrAlreadyInit(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, AlreadyInit, AlreadyInitMessage)
}

// ErrTooManyArguments throws an error
func ErrTooManyArguments(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, TooManyArguments, TooManyArgumentsMessage)
}

// ErrNotEnoughArguments throws an error
func ErrNotEnoughArguments(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, NotEnoughArguments, NotEnoughArgumentsMessage)
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
func ErrRequestPays(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, RequestPays, RequestPaysMessage)
}

// ErrRequestValue throws an error
func ErrRequestValue(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, RequestValue, RequestValueMessage)
}

// ErrRequestSpends throws an error
func ErrRequestSpends(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, RequestSpends, RequestSpendsMessage)
}

// ErrExternal converts any external error into an sdk error
func ErrExternal(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, ExternalError, err.Error())
}
