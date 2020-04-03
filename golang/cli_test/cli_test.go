package clitest

import (
	"testing"
    "fmt"
	"encoding/json"
	"encoding/hex"

	"github.com/stretchr/testify/require"
	rtypes "github.com/summa-tx/relays/golang/x/relay/types"
)

// func TestRelayCLIGetRelayGenesis(t *testing.T) {
// 	// Get Expected Value
// 	var genesisHeaders []rtypes.BitcoinHeader
// 	genesisJSON := readJSONFile(t, "genesis")
// 	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
// 	require.NoError(t, err)
// 	expected := genesisHeaders[1].HashLE
//
// 	// Query Chain for Actual Value
//     f := InitFixtures(t)
// 	proc := f.RelayDStart()
// 	defer proc.Stop(false)
// 	fooAddr := f.KeyAddress(keyFoo)
// 	genesisRelay := f.QueryGetRelayGenesis(fooAddr)
// 	actual := genesisRelay.Res
//
// 	// Condition
// 	require.Equal(t, expected, actual)
//
// 	//Cleanup
// 	f.Cleanup()
// }
//
// func TestRelayCLIGetLastReorgLCA(t *testing.T) {
// 	// Get Expected Value
// 	var genesisHeaders []rtypes.BitcoinHeader
// 	genesisJSON := readJSONFile(t, "genesis")
// 	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
// 	require.NoError(t, err)
// 	expected := genesisHeaders[1].HashLE
//
// 	// Query Chain for Actual Value
// 	f := InitFixtures(t)
// 	proc := f.RelayDStart()
// 	defer proc.Stop(false)
// 	fooAddr := f.KeyAddress(keyFoo)
// 	lastReorgLCA := f.QueryGetLastReorgLCA(fooAddr)
// 	actual := lastReorgLCA.Res
//
// 	// Condition
// 	require.Equal(t, expected, actual)
//
// 	//Cleanup
// 	f.Cleanup()
// }


func TestRelayCLITXIngestDiffChange(t *testing.T) {
	// Extracted needed data for transaction
	var genesisHeaders []rtypes.BitcoinHeader
	genesisJSON := readJSONFile(t, "genesis")
	err := json.Unmarshal([]byte(genesisJSON), &genesisHeaders)
	require.NoError(t, err)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	// prevEpochStart := "5c2078d0388e3844fe6241723e9543074bd3a974c16611000000000000000000"

	// hexPrevEpoch := hex.EncodeToString(prevEpochStart)
	fmt.Println("HERE!!!!")
	fmt.Println(prevEpochStart)


	// Transact with Chain
	f := InitFixtures(t)
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	success, stout, sterr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	if len(stout) > 1 {
		fmt.Println("STOUT!!!!!!!!!!")
		fmt.Println(stout)
	}
	if len(sterr) > 1 {
		fmt.Println("STERR!!!!!!!!!!")
		fmt.Println(sterr)
	}
	fmt.Println("SUCCESS!!!!!!!!!!")
	fmt.Println(success)
	require.True(t, success)


}
