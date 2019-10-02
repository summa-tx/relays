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
	// NewMsgSetLink is what is says on the tin
	NewMsgSetLink = types.NewMsgSetLink
	// RegisterCodec is what is says on the tin
	RegisterCodec = types.RegisterCodec
	// ModuleCdc is what is says on the tin
	ModuleCdc = types.ModuleCdc
)

type (
	// Keeper is what is says on the tin
	Keeper = keeper.Keeper

	// Hash256Digest 32-byte double-sha2 digest
	Hash256Digest = types.Hash256Digest

	// MsgSetName is what is says on the tin
	MsgSetName = types.MsgSetLink
	// QueryResGetParent is what is says on the tin
	QueryResGetParent = types.QueryResGetParent
)
