package relay

import (
	"github.com/summa-tx/relays/golang/x/relay/keeper"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

const (
	// ModuleName is what it says on the tin
	ModuleName = types.ModuleName
	// RouterKey is what it says on the tin
	RouterKey = types.RouterKey
	//StoreKey is what it says on the tin
	StoreKey = types.StoreKey
)

var (
	// NewKeeper is what is says on the tin
	NewKeeper = keeper.NewKeeper
	// NewQuerier is what is says on the tin
	NewQuerier = keeper.NewQuerier
	// NewMsgIngestHeaderChain is what is says on the tin
	NewMsgIngestHeaderChain = types.NewMsgIngestHeaderChain
	// NewMsgIngestDifficultyChange is what is says on the tin
	NewMsgIngestDifficultyChange = types.NewMsgIngestDifficultyChange
	// NewMsgMarkNewHeaviest is what is says on the tin
	NewMsgMarkNewHeaviest = types.NewMsgMarkNewHeaviest
	// NewMsgNewRequest is what is says on the tin
	NewMsgNewRequest = types.NewMsgNewRequest
	// NewMsgProvideProof is what is says on the tin
	NewMsgProvideProof = types.NewMsgProvideProof
	// RegisterCodec is what is says on the tin
	RegisterCodec = types.RegisterCodec
	// ModuleCdc is what is says on the tin
	ModuleCdc = types.ModuleCdc
)

type (
	// Keeper is what is says on the tin
	Keeper = keeper.Keeper

	// TODO: add query structs here

	// Hash256Digest 32-byte double-sha2 digest
	Hash256Digest = types.Hash256Digest

	// Hash160Digest is a 20-byte ripemd160+sha2 hash
	Hash160Digest = types.Hash160Digest

	// RawHeader is an 80-byte raw header
	RawHeader = types.RawHeader

	// HexBytes is a type alias to make JSON hex ser/deser easier
	HexBytes = types.HexBytes

	// BitcoinHeader is a parsed Bitcoin header
	BitcoinHeader = types.BitcoinHeader

	// SPVProof is the base struct for an SPV proof
	SPVProof = types.SPVProof

	// ProofHandler is an interface to which the keepers dispatches valid proofs
	ProofHandler = types.ProofHandler

	// NullHandler does nothing
	NullHandler = types.NullHandler
)
