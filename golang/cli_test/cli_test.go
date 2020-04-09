package clitest

import (
	"testing"
	"encoding/hex"
	"github.com/stretchr/testify/suite"
)

type UtilsSuite struct {
	suite.Suite
	TestData TestData
}

// Runs the whole test suite
func TestRelay(t *testing.T) {

	utilsSuite := new(UtilsSuite)
	utilsSuite.TestData = GrabTestData(t)

	suite.Run(t, utilsSuite)
}

func (suite *UtilsSuite) TestRelayCLIIsAncestor() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders
	newDiffHeaders := suite.TestData.NewDiffHeaders

	// Initialize CHain
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)

	// define param values
	fooAddr        := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	ancestor       := hex.EncodeToString(newDiffHeaders[0].HashLE[:])
	digest         := hex.EncodeToString(newDiffHeaders[1].HashLE[:])
	limit          := "5"

	// must ingest headers in order to perform query
	success, _, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	suite.True(success, stderr)

	// query chain for ancestor
	isancestor := f.QueryIsAncestor(digest, ancestor, limit)

	// Condition
	expected := true
	actual := isancestor.Res
	suite.Equal(expected, actual)

	//Cleanup
	f.Cleanup()
}

func (suite *UtilsSuite) TestRelayCLIGetRelayGenesis() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders

	// Query Chain for Actual Value
    f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	genesisRelay := f.QueryGetRelayGenesis(fooAddr)
	actual := genesisRelay.Res

	// Condition
	expected := genesisHeaders[1].HashLE
	suite.Equal(expected, actual)

	//Cleanup
	f.Cleanup()
}

func (suite *UtilsSuite) TestRelayCLIGetLastReorgLCA() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders

	// Query Chain for Actual Value
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	lastReorgLCA := f.QueryGetLastReorgLCA(fooAddr)
	actual := lastReorgLCA.Res

	// Condition
	expected := genesisHeaders[1].HashLE
	suite.Equal(expected, actual)

	//Cleanup
	f.Cleanup()
}

func (suite *UtilsSuite) TestRelayCLIGetBestDigest() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders

	// Query Chain for Actual Value
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr := f.KeyAddress(keyFoo)
	bestDigest := f.QueryGetBestDigest(fooAddr)

	// Condition
	actual := bestDigest.Res
	expected := genesisHeaders[1].HashLE
	suite.Equal(expected, actual)

	//Cleanup
	f.Cleanup()
}

func (suite *UtilsSuite) TestRelayCLIQueryFindAncestor() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders
	newDiffHeaders := suite.TestData.NewDiffHeaders

	// Initialize chain
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)

	// define paramater values
	fooAddr        := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	digest         := hex.EncodeToString(newDiffHeaders[1].HashLE[:])
	invalidOffset  := "5"
	validOffset    := "1"

	// ingest headers
	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	// Require findancestor fails if ancestor does not exist on relay
	f.QueryFindAncestorInvalid("could not find ancestor", digest, invalidOffset)

	// Require findancestor returns ancestor if valid query
	findancestor := f.QueryFindAncestor(digest, validOffset)
	expected := hex.EncodeToString(newDiffHeaders[0].HashLE[:])
	actual := hex.EncodeToString(findancestor.Res[:])
	suite.Equal(expected, actual)
}

func (suite *UtilsSuite) TestRelayCLIIsMostRecentCommonAncestor() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders
	newDiffHeaders := suite.TestData.NewDiffHeaders

	// Query Chain for Actual Value
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)

	// must ingest headers in order to perform query
	fooAddr            := f.KeyAddress(keyFoo)
	prevEpochStart     := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	success, _, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	suite.True(success, stderr)

	//perform query
	ancestor := hex.EncodeToString(newDiffHeaders[0].HashLE[:])
	left     := hex.EncodeToString(newDiffHeaders[0].HashLE[:])
	right    := hex.EncodeToString(newDiffHeaders[1].HashLE[:])
	limit    := "3"
	ismostrecentcommonancestor := f.QueryIsMostRecentCommonAncestor(ancestor, left, right, limit)

	// Condition
	expected := true
	actual   := ismostrecentcommonancestor.Res
	suite.Equal(expected, actual)

	// Require query returns error if invalid
	// ancestor = hex.EncodeToString(newDiffHeaders[1].HashLE[:])
	// ismostrecentcommonancestor = f.QueryIsMostRecentCommonAncestor(ancestor, left, right, limit)

	//Cleanup
	f.Cleanup()
}

func (suite *UtilsSuite) TestRelayCLIQueryHeaviestFromAncestor() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders
	newDiffHeaders := suite.TestData.NewDiffHeaders

	// Transact with Chain for Actual Value
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)

	// define paramteer values
	fooAddr        := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	ancestor       := hex.EncodeToString(genesisHeaders[1].HashLE[:])
	currentBest    := hex.EncodeToString(genesisHeaders[1].HashLE[:])
	validNewBest   := hex.EncodeToString(newDiffHeaders[1].HashLE[:])
	invalidNewBest := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	limit          := "10"

	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	// Query chain
	heaviestfromancestor := f.QueryHeaviestFromAncestor(ancestor, currentBest, validNewBest, limit)

	// Condition (query returns newBest)
	actual := hex.EncodeToString(heaviestfromancestor.Res[:])
	suite.Equal(validNewBest, actual)

	// Require heaviestfromancestor fails with invalid params
	f.QueryHeaviestFromAncestorInvalid("could not determine", ancestor, currentBest, invalidNewBest, limit)
}

func (suite *UtilsSuite) TestRelayCLIQueryCheckProof() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders

	// Initialize chain
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)
	fooAddr        := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])

	// Ingest headers
	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	// Require checkproof fails without headers associated with proof
	checkProof := f.QueryCheckProof("1_check_proof.json", "--inputfile")
	suite.Equal(false, checkProof.Valid)

	// Ingest associated header
	success, stdout, stderr = f.TxIngestHeaders(fooAddr, "2_ingest_headers.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	// Require proof is valid when associated header exists with valid transaction
	checkProof = f.QueryCheckProof("1_check_proof.json", "--inputfile")
	suite.Equal(true, checkProof.Valid)
}

func (suite *UtilsSuite) TestRelayCLITXIngestHeaders() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders

	// Initialize chain
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)

	// define parameter valuse
	fooAddr         := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])

	//Ingest Headers
	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	// Require successful IngestDiffChange
	success, stdout, stderr = f.TxIngestHeaders(fooAddr, "2_ingest_headers.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	//Cleanup
	f.Cleanup()
}

func (suite *UtilsSuite) TestRelayCLITXIngestDiffChange() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders

	// Initialize chain
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)

	// Define parameter values
	fooAddr        := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])

	// Require successful IngestDiffChange
	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	//Cleanup
	f.Cleanup()
}

func (suite *UtilsSuite) TestRelayCLITXProvideProof() {
	suite.T().Parallel()
	genesisHeaders := suite.TestData.GenesisHeaders

	// Initialize chain
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)

	// Define parameter values
	fooAddr        := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])

	// Ingest Headers w/ Diff Change
	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	// Ingest New Headers
	success, stdout, stderr = f.TxIngestHeaders(fooAddr, "2_ingest_headers.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	// require checkproof fails given invalid proof requests
	_, stdout, _ = f.TxProvideProof(fooAddr, "1_check_proof.json", "3_filled_requests.json", "--inputfile -y")
	suite.Contains(stdout, `"Request not found`)

	// submit proof request
	spends   := "0x"
	pays     := "0x17a91423737cd98bb6b2da5a11bcd82e5de36591d69f9f87"
	value    := "0"
	numConfs := "1"
	success, stdout, stderr = f.TxNewRequest(fooAddr, spends, pays, value, numConfs, "-y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	// checkproof succeeds given valid proof and requests
	success, stdout, stderr = f.TxProvideProof(fooAddr, "1_check_proof.json", "3_filled_requests.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	//Cleanup
	f.Cleanup()
}

func (suite *UtilsSuite) TestRelayCLITxMarkNewHeaviest() {
	suite.T().Parallel()

	genesisHeaders := suite.TestData.GenesisHeaders
	newDiffHeaders := suite.TestData.NewDiffHeaders

	// Initialize chain
	f := InitFixtures(suite.T())
	proc := f.RelayDStart()
	defer proc.Stop(false)

	// Define parameter values
	fooAddr        := f.KeyAddress(keyFoo)
	prevEpochStart := hex.EncodeToString(genesisHeaders[0].HashLE[:])
	ancestor       := hex.EncodeToString(genesisHeaders[1].HashLE[:])
	bestKnown      := hex.EncodeToString(genesisHeaders[1].Raw[:])
	newBest        := hex.EncodeToString(newDiffHeaders[1].Raw[:])
	limit          := "10"

	// Ingest new headers
	success, stdout, stderr := f.TxIngestDiffChange(fooAddr, prevEpochStart, "0_new_difficulty.json", "--inputfile -y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	// Mark new heaviest digest
	success, stdout, stderr = f.TxMarkNewHeaviest(fooAddr, ancestor, bestKnown, newBest, limit, "-y")
	suite.True(success, stderr)
	suite.Contains(stdout, `"success":true`)

	bestDigest := f.QueryGetBestDigest(fooAddr)

	// Condition
	expected := hex.EncodeToString(newDiffHeaders[1].HashLE[:])
	actual   := hex.EncodeToString(bestDigest.Res[:])
	suite.Equal(expected, actual)

	//Cleanup
	f.Cleanup()
}
