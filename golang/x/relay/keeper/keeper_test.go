package keeper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LinkTest struct{}

type ChainTest struct{}

type HeaderTest struct{}

type KeeperTestCases struct {
	LinkTestCases   []LinkTest   `json:"link"`
	HeaderTestCases []HeaderTest `json:"header"`
	ChainTestCases  []ChainTest  `json:"chain"`
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

func (suite *KeeperSuite) TestGetPrefixStore() {

}
