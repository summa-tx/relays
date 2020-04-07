package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Relay module event types
const (
	EventTypeExtension     = "extension"
	EventTypeReorg         = "reorg"
	EventTypeProofRequest  = "proof_request"
	EventTypeProofProvided = "proof_provided"

	AttributeKeyFirstBlock = "first_block"
	AttributeKeyLastBlock  = "last_block"

	AttributeKeyPreviousBest = "previous_best"
	AttributeKeyNewBest      = "new_best"
	AttributeKeyLatestCommon = "latest_common_ancestor"

	AttributeKeyRequestID = "request_id"
	AttributeKeyPays      = "pays"
	AttributeKeySpends    = "spends"
	AttributeKeyPaysValue = "value"
	AttributeKeyOrigin    = "origin"

	AttributeKeyTXID   = "txid"
	AttributeKeyFilled = "filled"
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
		sdk.NewAttribute(AttributeKeyFirstBlock, "0x"+hex.EncodeToString(first.Hash[:])),
		sdk.NewAttribute(AttributeKeyLastBlock, "0x"+hex.EncodeToString(last.Hash[:])),
	)
}

// NewProofRequestEvent instantiates a proof request event
func NewProofRequestEvent(pays, spends []byte, paysValue uint64, id RequestID, origin Origin) sdk.Event {
	return sdk.NewEvent(
		EventTypeProofRequest,
		sdk.NewAttribute(AttributeKeyRequestID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(AttributeKeyPays, "0x"+hex.EncodeToString(pays[:])),
		sdk.NewAttribute(AttributeKeySpends, "0x"+hex.EncodeToString(spends[:])),
		sdk.NewAttribute(AttributeKeyPaysValue, fmt.Sprintf("%d", paysValue)),
		sdk.NewAttribute(AttributeKeyOrigin, fmt.Sprintf("%d", origin)),
	)
}

// NewProofProvidedEvent instantiates a proof provided event
func NewProofProvidedEvent(txid Hash256Digest, filled []RequestID) sdk.Event {
	filledJSON, _ := json.Marshal(filled)
	return sdk.NewEvent(
		EventTypeProofProvided,
		sdk.NewAttribute(AttributeKeyTXID, "0x"+hex.EncodeToString(txid[:])),
		sdk.NewAttribute(AttributeKeyFilled, string(filledJSON)),
	)
}
