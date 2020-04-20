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

Cosmos modules expose queries (which read state) and messages (which modify
state).  These are available via CLI or REST.

### CLI

Queries

| Query | Description | Usage |
| ----- | ----------- | ------|
| IsAncestor | Deteremine if a block is an ancestor of another | `isancestor <digest> <ancestor> [limit]` |
| GetRelayGenesis | Get the trusted root of the relay | `getrelaygenesis` |
| GetLastReorgLCA | Get the LCA of the latest reorg | `getlastreorglca` |
| GetBestDigest | Get the best digest known to the relay | `getbestdigest` |
| FindAncestor | Find the nth ancestor of a block | `findancestor <digest> <offset>` |
| IsMostRecentCommonAncestor | Determine if a block is the LCA of two headers | `ismostrecentcommonancestor <ancestor> <left> <right> [limit]` |
| HeaviestFromAncestor | Check which of two descendents is heaviest from the LCA | `heaviestfromancestor <ancestor> <currentbest> <newbest> [limit]` |
| GetRequest | Get details of an SPV Proof Request | `getrequest <id>` |
| CheckProof | Check the syntactic validity of an SPV Proof | `checkproof <json proof>` |
| CheckRequests | Perform CheckProof and check the SPV Proof against a set of Requests | `checkrequests <json proof> <json list of requests>` |

Messages

| Message | Description | Usage |
| ------- | ----------- | ----- |
| IngestHeaderChain | Add a chain of headers to the relay | `ingestheaders <json list of headers>` |
| IngestDifficultyChange | Add a chain of headers to the relay with a difficulty change | `ingestdiffchange <prev epoch start> <json list of headers>` |
| MarkNewHeaviest | Mark a new best-known chain tip | `marknewheaviest <ancestor> <currentBest> <newBest> [limit]` |
| NewRequest | Register a new SPV Proof request | `newrequest <spends> <pays> <value> <numConfs>` |
| ProvideProof | Provide a proof that satisfies 1 or more requests | `provideproof <json proof> <json list of requests>` |

### REST Routes

Query routes

| Endpoint | Function | Description | Type |
| /isancestor/{digest}/{ancestor}/ | IsAncestor | Deteremine if a block is an ancestor of another | GET |
| /isancestor/{digest}/{ancestor}/{limit} | IsAncestor | Deteremine if a block is an ancestor of another | GET |
| /getrelaygenesis | GetRelayGenesis | Get the trusted root of the relay | GET |
| /getlastreorglca | GetLastReorgLCA | Get the LCA of the latest reorg | GET |
| /getbestdigest | GetBestDigest | Get the best digest known to the relay | GET |
| /findancestor/{digest}/{offset} | FindAncestor | Find the nth ancestor of a block | GET |
| /ismostrecentcommonancestor/{ancestor}/{left}/{right}/ | IsMostRecentCommonAncestor | Determine if a block is the LCA of two headers | GET |
| /ismostrecentcommonancestor/{ancestor}/{left}/{right}/{limit} | IsMostRecentCommonAncestor | Determine if a block is the LCA of two headers | GET |
| /heaviestfromancestor/{ancestor}/{currentBest}/{newBest}/ |  |  | GET |
| /heaviestfromancestor/{ancestor}/{currentBest}/{newBest}/{limit} |  |  | GET |
| /getrequest/{id} | GetRequest | Get details of an SPV Proof Request | GET |
| /checkrequests | CheckRequests | Perform CheckProof and check the SPV Proof against a set of Requests | POST |
| /checkproof | CheckProof | Check the syntactic validity of an SPV Proof | POST |

Message routes

| Endpoint | Function | Description | Type |
| -------- | -------- | ----------- | ---- |
| /ingestheaderchain | IngestHeaderChain | Add a chain of headers to the relay | POST |
| /ingestdiffchange | IngestDifficultyChange | Add a chain of headers to the relay with a difficulty change | POST |
| /marknewheaviest | MarkNewHeaviest | Mark a new best-known chain tip | POST |
| /newrequest | NewRequest | Register a new SPV Proof request | POST |
| /provideproof | ProvideProof | Provide a proof that satisfies 1 or more requests | POST |

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
