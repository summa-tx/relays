package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Relay module event types
const (
	EventTypeExtension = "extension"
	EventTypeReorg     = "reorg"

	AttributeKeyFirstBlock = "first_block"
	AttributeKeyLastBlock  = "last_block"

	AttributeKeyPreviousBest = "previous_best"
	AttributeKeyNewBest      = "new_best"
	AttributeKeyLatestCommon = "latest_common_ancestor"
)

// NewReorgEvent instantiates a reorg event
func NewReorgEvent(prev, new, lca BitcoinHeader) sdk.Event {
	return sdk.NewEvent(
		EventTypeReorg,
		sdk.NewAttribute(AttributeKeyPreviousBest, "0x"+hex.EncodeToString(prev.HashLE[:])),
		sdk.NewAttribute(AttributeKeyNewBest, "0x"+hex.EncodeToString(new.HashLE[:])),
		sdk.NewAttribute(AttributeKeyLatestCommon, "0x"+hex.EncodeToString(lca.HashLE[:])),
	)
}

// NewExtensionEvent instantiates an extension event
func NewExtensionEvent(first, last BitcoinHeader) sdk.Event {
	return sdk.NewEvent(
		EventTypeExtension,
		sdk.NewAttribute(AttributeKeyFirstBlock, "0x"+hex.EncodeToString(first.HashLE[:])),
		sdk.NewAttribute(AttributeKeyLastBlock, "0x"+hex.EncodeToString(last.HashLE[:])),
	)
}
