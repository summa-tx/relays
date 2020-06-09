### Summa Relay

This is a Bitcoin Relay. It uses 1 + 1/n slots per header relayed (n is
currently 4), and 2 slots to externalize useful information (best chain tip and
best shared ancestor of latest reorg).

Implementations are available in Solidity (for EVM chains) and Golang using the
cosmos-sdk framework.

### How does it work?

The core idea behind the relay is to minimize storage costs by increasing
calldata costs. Rather than storing headers, the relay stores the
`hashPrevBlock`field of each header and the height of every nth header. Should
the relay need to reference information in old headers (like the difficulty),
the header data is passed to the relay again, and validated against known
`hashPrevBlock` links. This allows the relay to check that newly submitted
blocks are valid extensions of existing blocks, without storing all past header
information.

As opposed to other relays, we separate the function of the relay into two
categories: "learning about new blocks" and "following the best chain tip."
Users may add new blocks in groups of at least 5 by calling `addHeaders` If the
block slice includes a difficulty retarget, users are required to call
`addHeadersWithRetarget`, which performs additional validation. The relay does
not update its tip unless it is specifically requested to do so by a user. The
user must call `markNewHeaviest` with the new heaviest, the old heaviest
header, and the digest of their most recent common ancestor (which may be the
old heaviest header.

As part of the process, the relay externalizes the most recent common ancestor,
which is to say, the heaviest header that both old and new heaviest tip
confirm. This is a metric of "subjective finality" for that block. During
normal operation without reorgs it lags behind the tip by 5 blocks. During
reorgs, it is the shared base of the competing branches (and as such may move
backwards!). This indicates that competing sets of miners both viewed it as
subjectively finalized. As such, it is a reasonable source of finalization
information for relay-consuming smart contracts.

This model provides large gas savings compared to previous relay designs (TODO:
benchmarking). It also gets especially attractive if EIP2028 activates,
reducing calldata gas costs.

### A Note on Endianness

Bitcoin internally uses little-endian representations of integers and digests.
Block explorers and most user-facing applications use the more common
big-endian representation. To minimize order swaps and prevent confusion, all
our tooling uses the LE representation exclusively. If using the JS, rust,
golang, or python tooling in [bitcoin-spv](http://bitcoin-spv.com), everything
will Just Work. If writing custom software using data from block explorers,
full nodes, or other data sources, make sure digests are LE before submitting
to the relay.

### Requests and Proofs

The Relay implementations here have an SPV request system built in. This allows
for abstraction of the off-chain proving software. Requesters don't need to
write a custom Bitcoin indexer, and existing Bitcoin indexers can work with any
requester, whether it's a module, a smart contract, or a user.

The relay coordinates an interaction between 3 roles:
1. Requester: creates a new SPV Proof request and designates a Handler
2. Handler: handles incoming SPV Proofs on the Requester's behalf
3. Indexer: watches requests, indexes Bitcoin, and provides SPV Proofs

While implementation details differ, the architecture is simple:

1. Requesters register a request for SPV Proofs.
    1. The request specifies a transaction filter and a proof handler.
    1. golang: submit a `MsgRequestProof`.
    1. golang CLI: `relaycli tx relay newrequest`.
    1. solidity: `OndemandSPV.request()`.
1. An event with request details is logged.
    1. golang: watch for `proof_request` events.
    1. solidity: subscribe to `NewProofRequest` events.
1. Indexers watch the Bitcoin chain for transactions that satisfy Requests.
    1. [Example](https://github.com/summa-tx/bcoin-relaylib).
1. Indexers create an SPV Proof and submit it to the relay.
    1. golang: submit a `MsgProvideProof`.
    1. golang CLI: `relaycli tx relay provideproof`.
    1. solidity: call `OnDemandSPV.provideProof()`.
1. The relay validates this proof.
1. If valid, on-chain handler dispatches tx info to the proof Handler
    1. golang: the module's `ProofHandler` routes info the the Handler
    1. solidity: the relay calls `spv()` on the handling contract

Essentially the requester is subscribing to a feed of Bitcoin transactions
matching a specific filter. This filter can specify which UTXO is being spent,
and/or an address that receives funds. The handler expects to receive
a stream of transactions that meet the filter's specifications.

**Note**: Due to solidity constraints, this filter system is unrelated to
existing Bitcoin filtering systems (e.g. BIP37 & BIP157) In the future,
the filter system may be upgraded to support more complex transaction
descriptions.

**Important**: All requests may be filled more than once. Setting a `spends`
filter is NOT sufficient to prevent this, as long reorgs may cause a UTXO to be
spent multiple times. There is NO WAY to ensure that only a single proof is
provided, so the handler should deal with multiple proofs gracefully.

### Misc Project Notes

Complete relays are available in Solidity, for EVM-based chains (like Ethereum)
and Golang using the cosmos-sdk framework.

The Python relay mainter in `./maintainer/` is not thoroughly tested, and does
not yet support the cosmos-sdk relay.
