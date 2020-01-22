package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Relay module event types
const (
	EventTypeExtension    = "extension"
	EventTypeReorg        = "reorg"
	EventTypeProofRequest = "proof_request"

	AttributeKeyFirstBlock = "first_block"
	AttributeKeyLastBlock  = "last_block"

	AttributeKeyPreviousBest = "previous_best"
	AttributeKeyNewBest      = "new_best"
	AttributeKeyLatestCommon = "latest_common_ancestor"

	AttributeKeyRequestID = "request_id"
	AttributeKeyPays      = "pays"
	AttributeKeySpends    = "spends"
	AttributeKeyPaysValue = "value"
)

// NewReorgEvent instantiates a reorg event
func NewReorgEvent(prev, new, lca Hash256Digest) sdk.Event {
	return sdk.NewEvent(
		EventTypeReorg,
		sdk.NewAttribute(AttributeKeyPreviousBest, "0x"+hex.EncodeToString(prev[:])),
		sdk.NewAttribute(AttributeKeyNewBest, "0x"+hex.EncodeToString(new[:])),
		sdk.NewAttribute(AttributeKeyLatestCommon, "0x"+hex.EncodeToString(lca[:])),
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

// NewExtensionEvent instantiates a proof request event
func NewProofRequestEvent(pays, spends []byte, paysValue uint64, id RequestID) sdk.Event {
	return sdk.NewEvent(
		EventTypeProofRequest,
		sdk.NewAttribute(AttributeKeyRequestID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(AttributeKeyPays, "0x"+hex.EncodeToString(pays[:])),
		sdk.NewAttribute(AttributeKeySpends, "0x"+hex.EncodeToString(spends[:])),
		sdk.NewAttribute(AttributeKeyPaysValue, fmt.Sprintf("%d", paysValue)),
	)
}
