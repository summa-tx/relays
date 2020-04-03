package clitest

import (
	"testing"
	"encoding/json"
	"encoding/hex"

	"github.com/stretchr/testify/require"
	rtypes "github.com/summa-tx/relays/golang/x/relay/types"
)

func TestRelayCLIGetRelayGenesis(t *testing.T) {
	// Get Expected Value
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)
	expected := genesisHeaders[1].HashLE

	// Query Chain for Actual Value
    f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	genesisRelay := f.QueryGetRelayGenesis(fooAddr)
	actual := genesisRelay.Res

	// Condition
	require.Equal(t, expected, actual)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLIGetLastReorgLCA(t *testing.T) {
	// Get Expected Value
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)
	expected := genesisHeaders[1].HashLE

	// Query Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	lastReorgLCA := f.QueryGetLastReorgLCA(fooAddr)
	actual := lastReorgLCA.Res

	// Condition
	require.Equal(t, expected, actual)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLIGetBestDigest(t *testing.T) {
	// Get Expected Value
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)
	expected := genesisHeaders[1].HashLE

	// Query Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	bestDigest := f.QueryGetBestDigest(fooAddr)
	actual := bestDigest.Res

	// Condition
	require.Equal(t, expected, actual)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLITXIngestDiffChange(t *testing.T) {
	// Extract data for transaction
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	ancestor := hex.EncodeToString(genesisHeaders[1].HashLE[:])
	newBest := hex.EncodeToString(newDifficultyHeaders[1].HashLE[:])

	// Transact with Chain
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	success, _, sterr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	require.True(t, success)

	//Cleanup
	f.Cleanup()
}

func TestRelayCLITxMarkNewHeaviest(t *testing.T) {
	// Extract data for transaction
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)

	var newDifficultyHeaders []rtypes.BitcoinHeader
	newDiffJSON := readJSONFile(t, "0_new_difficulty")
	err = json.Unmarshal([]byte(newDiffJSON), &newDifficultyHeaders)
	require.NoError(t, err)

	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	ancestor := hex.EncodeToString(genesisHeaders[1].HashLE[:])
	bestKnown := hex.EncodeToString(genesisHeaders[1].Raw[:])
	newBest := hex.EncodeToString(newDifficultyHeaders[1].Raw[:])
	limit := "10"

	// Get Expected Value
	expected := hex.EncodeToString(newDifficultyHeaders[1].HashLE[:])

	// Query Chain for Actual Value
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)

	success, _, _ := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	require.True(t, success)

	success, _, _ = f.TxMarkNewHeaviest(fooAddr, ancestor, bestKnown, newBest, limit, "-y")
	require.True(t, success)

	bestDigest := f.QueryGetBestDigest(fooAddr)
	actual := hex.EncodeToString(bestDigest.Res[:])

	// Condition
	require.Equal(t, expected, actual)

	//Cleanup
	f.Cleanup()
}
