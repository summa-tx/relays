### Summa Relay

This is a Bitcoin Relay. It uses 1 + 1/n slots per header relayed (n is
currently 4), and 2 slots to externalize useful information (best chain tip and
best shared ancestor of latest reorg).

At present, only a Solidity implementation is available, but we intend to add
more implementations soon :)

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

### Project Notes

Complete relays are available in Solidity, for EVM-based chains (like Ethereum)
and Golang using the cosmos-sdk framework.

The Python relay mainter in `./maintainer/` is not thoroughly tested, and does
not yet support the cosmos-sdk relay.
