# cosmos-relay-dashboard

## Description

The dashboard displays the Cosmos Relay chain data and verifies it against an external Bitcoin explorer (currently BlockStream).

--------------------------

## Getting Started

### Start Cosmos Relay

The dashboard connects with a locally run `relay`.

1. If you don't have Go installed, install Go.
2. If you haven't used the `go mod` before, add this to your environment:

```bash
$ mkdir -p $HOME/go/bin
$ echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
$ echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
$ source ~/.bash_profile
```

 > *Troubleshooting tip*
 >
 > If, after following steps 3 and 4 below, you are not able to successfully run `make install` or `make init` then try replacing the above lines with the following:
 >
 >```bash
 > $ export GOPATH=$HOME/go
 > $ export PATH=$GOPATH/bin:$PATH
 > ```
 >
 > Don't forget to run:
 > ```bash
 > $ source ~/.bash_profile
 > ```
 >
 > You may even need to restart your terminal.

3. Make sure you are in the `relays/golang` directory (one level up from here) and install the app into your `$GOBIN`.

```bash
$ make install
```

4. Initialize a new chain for testing.

```bash
$ make init
```

5. In the same folder, but in another terminal window, run the REST routes `rest-server`. This will make the relay application REST routes available on `http://localhost:1317`.

```bash
$ relaycli rest-server --chain-id relay
```

All routes are at `/relay/${route}`. For a list of available routes, see the golang README located at `relays/golang/README.md`.

[Relay Chain Instructions](https://github.com/summa-tx/relays/blob/master/golang/scripts/README.md).

### Dashboard

1. Install dependencies (`/relay/golang/dashboard`).

```base
$ npm install
```

2. Start dashboard.

```bash
$ npm run serve
```

View at http://localhost:8080 in your browser.

--------------------------

## Development

### Set Environment Variables

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

--------------------------

## Dashboard Overview: How Things Work

There are 2 sources used for the dashboard, the relay and an external source.  The Best Known Digest and Last Reorg Common Ancestor are polled every 2 minutes from the relay.  Information from the external source is polled every 3 minutes.

### Current Block

The user wants to know about new headers. In order to do that, we:

1. Get the best tip (most recent block height) from an external source.
2. Display the height, hash, and timestamp.

### Best Known Digest

This is most current and best block the relay knows about. It is updated approximately every 5 blocks, and will be behind the newest header.

1. Poll `/relay/getbestdigest`.
2. Store digest and display the height, hash, and timestamp.

This is conceptually equivalent to Github tags.

### Last (Reorg) Common Ancestor (LCA)

This is the latest block that is in the history of both the current best known digest, and the previous best known digest.

1. Poll `/relay/getlastreorglca`.
2. Store LCA and display the height, hash, and timestamp.

### Health Checks and Verification

The dashboard keeps track of and displays the following:

* **lastComms**: When was the last successful communication made?
  * **lastComms.relay** - Last successful communication from the relay.
  * **lastComms.external** - Last successful communication from the external source.

* **blockDifference**: What is the difference in blocks between the relay and the external source?

Health pulses are displayed as `TIME in MINUTES ago`.

### Networks Names

Displays the network for the relay and for the external source. For these purposes, the relay will always be a locally run relay, and the external source shows mainnet blocks.

--------------------------

## Relay

The following is mainly for informational purposes, rather than development.

### Relay updates

The relay is updated ~every 5 blocks.

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
