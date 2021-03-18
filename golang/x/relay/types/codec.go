package types

// import (
// 	"github.com/cosmos/cosmos-sdk/codec"
// 	"github.com/cosmos/cosmos-sdk/codec/types"
// 	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/types/msgservice"
// )
import (
	"github.com/summa-tx/relays/proto"
	"github.com/cosmos/cosmos-sdk/codec"
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	// "github.com/cosmos/cosmos-sdk/types/msgservice"
)

// // RegisterLegacyAminoCodec registers the necessary x/bank interfaces and concrete types
// // on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
// func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
// 	cdc.RegisterConcrete(&MsgSend{}, "cosmos-sdk/MsgSend", nil)
// 	cdc.RegisterConcrete(&MsgMultiSend{}, "cosmos-sdk/MsgMultiSend", nil)
// }

// RegisterLegacyAminoCodec registers the necessary x/bank interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&proto.MsgIngestHeaderChain{}, "relays/MsgIngestHeaderChain", nil)
	cdc.RegisterConcrete(&proto.MsgIngestDifficultyChange{}, "relays/MsgIngestDifficultyChange", nil)
	cdc.RegisterConcrete(&proto.MsgMarkNewHeaviest{}, "relays/MsgMarkNewHeaviest", nil)
	cdc.RegisterConcrete(&proto.MsgNewRequest{}, "relays/MsgNewRequest", nil)
	cdc.RegisterConcrete(&proto.MsgProvideProof{}, "relays/MsgProvideProof", nil)
}

// func RegisterInterfaces(registry types.InterfaceRegistry) {
// 	registry.RegisterImplementations((*sdk.Msg)(nil),
// 		&MsgSend{},
// 		&MsgMultiSend{},
// 	)

// 	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
// }
func RegisterInterfaces(registry ctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&proto.MsgIngestHeaderChain{},
		&proto.MsgIngestDifficultyChange{},
		&proto.MsgMarkNewHeaviest{},
		&proto.MsgNewRequest{},
		&proto.MsgProvideProof{},
	)

	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	sdk.RegisterInterfaces(registry)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
}
