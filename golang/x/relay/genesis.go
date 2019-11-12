package relay

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the genesis state
type GenesisState struct {
	Headers []btcspv.BitcoinHeader `json:"headers"`
}

// NewGenesisState instantiates a genesis state
func NewGenesisState(headers []btcspv.BitcoinHeader) GenesisState {
	return GenesisState{Headers: headers}
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

	return nil
}

// DefaultGenesisState makes a default empty genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Headers: []btcspv.BitcoinHeader{},
	}
}

// InitGenesis inits the app state based on the genesis state
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.IngestHeaderChain(ctx, data.Headers)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	panic("Not implemented")
}
