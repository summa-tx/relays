package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIngestHeaderChain{}, "relay/IngestHeaderChain", nil)
	cdc.RegisterConcrete(MsgIngestDifficultyChange{}, "relay/IngestDifficultyChange", nil)
	cdc.RegisterConcrete(MsgMarkNewHeaviest{}, "relay/MarkNewHeaviest", nil)

}
