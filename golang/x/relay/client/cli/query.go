package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

type QueryResGetParent struct {
	Digest types.Hash256Digest `json:"digest"`
}

func (d *QueryResGetParent) String() string {
	return fmt.Sprintf("%d", d.Digest)
}

// GetQueryCmd sets up query CLI commands
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	relayQueryCommand := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the relay module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	relayQueryCommand.AddCommand(client.GetCommands(
		GetCommandGetParent(storeKey, cdc),
	)...)
	return relayQueryCommand
}

// GetCommandGetParent queries information about a name
func GetCommandGetParent(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "getparent [digest]",
		Short: "getparent digest",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			digest := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getparent/%s", queryRoute, digest), nil)
			if err != nil {
				fmt.Printf("could not get parent - %s \n", digest)
				return nil
			}

			var out QueryResGetParent
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}
