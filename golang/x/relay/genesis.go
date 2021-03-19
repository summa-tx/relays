package relay

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	"github.com/summa-tx/relays/golang/x/relay/keeper"
	"github.com/summa-tx/relays/golang/x/relay/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the genesis state
type GenesisState struct {
	Headers     []types.BitcoinHeader `json:"headers"`
	PeriodStart types.BitcoinHeader   `json:"periodStart"`
}

// NewGenesisState instantiates a genesis state
func NewGenesisState(headers []types.BitcoinHeader, periodStart types.BitcoinHeader) GenesisState {
	return GenesisState{Headers: headers, PeriodStart: periodStart}
}

// ValidateGenesis validates a genesis state
func ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage, data GenesisState) error {
	raw := []byte{}
	for _, header := range data.Headers {
		_, err := header.Validate()
		if err != nil {
			return err
		}
		raw = append(raw, header.Raw[:]...)
	}

	_, err := btcspv.ValidateHeaderChain(raw)
	if err != nil {
		return err
	}

	// Genesis state must include first block of an epoch plus another block belonging to that same epoch
	if data.PeriodStart.Height != (data.Headers[0].Height - (data.Headers[0].Height % 2016)) {
		return errors.New("period start has incorrect height")
	}

	return nil
}

// DefaultGenesisState sets block 606210 as genesis
func DefaultGenesisState() GenesisState {
	periodStart, headers := getGenesisHeaders()
	return GenesisState{
		Headers:     headers,
		PeriodStart: periodStart,
	}
}

// InitGenesis inits the app state based on the genesis state
func InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessagekeeper keeper.Keeper, genData GenesisState) []abci.ValidatorUpdate {
	err := keeper.SetGenesisState(ctx, genData.Headers[0], genData.PeriodStart)
	if err != nil {
		panic("already init!")
	}
	if len(genData.Headers) > 1 {
		err = keeper.IngestHeaderChain(ctx, genData.Headers[1:])
		if err != nil {
			panic("Bad header chain in genesis state! " + err.Error())
		}
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the genesis state
// TODO: export GenesisState
//       May need special store keys for it
func ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage, k keeper.Keeper) GenesisState {
	panic("Not implemented")
}
