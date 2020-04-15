## cosmos-sdk Bitcoin Relay

This is a full-featured Bitcoin relay module for cosmos-sdk chains. It indexes
Bitcoin headers, provides information about the latest-known state of the
Bitcoin chain, and validates SPV Proofs against its view of the chain. It is a
critical component for many Cosmos applications to interact with Bitcoin.

For more information about the relay's architecture see the `README.md` in the
repo's root directory.

## Building the daemon and cli

```sh
$ make
# To additionally install them in your `$GOPATH/bin` directory:
$ make install
```

## Running tests

Run the unit tests as follows:

```sh
$ go test ./x/...
```

See the README in `./cli_test` for instructions on running the integration
tests.

Instructions for setting up manual testing can be found in the README in
`./scripts`.

## Project Status

- [X] Milestone 1
- [X] Milestone 2
- - [X] Expose best-known digest
- - [X] Expose LCA of reorg
- - [X] Follow API of existing Solidity Relay
- - [X] Validate SPV Proofs
- - [X] `ProvideProof` message
- [ ] Milestone 3
- - [X] Provide tooling for manual testing (scripts, docs, json test files)
- - [X] Integration Tests
- - [ ] Document relay design & architecture
- - [X] Document public interface
- - [X] Provide hooks to execute tasks + dispatch messages
- - [ ] Add a basic web dashboard with Relay health


## API

Cosmos modules expose messages, which modify state, and queries, which read
state.

### Queries

Queries are available via CLI or REST. For more information, see the
descriptions in the CLI.

| Query | Description |
| ----- | ----------- |
| IsAncestor | Deteremine if a block is an ancestor of another |
| GetRelayGenesis | Get the trusted root of the relay |
| GetLastReorgLCA | Get the LCA of the latest reorg |
| GetLastReorgLCA | Get the best digest known to the relay |
| FindAncestor | Find the nth ancestor of a block|
| IsMostRecentCommonAncestor | Determine if a block is the LCA of two headers|
| HeaviestFromAncestor | Check which of two descendents is heaviest from the LCA |
| GetRequest | Get details of an SPV Proof Request|
| CheckProof | Check the syntactic validity of an SPV Proof |
| CheckRequests | Perform CheckProof and check the SPV Proof against a set of Requests |

### Messages

Messages are available via CLI or REST. For more information, see the
descriptions in the CLI.

| Message | Description |
| ------- | ----------- |
| IngestHeaderChain | Add a chain of headers to the relay |
| IngestDifficultyChange | Add a chain of headers to the relay with a difficulty change|
| MarkNewHeaviest | Mark a new best-known chain tip |
| NewRequest | Register a new SPV Proof request |
| ProvideProof | Provide a proof that satisfies 1 or more requests |

## Project Overview

### Keeper
High-level overview of the project structure within the `keeper` file.

#### Keeper.go
Instantiates a `keeper` (what handles interaction with the store and contains most of the core functionality of the module). It also handles the genesis state for the relay.

#### Headers.go
Handles the storage and validation of Bitcoin Headers and Header Chains.

#### Chain.go
Checks and updates information about the chain.  Provides functionality to ensure we are using the heaviest chain.

#### Links.go
Sets and retrieves data about each link in the chain.  This is most commonly used to check information about ancestors.

#### Requests.go
Stores, retrieves, and validates requests.

#### Validator.go
Contains validation functions.  Currently, this can validate SPV Proofs and Requests.

#### Handler.go
Handles messages.

#### Querier.go
Handles queries.
