package keeper

import "github.com/summa-tx/relays/golang/x/relay/types"

func concatenateHeaders(array []types.BitcoinHeader) []byte {
	var current []byte
	for i := range array {
		current = append(current, array[i].Raw[:]...)
	}
	return current
}
