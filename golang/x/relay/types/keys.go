package types

const (
	// ModuleName is the name of the module
	ModuleName = "relay"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// LinkStorePrefix to be used when accessing links
	LinkStorePrefix = ModuleName + "-links-"

	// HeaderStorePrefix to be used when accessing headers
	HeaderStorePrefix = ModuleName + "-headers-"

	// RequestStorePrefix to be used when making requests
	RequestStorePrefix = ModuleName + "-requests-"

	// ChainStorePrefix to be used when accessing chain metadata
	ChainStorePrefix = ModuleName + "-chain-"

	// RelayGenesisStorage is the storage key for the relay genesis digest
	RelayGenesisStorage = "RelayGenesis"

	// BestKnownDigestStorage is the storage key for the best known digest
	BestKnownDigestStorage = "BestKnownDigest"

	// LastReorgLCAStorage is the storage key for the last reorg LCA
	LastReorgLCAStorage = "LastReorgLCA"

	// currentEpochDiffStorage is the storage key for the current epoch difficulty
	CurrentEpochDiffStorage = "currentEpochDifficulty"

	// prevEpochDiffStorage is the storage key for the prev epoch difficulty
	PrevEpochDiffStorage = "prevEpochDifficulty"

	// RequestIDTag is the storage key for the next Request ID to be used
	// when storing a request
	RequestIDTag = "id"
)
