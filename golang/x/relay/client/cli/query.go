package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/summa-tx/relays/proto"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

// TODO: Should this have the queryRoute param or not?
// GetQueryCmd sets up query CLI commands
func GetQueryCmd(queryRoute string) *cobra.Command {
	relayQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the relay module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	relayQueryCmd.AddCommand(
		GetCmdIsAncestor(queryRoute),
		GetCmdGetRelayGenesis(queryRoute),
		GetCmdGetLastReorgLCA(queryRoute),
		GetCmdGetBestDigest(queryRoute),
		GetCmdFindAncestor(queryRoute),
		GetCmdIsMostRecentCommonAncestor(queryRoute),
		GetCmdHeaviestFromAncestor(queryRoute),
		GetCmdCheckProof(queryRoute),
		GetCmdCheckRequests(queryRoute),
	)

	return relayQueryCmd
}

// GetCmdIsAncestor returns the CLI command struct for IsAncestor
func GetCmdIsAncestor(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
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
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			var limit uint32
			if len(args) == 3 {
				lim, err := strconv.ParseUint(args[2], 10, 32)
				if err != nil {
					fmt.Print(err.Error())
					return nil
				}
				limit = uint32(lim)
			}

			digestLE, sdkErr := types.BytesFromHex(args[0])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}
			ancestor, sdkErr := types.BytesFromHex(args[1])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			params := &proto.QueryParamsIsAncestor{
				DigestLE:            digestLE,
				ProspectiveAncestor: ancestor,
				Limit:               limit,
			}

			res, err := queryClient.IsAncestor(cmd.Context(), params)
			if err != nil {
				fmt.Printf("could not check if %s... is ancestor of %s... \n", args[1][:8], args[0][:8])
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdGetRelayGenesis returns the CLI command struct for GetRelayGenesis
func GetCmdGetRelayGenesis(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "getrelaygenesis",
		Example: "getrelaygenesis",
		Long:    "Get the first digest in the relay",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			res, err := queryClient.GetRelayGenesis(cmd.Context(), nil)
			if err != nil {
				fmt.Println("could not get the first digest in the relay")
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdGetLastReorgLCA returns the CLI command struct for GetLastReorgLCA
func GetCmdGetLastReorgLCA(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "getlastreorglca",
		Example: "getlastreorglca",
		Long:    "Returns the latest common ancestor of the last-known reorg",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			res, err := queryClient.GetLastReorgLCA(cmd.Context(), nil)
			if err != nil {
				fmt.Println("could not get the last Reorg LCA")
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdGetBestDigest returns the CLI command struct for GetBestDigest
func GetCmdGetBestDigest(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "getbestdigest",
		Example: "getbestdigest",
		Long:    "Returns the best known digest in the relay",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			res, err := queryClient.GetBestDigest(cmd.Context(), nil)
			if err != nil {
				fmt.Println("could not get best known digest")
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdFindAncestor returns the CLI command struct for FindAncestor
func GetCmdFindAncestor(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "findancestor <digest> <offset>",
		Example: "findancestor f8d0a038bfe4027e5de3b6bf07262122636fd2916d7503000000000000000000 2", // how do you use it?
		Long:    "Finds the digest <offset> blocks before <digest>. Errors if digest or the ancestor is unknown",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			digest, sdkErr := types.BytesFromHex(args[0])
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

			params := &proto.QueryParamsFindAncestor{
				DigestLE: digest,
				Offset:   offset,
			}

			res, err := queryClient.FindAncestor(cmd.Context(), params)
			if err != nil {
				fmt.Printf("could not find ancestor of %s... \n", args[0][:8])
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdIsMostRecentCommonAncestor returns the CLI command struct for IsMostRecentCommonAncestor
func GetCmdIsMostRecentCommonAncestor(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ismostrecentcommonancestor <ancestor> <left> <right> [limit]",
		Example: "ismostrecentcommonancestor 9be6406d5311123b6212b14b1a070276157364a5d5f004000000000000000000 0641238051855d1759da9b6603b156684a68a146d36a09000000000000000000 7f9923db1d3ad6a08054b4a80a5cd7478b57a9650eaf09000000000000000000 3", // how do you use it?
		Long:    "Checks if <ancestor> is the LCA of <left> and <right> digests",
		Args:    cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			ancestor, sdkErr := types.BytesFromHex(args[0])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			left, sdkErr := types.BytesFromHex(args[1])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			right, sdkErr := types.BytesFromHex(args[2])
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

			params := &proto.QueryParamsIsMostRecentCommonAncestor{
				Ancestor: ancestor,
				Left:     left,
				Right:    right,
				Limit:    limit,
			}

			res, err := queryClient.IsMostRecentCommonAncestor(cmd.Context(), params)
			if err != nil {
				fmt.Printf("could not check if %s... is the LCA of %s... and %s... \n", args[0][:8], args[1][:8], args[2][:8])
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdHeaviestFromAncestor returns the CLI command struct for HeaviestFromAncestor
func GetCmdHeaviestFromAncestor(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "heaviestfromancestor <ancestor> <currentbest> <newbest> [limit]",
		Example: "heaviestfromancestor 4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000 4c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000 0641238051855d1759da9b6603b156684a68a146d36a09000000000000000000 200", // how do you use it?
		Long:    "Determines the heavier descendant of a common ancestor",
		Args:    cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			ancestor, sdkErr := types.BytesFromHex(args[0])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			currentBest, sdkErr := types.BytesFromHex(args[1])
			if sdkErr != nil {
				fmt.Print(sdkErr.Error())
				return nil
			}

			newBest, sdkErr := types.BytesFromHex(args[2])
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

			params := &proto.QueryParamsHeaviestFromAncestor{
				Ancestor:    ancestor,
				CurrentBest: currentBest,
				NewBest:     newBest,
				Limit:       limit,
			}

			res, err := queryClient.HeaviestFromAncestor(cmd.Context(), params)
			if err != nil {
				fmt.Printf("could not determine if %s... or %s... is heaviest decendant of %s... \n", args[1][:8], args[2][:8], args[0][:8])
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdGetRequest returns the CLI command struct for getRequest
func GetCmdGetRequest(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "getrequest <id>",
		Example: "getrequest 12",
		Long:    "Get a proof request using the associated ID. ID can be an\n\"0x\" prepended hexbyte string or an integer",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			id, idErr := types.BytesFromHex(args[0])
			if idErr != nil {
				return idErr
			}

			params := &proto.QueryParamsGetRequest{
				ID: id,
			}

			res, err := queryClient.GetRequest(cmd.Context(), params)
			if err != nil {
				fmt.Printf("could not find request associated with id: %s... \n", args[0])
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdCheckRequests returns the CLI command struct for checkRequests
func GetCmdCheckRequests(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "checkrequests <json proof> <json list of requests>",
		Example: "checkrequests 1_check_proof.json 3_filled_requests.json --inputfile",
		Long: `Check whether proof successfully validates a set of requests.
Use flag --inputfile to submit a json filename as input from scripts/seed_data directory`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			var proof proto.SPVProof
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

			filledRequests := proto.FilledRequests{
				proof,
				requests,
			}

			params := &proto.QueryParamsCheckRequests{
				Filled: &filledRequests,
			}

			res, err := queryClient.CheckRequests(cmd.Context(), params)
			if err != nil {
				fmt.Printf("error processing checkrequests: %s \n", err)
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	attachFlagFileinput(cmd)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdCheckProof returns the CLI command struct for checkProof
func GetCmdCheckProof(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "checkproof <json proof>",
		Example: "checkproof 1_check_proof.json --inputfile",
		Long: `check proof has valid parameters.
Use flag --inputfile to submit a json filename as input from scripts/seed_data directory.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := proto.NewQueryClient(clientCtx)

			var proof proto.SPVProof
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

			params := &proto.QueryParamsCheckProof{
				Proof: proof,
			}

			res, err := queryClient.CheckProof(cmd.Context(), params)
			if err != nil {
				fmt.Printf("error processing checkproof: %s \n", err)
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	attachFlagFileinput(cmd)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
