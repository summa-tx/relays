package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is a name for the router
const RouterKey = ModuleName // this was defined in your key.go file

/***** SetLink *****/

// MsgSetLink defines a SetLink message
type MsgSetLink struct {
	Signer sdk.AccAddress `json:"signer"`
	Header string         `json:"header"`
}

// NewMsgSetLink is a constructor function for MsgSetLink
func NewMsgSetLink(address sdk.AccAddress, header string, owner sdk.AccAddress) MsgSetLink {
	return MsgSetLink{
		address,
		header,
	}
}

// GetSigners defines whose signature is required
func (msg MsgSetLink) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// Type should return the action
func (msg MsgSetLink) Type() string { return "set_link" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetLink) ValidateBasic() sdk.Error {
	b, err := hex.DecodeString(msg.Header)
	if err != nil || len(b) != 80 {
		return FromBTCSPVError(DefaultCodespace, err)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetLink) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route should return the name of the module
func (msg MsgSetLink) Route() string { return RouterKey }

/***** IngestHeaderChain *****/

// MsgIngestHeaderChain defines a IngestHeaderChain message
type MsgIngestHeaderChain struct {
	Signer  sdk.AccAddress  `json:"signer"`
	Headers []BitcoinHeader `json:"headers"`
}

func NewMsgIngestHeaderChain(address sdk.AccAddress, headers []BitcoinHeader) MsgIngestHeaderChain {
	return MsgIngestHeaderChain{
		address,
		headers,
	}
}

func (msg MsgIngestHeaderChain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

func (msg MsgIngestHeaderChain) Type() string { return "ingest_header_chain" }

func (msg MsgIngestHeaderChain) ValidateBasic() sdk.Error {
	for i := range msg.Headers {
		valid, err := msg.Headers[i].Validate()
		if valid == false || err != nil {
			return FromBTCSPVError(DefaultCodespace, err)
		}
	}
	return nil
}

func (msg MsgIngestHeaderChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIngestHeaderChain) Route() string { return RouterKey }

/***** IngestDifficultyChange *****/

// MsgIngestDifficultyChange defines a IngestDifficultyChange message
type MsgIngestDifficultyChange struct {
	Signer  sdk.AccAddress  `json:"signer"`
	Start   Hash256Digest   `json:"prevEpochStartLE"`
	Headers []BitcoinHeader `json:"headers"`
}

func NewMsgIngestDifficultyChange(address sdk.AccAddress, start Hash256Digest, headers []BitcoinHeader) MsgIngestDifficultyChange {
	return MsgIngestDifficultyChange{
		address,
		start,
		headers,
	}
}

func (msg MsgIngestDifficultyChange) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

func (msg MsgIngestDifficultyChange) Type() string { return "ingest_difficulty_change" }

func (msg MsgIngestDifficultyChange) ValidateBasic() sdk.Error {
	for i := range msg.Headers {
		valid, err := msg.Headers[i].Validate()
		if valid == false || err != nil {
			return ErrBadHeaderLength(DefaultCodespace)
		}
	}
	return nil
}

func (msg MsgIngestDifficultyChange) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

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

func NewMsgMarkNewHeaviest(address sdk.AccAddress, ancestor Hash256Digest, currentBest RawHeader, newBest RawHeader, limit uint32) MsgMarkNewHeaviest {
	return MsgMarkNewHeaviest{
		address,
		ancestor,
		currentBest,
		newBest,
		limit,
	}
}

func (msg MsgMarkNewHeaviest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

func (msg MsgMarkNewHeaviest) Type() string { return "mark_new_heaviest" }

func (msg MsgMarkNewHeaviest) ValidateBasic() sdk.Error {
	if len(msg.CurrentBest) != 80 || len(msg.NewBest) != 80 {
		return ErrBadHeaderLength(DefaultCodespace)
	}

	return nil
}

func (msg MsgMarkNewHeaviest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMarkNewHeaviest) Route() string { return RouterKey }
