package types

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

// event Extension(bytes32 indexed _first, bytes32 indexed _last);
// event Reorg(bytes32 indexed _from, bytes32 indexed _to, bytes32 indexed _gcd);
