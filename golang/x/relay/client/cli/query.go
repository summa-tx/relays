package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		GetCmdIsAncestor(queryRoute, cdc),
		GetCmdGetRelayGenesis(queryRoute, cdc),
		GetCmdGetLastReorgLCA(queryRoute, cdc),
		GetCmdGetBestDigest(queryRoute, cdc),
		GetCmdFindAncestor(queryRoute, cdc),
		GetCmdIsMostRecentCommonAncestor(queryRoute, cdc),
		GetCmdHeaviestFromAncestor(queryRoute, cdc),
		GetCmdCheckProof(queryRoute, cdc),
		GetCmdCheckRequests(queryRoute, cdc),
	)...)
	return relayQueryCommand
}

// GetCmdIsAncestor returns the CLI command struct for IsAncestor
func GetCmdIsAncestor(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		// what are the arguments. <> for required, [] for optional
		Use:     "isancestor <digest> <ancestor> [limit]",
		Example: "isancestor 0641238051855d1759da9b6603b156684a68a146d36a09000000000000000000 9be6406d5311123b6212b14b1a070276157364a5d5f004000000000000000000 200", // how do you use it?
		// a help message. shows on `help isancestor`
		Long: "Check if the second argument is an ancestor of the the argument. Optionally set a limit on block traversal",
		// how many arguments does it take?
		// also useful: cobra.ExactArgs(3)
		Args: cobra.RangeArgs(2, 3),
		// what does it do when run?
		RunE: func(cmd *cobra.Command, args []string) error {
			// spin up a context
			cliCtx, err := client.GetClientQueryContext(cmd)

			var limit uint32
			if len(args) == 3 {
				lim, err := strconv.ParseUint(args[2], 10, 32)
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

// GetCmdGetRelayGenesis returns the CLI command struct for GetRelayGenesis
func GetCmdGetRelayGenesis(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "getrelaygenesis",
		Example: "getrelaygenesis",
		Long:    "Get the first digest in the relay",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)

			res, _, err := cliCtx.QueryWithData("custom/relay/getrelaygenesis", nil)

			if err != nil {
				fmt.Println("could not get the first digest in the relay")
				return nil
			}

			var out types.QueryResGetRelayGenesis
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}

// GetCmdGetLastReorgLCA returns the CLI command struct for GetLastReorgLCA
func GetCmdGetLastReorgLCA(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "getlastreorglca",
		Example: "getlastreorglca",
		Long:    "Returns the latest common ancestor of the last-known reorg",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)

			res, _, err := cliCtx.QueryWithData("custom/relay/getlastreorglca", nil)

			if err != nil {
				fmt.Println("could not get the last Reorg LCA")
				return nil
			}

			var out types.QueryResGetLastReorgLCA
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}

// GetCmdGetBestDigest returns the CLI command struct for GetBestDigest
func GetCmdGetBestDigest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "getbestdigest",
		Example: "getbestdigest",
		Long:    "Returns the best known digest in the relay",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)

			res, _, err := cliCtx.QueryWithData("custom/relay/getbestdigest", nil)

			if err != nil {
				fmt.Println("could not get best known digest")
				return nil
			}

			var out types.QueryResGetBestDigest
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}

// GetCmdFindAncestor returns the CLI command struct for FindAncestor
func GetCmdFindAncestor(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "findancestor <digest> <offset>",
		Example: "findancestor f8d0a038bfe4027e5de3b6bf07262122636fd2916d7503000000000000000000 2", // how do you use it?
		Long:    "Finds the digest <offset> blocks before <digest>. Errors if digest or the ancestor is unknown",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)

			digest, sdkErr := types.Hash256DigestFromHex(args[0])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			var offset uint32
			off, err := strconv.ParseUint(args[1], 10, 32)
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

			res, _, err := cliCtx.QueryWithData("custom/relay/findancestor", queryData)

			if err != nil {
				fmt.Printf("could not find ancestor of %s... \n", args[0][:8])
				return nil
			}

			var out types.QueryResFindAncestor
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}

// GetCmdIsMostRecentCommonAncestor returns the CLI command struct for IsMostRecentCommonAncestor
func GetCmdIsMostRecentCommonAncestor(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "ismostrecentcommonancestor <ancestor> <left> <right> [limit]",
		Example: "ismostrecentcommonancestor 9be6406d5311123b6212b14b1a070276157364a5d5f004000000000000000000 0641238051855d1759da9b6603b156684a68a146d36a09000000000000000000 7f9923db1d3ad6a08054b4a80a5cd7478b57a9650eaf09000000000000000000 3", // how do you use it?
		Long:    "Checks if <ancestor> is the LCA of <left> and <right> digests",
		Args:    cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)

			ancestor, sdkErr := types.Hash256DigestFromHex(args[0])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			left, sdkErr := types.Hash256DigestFromHex(args[1])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			right, sdkErr := types.Hash256DigestFromHex(args[2])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			var limit uint32
			if len(args) == 4 {
				lim, err := strconv.ParseUint(args[3], 10, 32)
				if err != nil {
					fmt.Print(err.Error())
					return nil
				}
				limit = uint32(lim)
			}

			params := types.QueryParamsIsMostRecentCommonAncestor{
				Ancestor: ancestor,
				Left:     left,
				Right:    right,
				Limit:    limit,
			}

			queryData, err := cdc.MarshalJSON(params)
			if err != nil {
				fmt.Print(err.Error())
				return nil
			}

			res, _, err := cliCtx.QueryWithData("custom/relay/ismostrecentcommonancestor", queryData)

			if err != nil {
				fmt.Printf("could not check if %s... is the LCA of %s... and %s... \n", args[0][:8], args[1][:8], args[2][:8])
				return nil
			}

			var out types.QueryResIsMostRecentCommonAncestor
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}

// GetCmdHeaviestFromAncestor returns the CLI command struct for HeaviestFromAncestor
func GetCmdHeaviestFromAncestor(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "heaviestfromancestor <ancestor> <currentbest> <newbest> [limit]",
		Example: "heaviestfromancestor 4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000 4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000 0641238051855d1759da9b6603b156684a68a146d36a09000000000000000000 200", // how do you use it?
		Long:    "Determines the heavier descendant of a common ancestor",
		Args:    cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)

			ancestor, sdkErr := types.Hash256DigestFromHex(args[0])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			currentBest, sdkErr := types.Hash256DigestFromHex(args[1])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			newBest, sdkErr := types.Hash256DigestFromHex(args[2])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			var limit uint32
			if len(args) == 4 {
				lim, err := strconv.ParseUint(args[3], 10, 32)
				if err != nil {
					fmt.Print(err.Error())
					return nil
				}
				limit = uint32(lim)
			}

			params := types.QueryParamsHeaviestFromAncestor{
				Ancestor:    ancestor,
				CurrentBest: currentBest,
				NewBest:     newBest,
				Limit:       limit,
			}

			queryData, err := cdc.MarshalJSON(params)
			if err != nil {
				fmt.Print(err.Error())
				return nil
			}

			res, _, err := cliCtx.QueryWithData("custom/relay/heaviestfromancestor", queryData)

			if err != nil {
				fmt.Printf("could not determine if %s... or %s... is heaviest decendant of %s... \n", args[1][:8], args[2][:8], args[0][:8])
				return nil
			}

			var out types.QueryResHeaviestFromAncestor
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}

// GetCmdGetRequest returns the CLI command struct for getRequest
func GetCmdGetRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "getrequest <id>",
		Example: "getrequest 12",
		Long:    "Get a proof request using the associated ID. ID can be an\n\"0x\" prepended hexbyte string or an integer",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)

			id, idErr := types.RequestIDFromString(args[0])
			if idErr != nil {
				return idErr
			}

			params := types.QueryParamsGetRequest{
				ID: id,
			}

			queryData, err := cdc.MarshalJSON(params)
			if err != nil {
				fmt.Print(err.Error())
				return nil
			}

			res, _, err := cliCtx.QueryWithData("custom/relay/getrequest", queryData)

			if err != nil {
				fmt.Printf("could not find request associated with id: %s... \n", args[0])
				return nil
			}

			var out types.QueryResGetRequest
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}
}

// GetCmdCheckRequests returns the CLI command struct for checkRequests
func GetCmdCheckRequests(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "checkrequests <json proof> <json list of requests>",
		Example: "checkrequests 1_check_proof.json 3_filled_requests.json --inputfile",
		Long: `Check whether proof successfully validates a set of requests.
Use flag --inputfile to submit a json filename as input from scripts/seed_data directory`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)

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

			params := types.QueryParamsCheckRequests{
				Filled: filledRequests,
			}

			queryData, err := cdc.MarshalJSON(params)
			if err != nil {
				fmt.Print(err.Error())
				return nil
			}

			res, _, err := cliCtx.QueryWithData("custom/relay/checkrequests", queryData)

			if err != nil {
				fmt.Printf("error processing checkrequests: %s \n", err)
				return nil
			}

			var out types.QueryResCheckRequests
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}

	attachFlagFileinput(cmd)
	return cmd
}

// GetCmdCheckProof returns the CLI command struct for checkProof
func GetCmdCheckProof(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "checkproof <json proof>",
		Example: "checkproof 1_check_proof.json --inputfile",
		Long: `check proof has valid parameters.
Use flag --inputfile to submit a json filename as input from scripts/seed_data directory.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)

			var proof types.SPVProof
			if viper.GetBool("inputfile") {
				jsonFile, err := readJSONFromFile(args[0])
				if err != nil {
					return err
				}
				jsonErr := json.Unmarshal([]byte(jsonFile), &proof)
				if jsonErr != nil {
					return jsonErr
				}
			} else {
				jsonErr := json.Unmarshal([]byte(args[0]), &proof)
				if jsonErr != nil {
					return jsonErr
				}
			}

			params := types.QueryParamsCheckProof{
				Proof: proof,
			}

			queryData, err := cdc.MarshalJSON(params)
			if err != nil {
				fmt.Print(err.Error())
				return nil
			}

			res, _, err := cliCtx.QueryWithData("custom/relay/checkproof", queryData)

			if err != nil {
				fmt.Printf("error processing checkproof: %s \n", err)
				return nil
			}

			var out types.QueryResCheckProof
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(&out)
		},
	}

	attachFlagFileinput(cmd)
	return cmd
}
