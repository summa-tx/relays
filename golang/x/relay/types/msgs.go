package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"

	"github.com/summa-tx/relays/proto"
)

var _ proto.MsgServer;

// MsgIngestHeaderChain defines a IngestHeaderChain message
type MsgIngestHeaderChain struct {
	Signer  sdk.AccAddress
	Headers []BitcoinHeader
}

// FromProto populates a MsgIngestHeaderChain from a protobuf
func (msg *MsgIngestHeaderChain) FromProto(m *proto.MsgIngestHeaderChain) (error) {
	// Do any parsing/translation work
	address, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return err
	}

	headers, err := headerSliceFromProto(m.Headers)
	if err != nil {
		return err
	}

	msg.Signer = address
	msg.Headers = headers

	return nil
}

// MsgIngestDifficultyChange defines a IngestDifficultChange message
type MsgIngestDifficultyChange struct {
	Signer  sdk.AccAddress
	Start   btcspv.Hash256Digest
	Headers []BitcoinHeader
}

// FromProto populates a MsgIngestDifficultyChange from a protobuf
func (msg *MsgIngestDifficultyChange) FromProto(m *proto.MsgIngestDifficultyChange) (error) {

	address, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return err
	}

	start, err := bufToH256(m.Start)
	if err != nil {
		return err
	}

	headers, err := headerSliceFromProto(m.Headers)
	if err != nil {
		return err
	}

	msg.Signer = address
	msg.Start = start
	msg.Headers = headers

	return nil
}

// MsgMarkNewHeaviest defines a MarkNewHeaviest message
type MsgMarkNewHeaviest struct {
	Signer      sdk.AccAddress
	Ancestor    btcspv.Hash256Digest
	CurrentBest btcspv.RawHeader
	NewBest     btcspv.RawHeader
	Limit       uint32
}

// FromProto populates a MsgMarkNewHeaviest from a protobuf
func (msg *MsgMarkNewHeaviest) FromProto(m *proto.MsgMarkNewHeaviest) (error) {
	address, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return err
	}

	ancestor, err := bufToH256(m.Ancestor)
	if err != nil {
		return err
	}

	currentBest, err := bufToRawHeader(m.CurrentBest)
	if err != nil {
		return err
	}

	newBest, err := bufToRawHeader(m.NewBest)
	if err != nil {
		return err
	}

	msg.Signer = address
	msg.Ancestor = ancestor
	msg.CurrentBest = currentBest
	msg.NewBest = newBest
	msg.Limit = m.Limit

	return nil
}

// MsgNewRequest defines a NewRequest message
type MsgNewRequest struct {
	Signer    sdk.AccAddress
	Spends    HexBytes
	Pays      HexBytes
	PaysValue uint64
	NumConfs  uint8
	Origin    Origin
	Action    HexBytes
}


// FromProto populates a MsgNewRequest from a protobuf
func (msg *MsgNewRequest) FromProto(m *proto.MsgNewRequest) (error) {
	address, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return err
	}

	msg.Signer = address
	msg.Spends = btcspv.HexBytes(m.Spends)
	msg.Pays = btcspv.HexBytes(m.Pays)
	msg.PaysValue = m.PaysValue
	msg.NumConfs = uint8(m.NumConfs)
	msg.Origin = Origin(uint8(m.Origin))
	msg.Action = btcspv.HexBytes(m.Action)

	return nil
}

// MsgProvideProof defines a NewRequest message
type MsgProvideProof struct {
	Signer sdk.AccAddress
	Filled FilledRequests
}

// FromProto populates a MsgProvideProof from a protobuf
func (msg *MsgProvideProof) FromProto(m *proto.MsgProvideProof) (error) {

	address, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return err
	}

	var filled FilledRequests;
	err = filled.FromProto(m.Filled)
	if err != nil {
		return err
	}

	msg.Signer = address
	msg.Filled = filled

	return nil
}
