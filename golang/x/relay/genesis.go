package relay

import (
	"encoding/hex"
	"fmt"

	"github.com/summa-tx/relays/golang/x/relay/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the genesis state
type GenesisState struct {
	Headers []string     `json:"headers"`
	Links   []types.Link `json:"links"`
}

// NewGenesisState instantiates a genesis state
func NewGenesisState(headers []string) GenesisState {
	return GenesisState{Headers: headers, Links: []types.Link{}}
}

// ValidateGenesis Validates a genesis state
func ValidateGenesis(data GenesisState) error {
	for _, header := range data.Headers {
		if len(header) != 160 {
			return fmt.Errorf("Invalid header, bad length")
		}

		_, err := hex.DecodeString(header)
		if err != nil {
			return err
		}
	}
	return nil
}

// DefaultGenesisState makes a default empty genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Headers: []string{},
		Links:   []types.Link{},
	}
}

// InitGenesis inits the app state based on the genesis state
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, header := range data.Headers {
		// can't fail if passes ValidateGenesis
		h, _ := hex.DecodeString(header)
		keeper.SetLink(ctx, h)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	panic("Not implemented")
}
