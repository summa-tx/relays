package relay

import (
	"encoding/json"
	"io/ioutil"
)

func getGenesisHeaders() (BitcoinHeader, []BitcoinHeader) {
	headerJSON, jsonErr := ioutil.ReadFile("scripts/json_data/genesis.json")
	if jsonErr != nil {
		panic("could not retreive data in gen_state: " + jsonErr.Error())
	}

	var genesisHeaders []BitcoinHeader
	err := json.Unmarshal([]byte(headerJSON), &genesisHeaders)
	if err != nil {
		panic("bad json in gen_state: " + err.Error())
	}
	return genesisHeaders[0], genesisHeaders[1:]
}
