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
To run the CLI for manual testing, you can run `./scripts/init_chain.sh` to initialize a new chain.
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
use the flag ` --broadcast-mode block` to get errors synchronously upon transactions. Otherwise errors could get swallowed resulting in false positive success 
```bash
# Add more recent bitcoin headers to relay
relaycli tx relay ingestheaders headers_init.json --inputfile --from me --broadcast-mode block

# Submit Proof Requests
relaycli tx relay ingestheaders request.json --inputfile --from me --broadcast-mode block
```
