package relay

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the genesis state
type GenesisState struct {
	Headers     []BitcoinHeader `json:"headers"`
	PeriodStart BitcoinHeader   `json:"periodStart"`
}

// NewGenesisState instantiates a genesis state
func NewGenesisState(headers []BitcoinHeader, periodStart BitcoinHeader) GenesisState {
	return GenesisState{Headers: headers, PeriodStart: periodStart}
}

// ValidateGenesis validates a genesis state
func ValidateGenesis(data GenesisState) error {
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
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	err := keeper.SetGenesisState(ctx, data.Headers[0], data.PeriodStart)
	if err != nil {
		panic("already init!")
	}
	err = keeper.IngestHeaderChain(ctx, data.Headers[1:])
	if err != nil {
		panic("Bad header chain in genesis state!")
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the genesis state
// TODO: export GenesisState
//       May need special store keys for it
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	panic("Not implemented")
}
