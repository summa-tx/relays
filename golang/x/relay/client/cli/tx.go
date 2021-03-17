package cli

import (
	"encoding/json"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/summa-tx/relays/golang/x/relay/types"

	btcspv "github.com/summa-tx/bitcoin-spv/golang/btcspv"
)

// GetTxCmd sets up transaction CLI commands
func GetTxCmd(storeKey string, cdc *codec.LegacyAmino) *cobra.Command {
	relayTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Relay transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	relayTxCmd.AddCommand(client.PostCommands(
		GetCmdIngestHeaderChain(cdc),
		GetCmdIngestDifficultyChange(cdc),
		GetCmdNewRequest(cdc),
		GetCmdProvideProof(cdc),
		GetCmdMarkNewHeaviest(cdc),
	)...)

	return relayTxCmd
}

// GetCmdIngestHeaderChain creates a CLI command to ingest a header chain
func GetCmdIngestHeaderChain(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ingestheaders <json list of headers>",
		Example: "ingestheaders  2_ingest_headers.json --inputfile --from me",
		Short:   "Ingest a set of headers",
		Long: `Ingest a set of headers. The headers must be in order, and the header immediately before the first must already be known to the relay.
Use flag --inputfile to submit a json filename as input from scripts/seed_data directory`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)

			var headers = make([]types.BitcoinHeader, 0)

			if viper.GetBool("inputfile") {
				jsonFile, err := readJSONFromFile(args[0])
				if err != nil {
					return err
				}
				jsonErr := json.Unmarshal([]byte(jsonFile), &headers)
				if jsonErr != nil {
					return jsonErr
				}
			} else {
				jsonErr := json.Unmarshal([]byte(args[0]), &headers)
				if jsonErr != nil {
					return jsonErr
				}
			}

			msgs := []sdk.Msg{types.NewMsgWithdrawDelegatorReward(cliCtx.GetFromAddress(), headers)}

			// err := msg.ValidateBasic()
			// if err != nil {
			// 	return err
			// }

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msgs...)

		},
	}

	attachFlagFileinput(cmd)
	return cmd
}

// GetCmdIngestDifficultyChange creates a CLI command to ingest a difficulty change
func GetCmdIngestDifficultyChange(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ingestdiffchange <prev epoch start> <json list of headers>",
		Example: "ingestdiffchange ef8248820b277b542ac2a726ccd293e8f2a3ea24c1fe04000000000000000000  0_new_difficulty.json --inputfile --from me",
		Short:   "Ingest a difficulty change.",
		Long:    "Ingest a difficulty change. Prev Epoch Start is a hex digest.",
		Args:    cobra.ExactArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)

			prevEpochStart, err := btcspv.NewHash256Digest(btcspv.DecodeIfHex(args[0]))
			if err != nil {
				return types.FromBTCSPVError(types.DefaultCodespace, err)
			}

			var headers []types.BitcoinHeader
			if viper.GetBool("inputfile") {
				jsonFile, err := readJSONFromFile(args[1])
				if err != nil {
					return err
				}
				jsonErr := json.Unmarshal([]byte(jsonFile), &headers)
				if jsonErr != nil {
					return jsonErr
				}
			} else {
				jsonErr := json.Unmarshal([]byte(args[1]), &headers)
				if jsonErr != nil {
					return jsonErr
				}
			}

			msgs := []sdk.Msg{types.NewMsgIngestDifficultyChange(
				cliCtx.GetFromAddress(),
				prevEpochStart,
				headers,
			)}
			// err = msg.ValidateBasic()
			// if err != nil {
			// 	return err
			// }

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msgs...)

		},
	}

	attachFlagFileinput(cmd)
	return cmd
}

// GetCmdNewRequest stores a new proof request
func GetCmdNewRequest(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:     "newrequest <spends> <pays> <value> <numConfs>",
		Example: "newrequest 0x 17a91423737cd98bb6b2da5a11bcd82e5de36591d69f9f87 0 1 --from me",
		Short:   "Stores a new proof request",
		Long:    "Stores a new proof request",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)

			spends := btcspv.DecodeIfHex(args[0])
			pays := btcspv.DecodeIfHex(args[1])
			paysValue, valueErr := strconv.ParseUint(args[2], 10, 64)
			if valueErr != nil {
				return valueErr
			}
			numConfs, confsErr := strconv.ParseUint(args[3], 10, 8)
			if confsErr != nil {
				return confsErr
			}

			msgs := []sdk.Msg{types.NewMsgNewRequest(
				cliCtx.GetFromAddress(),
				spends,
				pays,
				paysValue,
				uint8(numConfs),
				types.Local,
				nil,
			)}
			// err := msg.ValidateBasic()
			// if err != nil {
			// 	return err
			// }

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msgs...)

		},
	}
}

// GetCmdProvideProof stores a new proof request
func GetCmdProvideProof(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "provideproof <json proof> <json list of requests>",
		Example: "provideproof 1_check_proof.json 3_filled_requests.json --inputfile --from me",
		Short:   "validates proof of given requests",
		Long: `Validates proof of given requests. Useful for validating proofs before spending gas on submission transaction.
Use flag --inputfile to submit a json filename as input from scripts/seed_data directory`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)

			var proof types.SPVProof
			var requests []types.FilledRequestInfo
			if viper.GetBool("inputfile") {
				jsonFileProof, err := readJSONFromFile(args[0])
				if err != nil {
					return err
				}
				jsonFileReq, err := readJSONFromFile(args[1])
				if err != nil {
					return err
				}
				jsonErr := json.Unmarshal([]byte(jsonFileProof), &proof)
				if jsonErr != nil {
					return jsonErr
				}
				jsonErr = json.Unmarshal([]byte(jsonFileReq), &requests)
				if jsonErr != nil {
					return jsonErr
				}
			} else {
				jsonErr := json.Unmarshal([]byte(args[0]), &proof)
				if jsonErr != nil {
					return jsonErr
				}
				jsonErr = json.Unmarshal([]byte(args[1]), &requests)
				if jsonErr != nil {
					return jsonErr
				}
			}

			filledRequests := types.NewFilledRequests(
				proof,
				requests,
			)

			msgs := []sdk.Msg{types.NewMsgProvideProof(
				cliCtx.GetFromAddress(),
				filledRequests,
			)}

			// err := msg.ValidateBasic()
			// if err != nil {
			// 	return err
			// }

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msgs...)

		},
	}

	attachFlagFileinput(cmd)
	return cmd
}

// GetCmdMarkNewHeaviest creates a CLI command to update best known digest and LCA
func GetCmdMarkNewHeaviest(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:     "marknewheaviest <ancestor> <currentBest (raw header)> <newBest (raw header)> [limit]",
		Example: "marknewheaviest 0x4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000 0x0000c020954ea1d980abc34fd5c260205e025a405f59cdf510960c000000000000000000ad864d04a6ca14e597da45c4936dd3a07946e7d72aab72a3ed7444f0f6da618dd150425eff3212173f0c982d 0x0000c020bc00d40ffb1b0e8850475b0ff71d990080bb0e8203d1090000000000000000008a317b377cc53010ed4c741bd6bcea5fe6748665a6a9374510ff77e5cdfac7e3b971425ed41a12174334a315 0 --from me",
		Short:   "Updates best known digest and LCA",
		Long:    "Updates best known digest and LCA.\nAncestor, current best, and new best are hex.\nLimit is an integer.",
		Args:    cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)


			// TODO: Set default limit if limit is not provided?

			ancestor, err := btcspv.NewHash256Digest(btcspv.DecodeIfHex(args[0]))
			if err != nil {
				return types.FromBTCSPVError(types.DefaultCodespace, err)
			}

			currentBest, curBestErr := btcspv.NewRawHeader(btcspv.DecodeIfHex(args[1]))
			if curBestErr != nil {
				return types.FromBTCSPVError(types.DefaultCodespace, curBestErr)
			}

			newBest, newBestErr := btcspv.NewRawHeader(btcspv.DecodeIfHex(args[2]))
			if newBestErr != nil {
				return types.FromBTCSPVError(types.DefaultCodespace, newBestErr)
			}

			limit, err := strconv.ParseUint(args[3], 10, 32)
			if err != nil {
				return err
			}

			msgs := []sdk.Msg{types.NewMsgMarkNewHeaviest(
				cliCtx.GetFromAddress(),
				ancestor,
				currentBest,
				newBest,
				uint32(limit),
			)}
			// err = msg.ValidateBasic()
			// if err != nil {
			// 	return err
			// }

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msgs...)
		},
	}
}
