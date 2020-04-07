## TODOs:

- [X] Milestone 1
- [X] Milestone 2
- - [X] Expose best-known digest
- - [X] Expose LCA of reorg
- - [X] Follow API of existing Solidity Relay
- - [X] Validate SPV Proofs
- - [X] `ProvideProof` message
- [ ] Milestone 3
- - [ ] Provide tooling for manual testing (scripts, docs, json test files)
- - [ ] Integration Tests
- - [ ] Document relay design & architecture
- - [ ] Document public interface
- - [ ] Provide hooks to execute tasks + dispatch messages
- - [ ] Add a basic web dashboard with Relay health

## Messages added to CLI

| Message | Status | Description |
| ------- | ------ | ----------- |
| IngestHeaderChain | Completed | Add a chain of headers to the relay |
| IngestDifficultyChange | Completed | Add a chain of headers to the relay with a difficulty change|
| MarkNewHeaviest | Completed | Mark a new best-known chain tip |
| NewRequest | Completed | Register a new SPV Proof request|
| ProvideProof | Completed | Provide a proof that satisfies 1 or more requests |

## Queries added to CLI

| Query | Status | Description |
| ----- | ------ | ----------- |
| IsAncestor | Completed | Deteremine if a block is an ancestor of another |
| GetRelayGenesis | Completed | Get the trusted root of the relay |
| GetLastReorgLCA | Completed | Get the LCA of the latest reorg |
| FindAncestor | Completed | Find the nth ancestor of a block|
| IsMostRecentCommonAncestor | Completed | Determine if a block is the LCA of two headers|
| HeaviestFromAncestor | Completed | Check which of two descendents is heaviest from the LCA |
| GetRequest | Completed | Get details of an SPV Proof Request|
| CheckProof | Completed | Check the syntactic validity of an SPV Proof |
| CheckRequests | Completed | Perform CheckProof and check the SPV Proof against a set of Requests |

## How to add a view function (queries)
1. Add necessary getter(s) in `x/relay/keeper/keeper.go`
1. Add response type to `x/relay/types/querier.go`
    1. Add new string tag for the new query
    1. Response type is a struct with the return values
    1. Implement `String()` for the response type
1. Add function to querier `x/relay/keeper/querier.go`
    1. Add new `query___` function
    1. Add new case block to `switch` in `NewQuerier()`
1. Add to CLI  
    1. add to `x/relay/client/cli/query.go`
      1. `func GetCmd______`
      1. returns a `cobra.Command` object
      1. define `Use` `Example` `Short` `Long` `Args` and `RunE`
      1. `RunE` parses args, returns errors, and calls `cliCtx.QueryWithData`
      1. parses the output and returns it with `cliCtx.PrintOutput`
1. Add to REST
    1. add to `x/relay/client/rest/query.go`
    1. new function `_____Handler`
      1. parse args and build structs
      1. cliCtx.QueryWithData
      1. return errors with `rest.WriteErrorResponse`
      1. return query result with `rest.PostProcessResponse`
    1. add GET route to `x/relay/client/rest/rest.go`
      1. new `s.HandleFunc` with the route and arguments
      1. `.Methods("GET")`
      1. duplicate for optional args (see `isancestor` for example)


## How to add a non-view function (messages)
1. Add necessary getters/setters in `x/relay/keeper/keeper.go`
1. Add msg type in `x/relay/types/msgs.go`
    1. Message type is a struct with the arguments
    1. Implement `New___()`
    1. Implement `GetSigners()` <--- Ask me about this later
    1. Implement `Type()`
    1. Implement `ValidateBasic()`
    1. Implement `GetSignBytes()`
    1. Implement `Route()`
1. Add to handler
    1. Add new `handle____` function
    1. Add new case block to `switch` in `NewHandler()`
1. Add aliases in `x/relay/alias.go`
    1. Add alias in `var` block
    1. Add alias in `type` block
1. Add to CLI  
    1. add to `x/relay/client/cli/tx.go`
      1. `func GetCmd______`
      1. returns a `cobra.Command` object
      1. define `Use` `Example` `Short` `Long` `Args` and `RunE`
      1. `RunE` parses args, returns errors, and calls `utils.GenerateOrBroadcastMsgs`
1. Add to REST
    1. add to `x/relay/client/rest/tx.go`
      1. new http request struct `______Req`
        1. `BaseReq` + the struct from `x/relay/types/msgs.go`
      1. new function `_____Handler`
        1. parse args and build structs
        1. return errors with `rest.WriteErrorResponse`
        1. make the tx with `utils.WriteGenerateStdTxResponse`
    1. add POST route to `x/relay/client/rest/rest.go`
      1. new `s.HandleFunc` with the route and arguments
      1. `.Methods("POST")`
