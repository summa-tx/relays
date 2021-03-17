## Adding new functionality

This is a cosmos-sdk module. It can be extended with new messages and/or
queries. Generally, this module is feature-complete, and should not be
extended. The main exception is the WIP hooks system on proof validation. All
other functionality should likely be put into a separate module.

### Integrating with other modules

The relay keeper keeps a reference to an object that implements the following
interface (found in `x/types/types.go`).

```go
type ProofHandler interface {
	HandleValidProof(ctx sdk.Context, filled FilledRequests, requests []ProofRequest)
}
```

The `FilledRequests` struct contains an `SPVProof` and supporting information
about the transaction that fulfills the request.
It can be found in `x/types/validator.go`. `requests []ProofRequest` is a slice
of `ProofRequest`s that have been filled.

When the keeper validates a proof, it will call the `HandleValidProof` function
with the valid `FilledRequests` struct and the `ProofRequests` that have been
filled.

First, instantiate a `handler` that fulfills the `ProofHandler` interface. Then
add an instance of `relay.Keeper` to your app in `app.go`. It can be
instantiated as follows:

```go
handler = types.NewNullHandler()  // or your preferred handler

app.relayKeeper = relay.NewKeeper(
  keys[relay.StoreKey],
  app.cdc,
  true,
  handler
)
```

After that, the relay can be accessed via the Keeper's public interface.

### Extending this module

In order to extend this module, follow these steps:

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
      1. parses the output and returns it with `cliCtx.PrintProto`
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
