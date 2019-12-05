package cli

import (
	"encoding/json"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/summa-tx/relays/golang/x/relay/types"
)

// GetTxCmd sets up transaction CLI commands
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	relayTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Relay transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	relayTxCmd.AddCommand(client.PostCommands(
		GetCmdIngestHeaderChain(cdc),
	)...)

	return relayTxCmd
}

// GetCmdIngestHeaderChain creates a CLI command to ingest a header chain
func GetCmdIngestHeaderChain(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ingestheaders <json list of headers>",
		Short: "ingest a set of headers",
		Long:  "Ingest a set of headers. The headers must be in order, and the header immediately before the first must already be known to the relay",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// note: unmarshalling json here is undesirable, but necessary
			// 	     for MarkNewHeaviest we want to accept text instead of json
			var headers []types.BitcoinHeader
			jsonErr := json.Unmarshal([]byte(args[0]), &headers)
			if jsonErr != nil {
				return jsonErr
			}

			msg := types.NewMsgIngestHeaderChain(
				cliCtx.GetFromAddress(),
				headers,
			)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})

		},
	}
}

// // GetCmdSetLink is the CLI command for sending a SetLink transaction
// func GetCmdSetLink(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "setlink [hex header] [signer address]",
// 		Short: "Set a link",
// 		Args:  cobra.ExactArgs(2),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
//
// 			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
//
// 			addr := sdk.AccAddress(args[1])
// 			msg := types.NewMsgSetLink(addr, args[0], cliCtx.GetFromAddress())
// 			err := msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}
//
// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }
