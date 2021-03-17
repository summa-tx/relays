package relay

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/summa-tx/relays/golang/x/relay/types"
)

func getGenesisHeaders() (types.BitcoinHeader, []types.BitcoinHeader) {
	// get path to root directory
	path, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	// if running this function from cli_test directory do not include it in golang path
	path = strings.TrimSuffix(path, "/cli_test")

	headerJSON, jsonErr := ioutil.ReadFile("/" + path + "/scripts/json_data/genesis.json")
	if jsonErr != nil {
		panic("could not retreive data in gen_state: " + jsonErr.Error())
	}

	var genesisHeaders []types.BitcoinHeader
	err = json.Unmarshal([]byte(headerJSON), &genesisHeaders)
	if err != nil {
		panic("bad json in gen_state: " + err.Error())
	}
	return genesisHeaders[0], genesisHeaders[1:]
}
