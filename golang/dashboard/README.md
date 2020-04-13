# cosmos-relay-dashboard

## TODOs

Functional:

- [ ] Convert socket calls to REST. [Available Routes](https://github.com/summa-tx/relays/blob/master/golang/x/relay/client/rest/rest.go)
- [ ] Add poll time for REST calls

Mostly pretty but also functional:

- [ ] Show height for best known digest and updatedAt
- [ ] Show height for lca and updatedAt
- [ ] Format date strings
- [ ] Add info tooltip explanations
- [ ] Separate info at bottom into two areas (relay, external source)
- [ ] Click to copy broken

## Description

The dashboard displays the Cosmos Relay chain data and verifies it against an external Bitcoin explorer (currently BlockCypher).

--------------------------

## Getting Started

### Cosmos Relay

The dashboard connects with a locally run `relay`. Follow [these instructions](https://github.com/summa-tx/relays/blob/master/golang/scripts/README.md) to build and run the the chain.

Once the chain is up, start the `REST` routes in a separate terminal. This will make the relay application available on `http://localhost:1317`.

### Dashboard

#### Install dependencies:

```
$ npm install
```

#### Start dashboard

```sh
$ npm run serve
```

View at http://localhost:8080 in your browser.

### Development

#### Set Environment Variables

If no `.env` file is present, defaults are used. See `/src/config.js`.

### Commands

#### Start dashboard

Compiles and hot-reloads.

```sh
$ npm run serve
```

#### Run your tests

```sh
$ npm run test
```

#### Lints and fixes files

```sh
$ npm run lint
```

#### Compiles and minifies for production

```sh
$ npm run build
```

### Customize configuration

See [Configuration Reference](https://cli.vuejs.org/config/).

## How Things Work

### Dashboard Overview

### Dashboard - Newest Header

The user wants to know about new headers. In order to do that, we:

1. Get the best tip (most recent block height) from an external source
2. Check if relay can this block with `findHeight`
    1. If yes, then display this height along with the block hash
    2. Otherwise, show flag that this isn't verified by the relay.

These are conceptually equivalent to Github commits.

### Dashboard - Best Known Digest

This is the block that is the best. It is updated approximately every 5 blocks, and will be behind the newest header.
We should buffer against this.

1. Listen for "Reorg" contract events
2. Display and store digest when it changes. Also display height.

This is conceptually equivalent to Github tags.

#### Dashboard - Last (Reorg) Common Ancestor (LCA)

This is the latest block that is in the history of both the current best known digest, and the previous best known digest.

1. Call `getLastReorgCommonAncestor()` to get LCA

#### Dashboard - Health Checks and Verification

The dashboard keeps track of and displays the following:

* **lastComms**: When was the last successful communication made?
  * **lastComms.relay**
  * **lastComms.external**

* **currentBlock.verifiedAt**: When was the current block verified i.e. When did `findHeight` return true?

* **previousBlock.verifiedAt**: When was the previous block verified?

Health pulses are displayed as `TIME in MINUTES ago`.

### Dashboard - Networks Names
Display for the user what network they are on in the format `eth_network_name-bitcoin_network_name`.

Examples:
- If we are on Ropsten and BTC Testnet, the network should be displayed as `ropsten-test`
- If we are on Celo BTC Mainnet, the network should be displayed as `alfajores-test`

--------------------------

### Relay

The following is mainly for informational purposes, rather than development.

#### Relay updates

The relay is updated ~every 5 blocks

**Advance chaining:**
Suppose this happens:

```
                                     BEST
                                      V
500 <- 501 <- 502 <- 503a <- 504  <- 505
               ^
               | --- 503b <- 504b <- 505b <- 506b <- 507b
```

we would update to this:
```
              LCA
               V
500 <- 501 <- 502 <- 503a <- 504  <- 505
               ^
               | --- 503b <- 504b <- 505b <- 506b <- 507b
                                                      ^
                                                     BEST
```

--------------------------
