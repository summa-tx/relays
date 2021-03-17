package types

import (
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	// "github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry ctypes.InterfaceRegistry) {
	// registry.RegisterImplementations((*sdk.Msg)(nil),
	// 	&MsgIngestHeaderChain{},
	// 	&MsgIngestDifficultyChange{},
	// 	&MsgMarkNewHeaviest{},
	// 	&MsgNewRequest{},
	// 	&MsgProvideProof{},
	// )

	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	sdk.RegisterInterfaces(registry)
	txtypes.RegisterInterfaces(registry)
	// cryptocodec.RegisterInterfaces(registry)
}
