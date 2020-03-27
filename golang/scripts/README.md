# Build and Run App

## Setup
If you have never used the `go mod` before, you must add some parameters to your environment.

```bash
mkdir -p $HOME/go/bin
echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
source ~/.bash_profile
```

Now, you can install and run the application.

```bash
# Clone repository
git clone https://github.com/summa-tx/relays.git
cd relays/golang

# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands:
relayd help
relaycli help
```
## Running the CLI
To run the CLI for manual testing, you can run `./scripts/init_chain.sh` to initialize a new chain.<br><br>
All chain related data lives in `scripts/json_data`. Edit `scripts/json_data/genesis.json` to generate a customized genesis state. This JSON must be a list of block headers pertaining to one epoch. The first header must be the first block of the epoch. The remaining headers must be ordered headers beginning at any height in the epoch.
```bash
# Set the executable rights if not done already
chmod +x scripts/init_chain.sh

#run script
./scripts/init_chain.sh
```
Open up a new terminal tab in the same directory to begin interacting with the chain. As per the setup script, you can now interact via username/password `me / 12345678` such that when submitting transactions using flag `--from me` when prompted for the password enter: `12345678`

### Query CLI
Querying neither requires the `--from` flag nor a password.
```bash
# Retrieve the first digest of the relay
relaycli query relay getrelaygenesis

# Retrieve the best known digest
relaycli query relay getlastreorglca

# List other query options
relaycli query relay
```

### Transact with CLI
Transactions require the `--from` flag and password.<br><br>
JSON parameters can be accepted as either raw json or json files. including the `--inputfile` flag will interpret all json parameters as json files from directory `scripts/json_data` <br><br>
use the flag ` --broadcast-mode block` to get errors synchronously upon transactions. Otherwise errors could get swallowed resulting in false positive success <br><br>
Here are some transactions and queries you can run upon initializing the chain with the default genesis state:

```bash
# Add the following bitcoin headers which also correspond with a difficulty change in the bitcoin change
relaycli tx relay ingestdiffchange ef8248820b277b542ac2a726ccd293e8f2a3ea24c1fe04000000000000000000  0_new_difficulty.json --inputfile --from me --broadcast-mode block

# Submit Proof Request
relaycli tx relay newrequest 0x 0x17a91423737cd98bb6b2da5a11bcd82e5de36591d69f9f87 0  1  --broadcast-mode block --from me

# Check whether given proof is valid: It will not because block with transaction has not been ingested yet
relaycli query relay checkproof 1_check_proof.json --inputfile

# Ingest new headers to relay (without any change in difficulty)
relaycli tx relay ingestheaders 2_ingest_headers.json  --from me  --inputfile --broadcast-mode block

# Check whether given proof is valid: It will will be valid with new headers from previous tx
relaycli query relay checkproof 1_check_proof.json --inputfile

# Provide valid proof that fulfils a proof request
relaycli tx relay provideproof 1_check_proof.json 3_filled_requests.json  --from me  --inputfile --broadcast-mode block

# Ingest remaining headers to relay (without any change in difficulty)
relaycli tx relay ingestheaders 4_ingest_headers.json  --from me  --inputfile --broadcast-mode block

```
