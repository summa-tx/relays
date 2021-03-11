package types

import (
	"fmt"

	"github.com/summa-tx/bitcoin-spv/golang/btcspv"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/summa-tx/relays/golang/x/relay/types"

	"github.com/summa-tx/realys/proto"
)

var _ MsgServer;

func stringToError(str string) (sdk.Error) {
	// TODO: How to handle CodeType?
	return sdk.NewError(DefaultCodespace, ExternalError, str)
}

func bufToH256(buf []byte) (btcspv.Hash256Digest, error) {
	var h btcspv.Hash256Digest;
	if len(buf) != 32 {
		return h, fmt.Errorf("Expected 32 bytes, got %d bytes", len(buf))
	}

	copy(h[:], buf)

	return h, nil
}

func bufToRawHeader(buf []byte) (btcspv.RawHeader, error) {
	var h btcspv.RawHeader;
	if len(buf) != 80 {
		return h, fmt.Errorf("Expected 80 bytes, got %d bytes", len(buf))
	}

	copy(h[:], buf)

	return h, nil
}

func (m *MsgIngestHeaderChain) translate() (types.MsgIngestHeaderChain, error) {
	var msg types.MsgIngestHeaderChain;

	// Do any parsing/translation work
	address, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return msg, err
	}

	msg.Signer = address
	msg.Headers = m.Headers

	return msg, nil
}

func (m *MsgIngestDifficultyChange) translate() (types.MsgIngestDifficultyChange, error) {
	var msg types.MsgIngestDifficultyChange;

	address, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return msg, err
	}

	start, err := bufToH256(q.Start)
	if err != nil {
		return msg, err
	}

	msg.Signer = address
	msg.Start = start
	msg.Headers = m.Headers

	return msg, nil
}

func (m *MsgMarkNewHeaviest) translate() (types.MsgMarkNewHeaviest, error) {
	var msg types.MsgMarkNewHeaviest;

	address, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return msg, err
	}

	ancestor, err := bufToH256(q.Ancestor)
	if err != nil {
		return msg, err
	}

	currentBest, err := bufToRawHeader(q.CurrentBest)
	if err != nil {
		return msg, err
	}

	newBest, err := bufToRawHeader(q.NewBest)
	if err != nil {
		return msg, err
	}

	msg.Signer = address
	msg.Ancestor = m.Ancestor
	msg.CurrentBest = currentBest
	msg.NewBest = newBest
	msg.Limit = m.Limit

	return msg, nil
}

func (m *MsgNewRequest) translate() (types.MsgNewRequest, error) {
	var msg types.MsgNewRequest;

	address, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return msg, err
	}

	msg.Signer = address
	msg.Spends = btcspv.HexBytes(m.Spends)
	msg.Pays = btcspv.HexBytes(m.Pays)
	msg.PaysValue = m.PaysValue
	msg.NumConfs = uint8(m.NumConfs)
	msg.Origin = m.Origin
	msg.Action = btcspv.HexBytes(m.Action)

	return msg, nil
}

func (m *MsgProvideProof) translate() (types.MsgProvideProof, error) {
	var msg types.MsgProvideProof;

	address, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return msg, err
	}

	msg.Signer = address
	msg.FilledRequests = m.FilledRequests

	return msg, nil
}
