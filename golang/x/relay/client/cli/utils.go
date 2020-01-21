package cli

import (
	"encoding/binary"
	"strconv"

	"github.com/summa-tx/relays/golang/x/relay/types"
)

func requestIDFromString(s string) (types.RequestID, error) {
	var idBytes types.RequestID
	var err error
	if s[:2] == "0x" {
		idBytes, err = types.RequestIDFromHex(s)
		if err != nil {
			return types.RequestID{}, err
		}
	} else {
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return types.RequestID{}, err
		}

		// convert to bytes
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, id)
		idBytes, err = types.NewRequestID(b)
		if err != nil {
			return types.RequestID{}, err
		}
	}
	return idBytes, err
}
