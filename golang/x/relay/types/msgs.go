package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is a name for the router
const RouterKey = ModuleName // this was defined in your key.go file

/***** IngestHeaderChain *****/

// MsgIngestHeaderChain defines a IngestHeaderChain message
type MsgIngestHeaderChain struct {
	Signer  sdk.AccAddress  `json:"signer"`
	Headers []BitcoinHeader `json:"headers"`
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

// GetSignBytes returns the sighash for the message
func (msg MsgIngestHeaderChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route key
func (msg MsgIngestHeaderChain) Route() string { return RouterKey }

/***** IngestDifficultyChange *****/

// MsgIngestDifficultyChange defines a IngestDifficultyChange message
type MsgIngestDifficultyChange struct {
	Signer  sdk.AccAddress  `json:"signer"`
	Start   Hash256Digest   `json:"prevEpochStartLE"`
	Headers []BitcoinHeader `json:"headers"`
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route key
func (msg MsgIngestDifficultyChange) Route() string { return RouterKey }

/***** MarkNewHeaviest *****/

// MsgMarkNewHeaviest defines a MarkNewHeaviest message
type MsgMarkNewHeaviest struct {
	Signer      sdk.AccAddress `json:"signer"`
	Ancestor    Hash256Digest  `json:"ancestor"`
	CurrentBest RawHeader      `json:"currentBest"`
	NewBest     RawHeader      `json:"newBest"`
	Limit       uint32         `json:"limit"`
}

// NewMsgMarkNewHeaviest instantiates a MsgMarkNewHeaviest
func NewMsgMarkNewHeaviest(address sdk.AccAddress, ancestor Hash256Digest, currentBest RawHeader, newBest RawHeader, limit uint32) MsgMarkNewHeaviest {
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
	if len(msg.CurrentBest) != 80 || len(msg.NewBest) != 80 {
		return ErrBadHeaderLength(DefaultCodespace)
	}

	return nil
}

// GetSignBytes returns the sighash for the message
func (msg MsgMarkNewHeaviest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route key
func (msg MsgMarkNewHeaviest) Route() string { return RouterKey }

/***** NewRequest *****/

// MsgNewRequest defines a NewRequest message
type MsgNewRequest struct {
	Signer    sdk.AccAddress `json:"signer"`
	Spends    []byte         `json:"spends"`
	Pays      []byte         `json:"pays"`
	PaysValue uint64         `json:"paysValue"`
	NumConfs  uint8          `json:"numConfs"`
}

// NewMsgNewRequest instantiates a MsgNewRequest
func NewMsgNewRequest(address sdk.AccAddress, spends, pays []byte, paysValue uint64, numConfs uint8) MsgNewRequest {
	return MsgNewRequest{
		address,
		spends,
		pays,
		paysValue,
		numConfs,
	}
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
	if len(msg.Spends) != 36 {
		return ErrSpendsLength(DefaultCodespace)
	}
	if len(msg.Pays) > 50 {
		return ErrPaysLength(DefaultCodespace)
	}
	return nil
}

// GetSignBytes returns the sighash for the message
func (msg MsgNewRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route key
func (msg MsgNewRequest) Route() string { return RouterKey }

/***** ProvideProof *****/

// MsgProvideProof defines a NewRequest message
type MsgProvideProof struct {
	Signer sdk.AccAddress `json:"signer"`
	Filled FilledRequests `json:"filled"`
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route key
func (msg MsgProvideProof) Route() string { return RouterKey }
