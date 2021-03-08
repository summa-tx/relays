package clitest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	app "github.com/summa-tx/relays/golang"
	rtypes "github.com/summa-tx/relays/golang/x/relay/types"

	clientkeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	keyFoo     = "foo"
	stakeDenom = "stake"
	fooDenom   = "footoken"
)

var startCoins = sdk.NewCoins(
	sdk.NewCoin(stakeDenom, sdk.TokensFromConsensusPower(1000000)),
	sdk.NewCoin(fooDenom, sdk.TokensFromConsensusPower(1000000)),
)

type TestData struct {
	GenesisHeaders []rtypes.BitcoinHeader
	NewDiffHeaders []rtypes.BitcoinHeader
	NewHeaders     []rtypes.BitcoinHeader
}

// Fixtures is used to setup the testing environment
type Fixtures struct {
	BinDir         string
	RootDir        string
	RelaydBinary   string
	RelaycliBinary string
	ChainID        string
	RPCAddr        string
	Port           string
	RelaydHome     string
	RelaycliHome   string
	P2PAddr        string
	T              *testing.T
}

//////////////////////////////////////////////////////////////////////////////////////
// Instantiation /////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

// NewFixtures creates a new instance of Fixtures with many vars set
func NewFixtures(t *testing.T) *Fixtures {
	tmpDir, err := ioutil.TempDir("", "relay_integration_"+strings.TrimPrefix(t.Name(), "TestRelay/")+"_")
	require.NoError(t, err)
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)

	p2pAddr, _, err := server.FreeTCPAddr()
	require.NoError(t, err)

	binDir := os.Getenv("GOPATH") + "/bin/"
	require.True(t, fileExists(binDir+"relayd"), "relayd binary does not exist")
	require.True(t, fileExists(binDir+"relaycli"), "relaycli binary does not exist")

	return &Fixtures{
		T:              t,
		BinDir:         binDir,
		RootDir:        tmpDir,
		RelaydBinary:   filepath.Join(binDir, "relayd"),
		RelaycliBinary: filepath.Join(binDir, "relaycli"),
		RelaydHome:     filepath.Join(tmpDir, ".relayd"),
		RelaycliHome:   filepath.Join(tmpDir, ".relaycli"),
		RPCAddr:        servAddr,
		P2PAddr:        p2pAddr,
		Port:           port,
	}
}

// InitFixtures is called at the beginning of a test  and initializes a chain
// with 1 validator.
func InitFixtures(t *testing.T) (f *Fixtures) {
	f = NewFixtures(t)

	// reset test state
	f.UnsafeResetAll()

	f.KeysAdd(keyFoo)

	// ensure that CLI output is in JSON format
	f.CLIConfig("output", "json")

	// NOTE: RelayDInit sets the ChainID
	f.RelayDInit(keyFoo)

	f.CLIConfig("chain-id", f.ChainID)
	f.CLIConfig("broadcast-mode", "block")
	f.CLIConfig("trust-node", "true")

	// start an account with tokens
	f.AddGenesisAccount(f.KeyAddress(keyFoo), startCoins)

	f.GenTx(keyFoo)
	f.CollectGenTxs()

	return f
}

func GrabTestData(t *testing.T) TestData {
	testData := TestData{}
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &testData.GenesisHeaders)
	require.NoError(t, err)

	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &testData.NewDiffHeaders)
	require.NoError(t, err)

	newHeadersJSON := readJSONFile(t, "2_ingest_headers")
	err = json.Unmarshal([]byte(newHeadersJSON), &testData.NewHeaders)
	require.NoError(t, err)

	return testData
}

//////////////////////////////////////////////////////////////////////////////////////
// Fixtures Interface ////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

// RelayDInit is relayd init
// NOTE: RelayDInit sets the ChainID for the Fixtures instance
func (f *Fixtures) RelayDInit(moniker string, flags ...string) {
	cmd := fmt.Sprintf("%s init -o --home=%s %s", f.RelaydBinary, f.RelaydHome, moniker)
	_, stderr := tests.ExecuteT(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)

	var chainID string
	var initRes map[string]json.RawMessage

	err := json.Unmarshal([]byte(stderr), &initRes)
	require.NoError(f.T, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(f.T, err)

	f.ChainID = chainID
}

// RelaydStart runs relayd start with the appropriate flags and returns a process
func (f *Fixtures) RelayDStart(flags ...string) *tests.Process {
	cmd := fmt.Sprintf("%s start --home=%s --rpc.laddr=%v --p2p.laddr=%v", f.RelaydBinary, f.RelaydHome, f.RPCAddr, f.P2PAddr)
	proc := tests.GoExecuteTWithStdout(f.T, addFlags(cmd, flags))
	tests.WaitForTMStart(f.Port)
	tests.WaitForNextNBlocksTM(1, f.Port)
	return proc
}

// Cleanup is meant to be run at the end of a test to clean up an remaining test state
func (f *Fixtures) Cleanup(dirs ...string) {
	clean := append(dirs, f.RootDir)
	for _, d := range clean {
		require.NoError(f.T, os.RemoveAll(d))
	}
}

// AddGenesisAccount is relayd add-genesis-account
func (f *Fixtures) AddGenesisAccount(address sdk.AccAddress, coins sdk.Coins, flags ...string) {
	cmd := fmt.Sprintf("%s add-genesis-account %s %s --home=%s", f.RelaydBinary, address, coins, f.RelaydHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

// GenesisFile returns the path of the generated genesis file
func (f Fixtures) GenesisFile() string {
	return filepath.Join(f.RelaydHome, "config", "genesis.json")
}

// Flags returns the flags necessary for making most CLI calls
func (f *Fixtures) Flags() string {
	return fmt.Sprintf("--home=%s --node=%s", f.RelaycliHome, f.RPCAddr)
}

// KeysAdd is relaycli keys add
func (f *Fixtures) KeysAdd(name string, flags ...string) {
	cmd := fmt.Sprintf("%s keys add --home=%s %s", f.RelaycliBinary, f.RelaycliHome, name)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// CLIConfig is relaycli config
func (f *Fixtures) CLIConfig(key, value string, flags ...string) {
	cmd := fmt.Sprintf("%s config --home=%s %s %s", f.RelaycliBinary, f.RelaycliHome, key, value)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

// CollectGenTxs is relayd collect-gentxs
func (f *Fixtures) CollectGenTxs(flags ...string) {
	cmd := fmt.Sprintf("%s collect-gentxs --home=%s", f.RelaydBinary, f.RelaydHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// KeysShow is relaycli keys show
func (f *Fixtures) KeysShow(name string, flags ...string) keys.KeyOutput {
	cmd := fmt.Sprintf("%s keys show --home=%s %s", f.RelaycliBinary,
		f.RelaycliHome, name)
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var ko keys.KeyOutput
	err := clientkeys.UnmarshalJSON([]byte(out), &ko)
	require.NoError(f.T, err)
	return ko
}

// UnsafeResetAll is relayd unsafe-reset-all
func (f *Fixtures) UnsafeResetAll(flags ...string) {
	cmd := fmt.Sprintf("%s --home=%s unsafe-reset-all", f.RelaydBinary, f.RelaydHome)
	executeWrite(f.T, addFlags(cmd, flags))
	err := os.RemoveAll(filepath.Join(f.RelaydHome, "config", "gentx"))
	require.NoError(f.T, err)
}

// KeyAddress returns the SDK account address from the key
func (f *Fixtures) KeyAddress(name string) sdk.AccAddress {
	ko := f.KeysShow(name)
	accAddr, err := sdk.AccAddressFromBech32(ko.Address)
	require.NoError(f.T, err)
	return accAddr
}

// GenTx is relayd gentx
func (f *Fixtures) GenTx(name string, flags ...string) {
	cmd := fmt.Sprintf("%s gentx --name=%s --home=%s --home-client=%s", f.RelaydBinary, name, f.RelaydHome, f.RelaycliHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

//////////////////////////////////////////////////////////////////////////////////////
// CLI Queries ///////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

// QueryGetRelayGenesis returns the relay genesis block Hash
func (f *Fixtures) QueryGetRelayGenesis(delAddr sdk.AccAddress) rtypes.QueryResGetRelayGenesis {
	cmd := fmt.Sprintf("%s query relay getrelaygenesis %s %s", f.RelaycliBinary, delAddr, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var relaygenesis rtypes.QueryResGetRelayGenesis
	err := cdc.UnmarshalJSON([]byte(res), &relaygenesis)
	require.NoError(f.T, err)
	return relaygenesis
}

// QueryGetLastReorgLCA returns Last Common Anscestor
func (f *Fixtures) QueryGetLastReorgLCA(delAddr sdk.AccAddress) rtypes.QueryResGetLastReorgLCA {
	cmd := fmt.Sprintf("%s query relay getlastreorglca %s %s", f.RelaycliBinary, delAddr, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var lastreorglca rtypes.QueryResGetLastReorgLCA
	err := cdc.UnmarshalJSON([]byte(res), &lastreorglca)
	require.NoError(f.T, err)
	return lastreorglca
}

// QueryGetBestDigest returns the Best Known Digest
func (f *Fixtures) QueryGetBestDigest(delAddr sdk.AccAddress) rtypes.QueryResGetBestDigest {
	cmd := fmt.Sprintf("%s query relay getbestdigest %s %s", f.RelaycliBinary, delAddr, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var bestknowndigest rtypes.QueryResGetBestDigest
	err := cdc.UnmarshalJSON([]byte(res), &bestknowndigest)
	require.NoError(f.T, err)
	return bestknowndigest
}

// QueryIsAncestor returns the Boolean
func (f *Fixtures) QueryIsAncestor(digest, ancestor, limit string) rtypes.QueryResIsAncestor {
	cmd := fmt.Sprintf("%s query relay isancestor %s %s %s %s", f.RelaycliBinary, digest, ancestor, limit, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var isancestor rtypes.QueryResIsAncestor
	err := cdc.UnmarshalJSON([]byte(res), &isancestor)
	require.NoError(f.T, err)
	return isancestor
}

// QueryFindAncestor returns the Boolean
func (f *Fixtures) QueryFindAncestor(digest, offset string) rtypes.QueryResFindAncestor {
	cmd := fmt.Sprintf("%s query relay findancestor %s %s %s", f.RelaycliBinary, digest, offset, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var findancestor rtypes.QueryResFindAncestor
	err := cdc.UnmarshalJSON([]byte(res), &findancestor)
	require.NoError(f.T, err)
	return findancestor
}

// QueryFindAncestorInvalid require proper response for invalid query
func (f *Fixtures) QueryFindAncestorInvalid(errStr, digest, offset string) {
	cmd := fmt.Sprintf("%s query relay findancestor %s %s %s", f.RelaycliBinary, digest, offset, f.Flags())
	res, _ := tests.ExecuteT(f.T, cmd, "")
	require.Contains(f.T, res, errStr)
}

// QueryIsMostRecentCommonAncestor returns a Boolean
func (f *Fixtures) QueryIsMostRecentCommonAncestor(ancestor, left, right, limit string) rtypes.QueryResIsMostRecentCommonAncestor {
	cmd := fmt.Sprintf("%s query relay ismostrecentcommonancestor %s %s %s %s %s", f.RelaycliBinary, ancestor, left, right, limit, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var ismostrecentcommonancestor rtypes.QueryResIsMostRecentCommonAncestor
	err := cdc.UnmarshalJSON([]byte(res), &ismostrecentcommonancestor)
	require.NoError(f.T, err)
	return ismostrecentcommonancestor
}

// QueryHeaviestFromAncestor returns a Boolean
func (f *Fixtures) QueryHeaviestFromAncestor(ancestor, currentBest, newBest, limit string) rtypes.QueryResHeaviestFromAncestor {
	cmd := fmt.Sprintf("%s query relay heaviestfromancestor %s %s %s %s %s", f.RelaycliBinary, ancestor, currentBest, newBest, limit, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var heaviestfromancestor rtypes.QueryResHeaviestFromAncestor
	err := cdc.UnmarshalJSON([]byte(res), &heaviestfromancestor)
	require.NoError(f.T, err)
	return heaviestfromancestor
}

// QueryHeaviestFromAncestorInvalid require proper response for invalid query
func (f *Fixtures) QueryHeaviestFromAncestorInvalid(errStr, ancestor, currentBest, newBest, limit string) {
	cmd := fmt.Sprintf("%s query relay heaviestfromancestor %s %s %s %s %s", f.RelaycliBinary, ancestor, currentBest, newBest, limit, f.Flags())
	res, _ := tests.ExecuteT(f.T, cmd, "")
	require.Contains(f.T, res, errStr)
}

// QueryCheckProof returns the Boolean
func (f *Fixtures) QueryCheckProof(proof string, flags ...string) rtypes.QueryResCheckProof {
	cmd := fmt.Sprintf("%s query relay checkproof %s %s", f.RelaycliBinary, proof, f.Flags())
	res, errStr := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var checkproof rtypes.QueryResCheckProof
	err := cdc.UnmarshalJSON([]byte(res), &checkproof)
	require.NoError(f.T, err)
	return checkproof
}

// QueryCheckRequests returns the Boolean
func (f *Fixtures) QueryCheckRequests(proof, requests string, flags ...string) rtypes.QueryResCheckRequests {
	cmd := fmt.Sprintf("%s query relay checkrequests %s %s %s", f.RelaycliBinary, proof, requests, f.Flags())
	res, errStr := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var checkrequests rtypes.QueryResCheckRequests
	err := cdc.UnmarshalJSON([]byte(res), &checkrequests)
	require.NoError(f.T, err)
	return checkrequests
}

/////////////////////////////////////////////////////////////////////
// CLI Transactions /////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////

// TxIngestDiffChange is relaycli tx that ingests headers with new difficulty
func (f *Fixtures) TxIngestDiffChange(delAddr sdk.AccAddress, prevEpochStart, jsonHeaders string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx relay ingestdiffchange %s %s --from %s %s", f.RelaycliBinary, prevEpochStart, jsonHeaders, delAddr, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// TxIngestHeaders is a relaycli tx that ingests headers with same difficulty as previous headers
func (f *Fixtures) TxIngestHeaders(delAddr sdk.AccAddress, headers string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx relay ingestheaders %s --from %s %s", f.RelaycliBinary, headers, delAddr, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// TxNewRequest is a relaycli tx that submits a new Proof Request
func (f *Fixtures) TxNewRequest(delAddr sdk.AccAddress, spends, pays, value, numConfs string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx relay newrequest %s %s %s %s --from %s %s", f.RelaycliBinary, spends, pays, value, numConfs, delAddr, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// TxProvideProof is a relaycli tx that submits a new Proof Request
func (f *Fixtures) TxProvideProof(delAddr sdk.AccAddress, proof, listofrequests string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx relay provideproof %s %s --from %s %s", f.RelaycliBinary, proof, listofrequests, delAddr, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

// TxMarkNewHeaviest returns Last Common Anscestor
func (f *Fixtures) TxMarkNewHeaviest(delAddr sdk.AccAddress, ancestor, currentBest, newBest, limit string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx relay marknewheaviest %s %s %s %s --from %s %s", f.RelaycliBinary, ancestor, currentBest, newBest, limit, delAddr, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

//////////////////////////////////////////////////////////////////////////////////////
// Executors /////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

func executeWriteCheckErr(t *testing.T, cmdStr string, writes ...string) {
	require.True(t, executeWrite(t, cmdStr, writes...))
}

func executeWrite(t *testing.T, cmdStr string, writes ...string) (exitSuccess bool) {
	exitSuccess, _, _ = executeWriteRetStdStreams(t, cmdStr, writes...)
	return
}

func executeWriteRetStdStreams(t *testing.T, cmdStr string, writes ...string) (bool, string, string) {
	proc := tests.GoExecuteT(t, cmdStr)

	// Enables use of interactive commands
	for _, write := range writes {
		_, err := proc.StdinPipe.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}

	// Read both stdout and stderr from the process
	stdout, stderr, err := proc.ReadAll()
	if err != nil {
		fmt.Println("Err on proc.ReadAll()", err, cmdStr)
	}

	// Log output.
	if len(stdout) > 0 {
		t.Log("Stdout:", string(stdout))
	}
	if len(stderr) > 0 {
		t.Log("Stderr:", string(stderr))
	}

	// Wait for process to exit
	proc.Wait()
	// Return succes, stdout, stderr
	return proc.ExitState.Success(), string(stdout), string(stderr)
}

//////////////////////////////////////////////////////////////////////////////////////
// utils /////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
func addFlags(cmd string, flags []string) string {
	for _, f := range flags {
		cmd += " " + f
	}
	return strings.TrimSpace(cmd)
}

func readJSONFile(t *testing.T, filename string) []byte {
	headerJSON, jsonErr := ioutil.ReadFile("../scripts/json_data/" + filename + ".json")
	require.NoError(t, jsonErr)
	return headerJSON
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
