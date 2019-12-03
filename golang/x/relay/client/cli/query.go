package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

// GetQueryCmd sets up query CLI commands
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	relayQueryCommand := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the relay module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	relayQueryCommand.AddCommand(client.GetCommands(
		GetCmdIsAncestor(queryRoute, cdc), // @Erin: add new query commands to this list
		GetCmdFindAncestor(queryRoute, cdc),
		// GetCmd(queryRoute, cdc)
	)...)
	return relayQueryCommand
}

// GetCmdIsAncestor returns the CLI command struct for IsAncestor
func GetCmdIsAncestor(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		// what are the arguments. <> for required, [] for optional
		Use:     "isancestor <digest> <ancestor> [limit]",
		Example: "isancestor 12..ab 34..cd 200", // how do you use it?
		// a help message. shows on `help isancestor`
		Long: "Check if the second argument is an ancestor of the the argument. Optionally set a limit on block traversal",
		// how many arguments does it take?
		// also useful: cobra.ExactArgs(3)
		Args: cobra.RangeArgs(2, 3),
		// what does it do when run?
		RunE: func(cmd *cobra.Command, args []string) error {
			// spin up a context
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var limit uint32
			if len(args) == 3 {
				lim, err := strconv.ParseUint(args[2], 0, 32)
				if err != nil {
					fmt.Print(err.Error())
					return nil
				}
				limit = uint32(lim)
			}

			digestLE, sdkErr := types.Hash256DigestFromHex(args[0])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}
			ancestor, sdkErr := types.Hash256DigestFromHex(args[1])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			params := types.QueryParamsIsAncestor{
				DigestLE:            digestLE,
				ProspectiveAncestor: ancestor,
				Limit:               limit,
			}

			queryData, err := cdc.MarshalJSON(params)
			if err != nil {
				fmt.Print(err.Error())
				return nil
			}

			// run the query. the routeString is passed as strings to our querier switch/case in `keeper/querier.go`
			res, _, err := cliCtx.QueryWithData("custom/relay/isancestor", queryData)

			if err != nil {
				fmt.Printf("could not check if %s... is ancestor of %s... \n", args[1][:8], args[0][:8])
				return nil
			}

			var out types.QueryResIsAncestor
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}

func GetCmdFindAncestor(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		// what are the arguments. <> for required, [] for optional
		Use:     "findancestor <digest> <offset>",
		Example: "isancestor 12..ab 2", // how do you use it?
		// a help message. shows on `help isancestor`
		Long: "Finds the ancestor of a given digest",
		// how many arguments does it take?
		// also useful: cobra.ExactArgs(3)
		Args: cobra.ExactArgs(2),
		// what does it do when run?
		RunE: func(cmd *cobra.Command, args []string) error {
			// spin up a context
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			digest, sdkErr := types.Hash256DigestFromHex(args[0])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			var offset uint32
			off, err := strconv.ParseUint(args[1], 0, 32)
			if err != nil {
				fmt.Print(err.Error())
				return nil
			}
			offset = uint32(off)

			params := types.QueryParamsFindAncestor{
				DigestLE: digest,
				Offset:   offset,
			}

			queryData, err := cdc.MarshalJSON(params)
			if err != nil {
				fmt.Print(err.Error())
				return nil
			}

			// run the query. the routeString is passed as strings to our querier switch/case in `keeper/querier.go`
			res, _, err := cliCtx.QueryWithData("custom/relay/isancestor", queryData)

			if err != nil {
				fmt.Printf("could not find ancestor of %s... \n", args[0][:8])
				return nil
			}

			var out types.QueryResIsAncestor
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}
