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

// DefaultGenesisState makes a default empty genesis state
// TODO: set recent block as default
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Headers:     []BitcoinHeader{},
		PeriodStart: BitcoinHeader{},
	}
}

// InitGenesis inits the app state based on the genesis state
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) ([]abci.ValidatorUpdate, sdk.Error) {
	err := keeper.SetGenesisState(ctx, data.Headers[0], data.PeriodStart)
	if err != nil {
		return []abci.ValidatorUpdate{}, err
	}
	err = keeper.IngestHeaderChain(ctx, data.Headers[1:])
	if err != nil {
		return []abci.ValidatorUpdate{}, err
	}
	return []abci.ValidatorUpdate{}, nil
}

// ExportGenesis exports the genesis state
// TODO: export GenesisState
//       May need special store keys for it
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	panic("Not implemented")
}
