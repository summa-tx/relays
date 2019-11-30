package keeper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/summa-tx/relays/golang/x/relay/types"
)

type LinkTest struct{}

type ChainTest struct{}

type IngestCase struct {
	Headers   []types.BitcoinHeader `json:"headers"`
	Anchor    types.BitcoinHeader   `json:"anchor"`
	Internal  bool                  `json:"internal"`
	IsMainnet bool                  `json:"isMainnet"`
	Output    sdk.CodeType          `json:"output"`
}

type DiffChangeCase struct {
	Headers        []types.BitcoinHeader `json:"headers"`
	PrevEpochStart types.BitcoinHeader   `json:"prevEpochStart"`
	Anchor         types.BitcoinHeader   `json:"anchor"`
	Output         sdk.CodeType          `json:"output"`
}

type HeaderTestCases struct {
	ValidateDiffChange []DiffChangeCase `json:"validateDifficultyChange"`
	ValidateChain      []IngestCase     `json:"validateHeaderChain"`
}

type KeeperTestCases struct {
	LinkTestCases   []LinkTest      `json:"link"`
	HeaderTestCases HeaderTestCases `json:"header"`
	ChainTestCases  []ChainTest     `json:"chain"`
}

type KeeperSuite struct {
	suite.Suite
	Fixtures KeeperTestCases
}

func logIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Runs the whole test suite
func TestKeeper(t *testing.T) {
	jsonFile, err := os.Open("../../../../testVectors.json")
	defer jsonFile.Close()
	logIfErr(err)

	byteValue, err := ioutil.ReadAll(jsonFile)
	logIfErr(err)

	var fixtures KeeperTestCases
	json.Unmarshal([]byte(byteValue), &fixtures)

	keeperSuite := new(KeeperSuite)
	keeperSuite.Fixtures = fixtures

	suite.Run(t, keeperSuite)
}
