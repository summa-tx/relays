package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/bitcoin-spv/golang/btcspv"

	"github.com/summa-tx/relays/proto"
)

var _ proto.MsgServer;

// RouterKey is a name for the router
const RouterKey = ModuleName // this was defined in your key.go file

/***** IngestHeaderChain *****/

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

// ToProto populates a protobuf MsgIngestHeaderChain from a MsgIngestHeaderChain
func (msg *proto.MsgIngestHeaderChain) ToProto(m *MsgIngestHeaderChain) {
	msg.Signer = m.Signer
	msg.Headers = []byte{m.Headers}
}

// NewMsgIngestHeaderChain instantiates a MsgIngestHeaderChain
func NewMsgIngestHeaderChain(address sdk.AccAddress, headers []BitcoinHeader) MsgIngestHeaderChain {
	return MsgIngestHeaderChain{
		address,
		headers,
	}
}

// GetSigners gets signers
func (msg MsgIngestHeaderChain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// Type returns an identifier
func (msg MsgIngestHeaderChain) Type() string { return "ingest_header_chain" }

// ValidateBasic runs stateless validation
func (msg MsgIngestHeaderChain) ValidateBasic() sdk.Error {
	for i := range msg.Headers {
		valid, err := msg.Headers[i].Validate()
		if !valid || err != nil {
			return FromBTCSPVError(DefaultCodespace, err)
		}
	}
	return nil
}

// TODO: toProto?
// GetSignBytes returns the sighash for the message
func (msg MsgIngestHeaderChain) GetSignBytes() []byte {
	var proto proto.MsgIngestHeaderChain
	proto.ToProto(&msg)

	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(proto))
}

// Route returns the route key
func (msg MsgIngestHeaderChain) Route() string { return RouterKey }

/***** IngestDifficultyChange *****/

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

// ToProto populates a protobuf MsgIngestDifficultyChange from a MsgIngestDifficultyChange
func (msg *proto.MsgIngestDifficultyChange) ToProto(m *MsgIngestDifficultyChange) {
	msg.Signer = m.Signer
	msg.Start = []byte{m.Start}
	msg.Headers = []byte{m.Headers}
}

// NewMsgIngestDifficultyChange instantiates a MsgIngestDifficultyChange
func NewMsgIngestDifficultyChange(address sdk.AccAddress, start Hash256Digest, headers []BitcoinHeader) MsgIngestDifficultyChange {
	return MsgIngestDifficultyChange{
		address,
		start,
		headers,
	}
}

// GetSigners gets signers
func (msg MsgIngestDifficultyChange) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// Type returns an identifier
func (msg MsgIngestDifficultyChange) Type() string { return "ingest_difficulty_change" }

// ValidateBasic runs stateless validation
func (msg MsgIngestDifficultyChange) ValidateBasic() sdk.Error {
	for i := range msg.Headers {
		valid, err := msg.Headers[i].Validate()
		if !valid || err != nil {
			return FromBTCSPVError(DefaultCodespace, err)
		}
	}
	return nil
}

// GetSignBytes returns the sighash for the message
func (msg MsgIngestDifficultyChange) GetSignBytes() []byte {
	var proto proto.MsgIngestDifficultyChange
	proto.ToProto(&msg)

	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(proto))
}

// Route returns the route key
func (msg MsgIngestDifficultyChange) Route() string { return RouterKey }

/***** MarkNewHeaviest *****/

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

// ToProto populates a protobuf MsgMarkNewHeaviest from a MsgMarkNewHeaviest
func (msg *proto.MsgMarkNewHeaviest) ToProto(m *MsgMarkNewHeaviest) {
	msg.Signer = m.Signer
	msg.Ancestor = []byte{m.Ancestor}
	msg.CurrentBest = []byte{m.CurrentBest}
	msg.NewBest = []byte{m.NewBest}
	msg.Limit = m.Limit
}

// NewMsgMarkNewHeaviest instantiates a MsgMarkNewHeaviest
func NewMsgMarkNewHeaviest(address sdk.AccAddress, ancestor Hash256Digest, currentBest RawHeader, newBest RawHeader, limit uint32) MsgMarkNewHeaviest {
	if limit == 0 {
		limit = DefaultLookupLimit
	}

	return MsgMarkNewHeaviest{
		address,
		ancestor,
		currentBest,
		newBest,
		limit,
	}
}

// GetSigners gets signers
func (msg MsgMarkNewHeaviest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// Type returns an identifier
func (msg MsgMarkNewHeaviest) Type() string { return "mark_new_heaviest" }

// ValidateBasic runs stateless validation
func (msg MsgMarkNewHeaviest) ValidateBasic() sdk.Error {
	if len(msg.CurrentBest) != 80 {
		return ErrBadHeaderLength(DefaultCodespace, "currentBest", msg.CurrentBest, len(msg.CurrentBest))
	}

	if len(msg.NewBest) != 80 {
		return ErrBadHeaderLength(DefaultCodespace, "newBest", msg.NewBest, len(msg.NewBest))
	}

	return nil
}

// GetSignBytes returns the sighash for the message
func (msg MsgMarkNewHeaviest) GetSignBytes() []byte {
	var proto proto.MsgMarkNewHeaviest
	proto.ToProto(&msg)

	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(proto))
}

// Route returns the route key
func (msg MsgMarkNewHeaviest) Route() string { return RouterKey }

/***** NewRequest *****/

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

// ToProto populates a protobuf MsgNewRequest from a MsgNewRequest
func (msg *proto.MsgNewRequest) ToProto(m *MsgNewRequest) {
	msg.Signer = m.Signer
	msg.Spends = []byte{m.Spends}
	msg.Pays = []byte{m.Pays}
	msg.PaysValue = m.PaysValue
	msg.NumConfs = uint32(m.NumConfs)
	msg.Origin = uint32(m.Origin)
	msg.Action = []byte{m.Action}
}

// GetSigners gets signers
func (msg MsgNewRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// Type returns an identifier
func (msg MsgNewRequest) Type() string { return "new_request" }

// ValidateBasic runs stateless validation
func (msg MsgNewRequest) ValidateBasic() sdk.Error {
	// TODO: validate output types
	if len(msg.Spends) != 36 && len(msg.Spends) != 0 {
		return ErrSpendsLength(DefaultCodespace)
	}
	if len(msg.Pays) > 50 {
		return ErrPaysLength(DefaultCodespace)
	}
	if len(msg.Action) > 500 {
		return ErrActionLength(DefaultCodespace)
	}
	return nil
}

// GetSignBytes returns the sighash for the message
func (msg MsgNewRequest) GetSignBytes() []byte {
	var proto proto.MsgNewRequest
	proto.ToProto(&msg)

	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(proto))
}

// Route returns the route key
func (msg MsgNewRequest) Route() string { return RouterKey }

/***** ProvideProof *****/

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

	filled, err := filledRequestsFromProto(m.Filled)
	if err != nil {
		return err
	}

	msg.Signer = address
	msg.Filled = filled

	return nil
}

// ToProto populates a protobuf MsgProvideProof from a MsgProvideProof
func (msg *proto.MsgProvideProof) ToProto(m *MsgProvideProof) {
	msg.Signer = m.Signer
	msg.Filled = filledRequestsToProto(m.Filled)
}

// NewMsgProvideProof instantiates a MsgProvideProof
func NewMsgProvideProof(address sdk.AccAddress, filledRequests FilledRequests) MsgProvideProof {
	return MsgProvideProof{
		address,
		filledRequests,
	}
}

// GetSigners gets signers
func (msg MsgProvideProof) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// ValidateBasic runs stateless validation
func (msg MsgProvideProof) ValidateBasic() sdk.Error {
	valid, err := msg.Filled.Proof.Validate()
	if !valid || err != nil {
		return FromBTCSPVError(DefaultCodespace, err)
	}

	return nil
}

// Type returns an identifier
func (msg MsgProvideProof) Type() string { return "provide_proof" }

// GetSignBytes returns the sighash for the message
func (msg MsgProvideProof) GetSignBytes() []byte {
	var proto proto.MsgProvideProof
	proto.ToProto(&msg)

	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(proto))
}

// Route returns the route key
func (msg MsgProvideProof) Route() string { return RouterKey }
