package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultCodespace is the default code space
	DefaultCodespace sdk.CodespaceType = ModuleName

	// 100-block -- shared errors

	// BadHeaderLength means Header array length not divisible by 80
	BadHeaderLength sdk.CodeType = 101

	// InsufficientWork means a block did not have sufficient work
	InsufficientWork sdk.CodeType = 102

	// UnknownBlock is the error code for unknown blocks
	UnknownBlock sdk.CodeType = 103

	// BadHeight occurs when a proposed descendant is below a proposed ancestor
	BadHeight sdk.CodeType = 104

	// 200-block -- AddHeaders

	// RetargetOnExternal indicates a retarget was seen in AddHeaders
	RetargetOnExternal sdk.CodeType = 201

	// UnexptectedRetarget indicates a retarget was seen during AddHeaders loop
	UnexptectedRetarget sdk.CodeType = 202

	// BadLink indicates a broken link was found in the header array during AddHeaders loop
	BadLink sdk.CodeType = 203

	// 300-block AddHeadersWithRetarget

	// WrongEnd means the end block is at the wrong height
	WrongEnd sdk.CodeType = 301

	// WrongStart means the start block is at the wrong height
	WrongStart sdk.CodeType = 302

	// PeriodMismatch means the start and end block do not have the same difficulty
	PeriodMismatch sdk.CodeType = 303

	// BadRetarget means the provided blocks did not create the expected retarget
	BadRetarget sdk.CodeType = 304

	// 400-block -- MarkNewBestHeight

	// NotBestKnown means a block should have been the best known, but wasn't
	NotBestKnown sdk.CodeType = 403

	// NotHeaviestAncestor means a later common ancestor was found
	NotHeaviestAncestor sdk.CodeType = 404

	// NotHeavier means the proposed new best is not heavier than the current best
	NotHeavier sdk.CodeType = 405
)

// ErrBadHeaderLength throws an error
func ErrBadHeaderLength(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadHeaderLength, "Header array length must be divisible by 80")
}

// ErrInsufficientWork throws an error
func ErrInsufficientWork(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, InsufficientWork, "Header work is insufficient")
}

// ErrUnknownBlock throws an error
func ErrUnknownBlock(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, UnknownBlock, "Unknown block")
}

// ErrBadHeight throws an error
func ErrBadHeight(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadHeight, "A descendant height is below the ancestor height")
}

// ErrRetargetOnExternal throws an error
func ErrRetargetOnExternal(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, RetargetOnExternal, "Unexpected retarget on external call")
}

// ErrUnexptectedRetarget throws an error
func ErrUnexptectedRetarget(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, UnexptectedRetarget, "Target changed unexpectedly")
}

// ErrBadLink throws an error
func ErrBadLink(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadLink, "Headers do not form a consistent chain")
}

// ErrWrongEnd throws an error
func ErrWrongEnd(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, WrongEnd, "Must provide the last header of the closing difficulty period")
}

// ErrWrongStart throws an error
func ErrWrongStart(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, WrongStart, "Must provide exactly 1 difficulty period")
}

// ErrPeriodMismatch throws an error
func ErrPeriodMismatch(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, PeriodMismatch, "Period header difficulties do not match")
}

// ErrBadRetarget throws an error
func ErrBadRetarget(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, BadRetarget, "Invalid retarget provided")
}

// ErrNotBestKnown throws an error
func ErrNotBestKnown(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, NotBestKnown, "New best is unknown")
}

// ErrNotHeaviestAncestor throws an error
func ErrNotHeaviestAncestor(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, NotHeaviestAncestor, "Ancestor must be heaviest common ancestor")
}

// ErrNotHeavier throws an error
func ErrNotHeavier(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, NotHeavier, "New best hash does not have more work than previous")
}
