package types

const (
	// ModuleName is the name of the module
	ModuleName = "relay"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// LinkStorePrefix to be used when accessing links
	LinkStorePrefix = ModuleName + "links-"

	// HeaderStorePrefix to be used when accessing headers
	HeaderStorePrefix = ModuleName + "headers-"

	// RelayGenesisStorage is the storage key for the relay genesis digest
	RelayGenesisStorage = "RelayGenesis"

	// BestKnownDigestStorage is the storage key for the best known digest
	BestKnownDigestStorage = "BestKnownDigest"
)
