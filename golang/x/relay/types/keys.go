package types

const (
	// ModuleName is the name of the module
	ModuleName = "relay"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RelayGenesisStorage is the storage key for the relay genesis digest
	RelayGenesisStorage = "RelayGenesis"

	// BestKnownDigestStorage is the storage key for the best known digest
	BestKnownDigestStorage = "BestKnownDigest"
)
