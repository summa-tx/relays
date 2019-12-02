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

type NamedCase interface {
	Name() string
}

type Case struct {
	NamedCase
	Comment string `json:"comment"`
}

type LinkTest struct{}

type ChainTest struct{}

type IngestCase struct {
	Case
	Headers   []types.BitcoinHeader `json:"headers"`
	Anchor    types.BitcoinHeader   `json:"anchor"`
	Internal  bool                  `json:"internal"`
	IsMainnet bool                  `json:"isMainnet"`
	Output    sdk.CodeType          `json:"output"`
}

type DiffChangeCase struct {
	Case
	Headers        []types.BitcoinHeader `json:"headers"`
	PrevEpochStart types.BitcoinHeader   `json:"prevEpochStart"`
	Anchor         types.BitcoinHeader   `json:"anchor"`
	Output         sdk.CodeType          `json:"output"`
}

type CompareCase struct {
	Case
	Full      sdk.Uint `json:"full"`
	Truncated sdk.Uint `json:"truncated"`
	Output    bool     `json:"output"`
}

type HeaderTestCases struct {
	ValidateDiffChange []DiffChangeCase `json:"validateDifficultyChange"`
	ValidateChain      []IngestCase     `json:"validateHeaderChain"`
	CompareTargets     []CompareCase    `json:"compareTargets"`
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

func (c Case) Name() string {
	return c.Comment
}

func logIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func logIfTestCaseError(tc NamedCase, err sdk.Error) {
	if err != nil {
		log.Printf("Unexpected Error\nIn case: %s\n%s\n", tc.Name(), err.Error())
	}
}

// Runs the whole test suite
func TestKeeper(t *testing.T) {
	jsonFile, err := os.Open("../../../../testVectors.json")
	defer jsonFile.Close()
	logIfError(err)

	byteValue, err := ioutil.ReadAll(jsonFile)
	logIfError(err)

	var fixtures KeeperTestCases
	err = json.Unmarshal([]byte(byteValue), &fixtures)
	logIfError(err)

	keeperSuite := new(KeeperSuite)
	keeperSuite.Fixtures = fixtures

	suite.Run(t, keeperSuite)
}
