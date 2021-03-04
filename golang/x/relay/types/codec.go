package types

import (
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterInterfaces(registry ctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIngestHeaderChain{},
		&MsgIngestDifficultyChange{},
		&MsgMarkNewHeaviest{},
		&MsgNewRequest{},
		&MsgProvideProof{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
