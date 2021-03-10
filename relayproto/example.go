package relayproto

import (
	"fmt"

	"github.com/summa-tx/bitcoin-spv/golang/btcspv"

	"github.com/summa-tx/relays/golang/x/relay/types"
)

var _ QueryServer;


func bufToH256(buf []byte) (btcspv.Hash256Digest, error) {
	var h btcspv.Hash256Digest;
	if len(buf) != 32 {
		return h, fmt.Errorf("Expected 32 bytes, got %d bytes", len(buf))
	}

	copy(h[:], buf)

	return h, nil
}

// translate is an example of how you might translate a protobuf struct to
// the existing structs that the relay expects

// receiver is the Protobuf-style.
// Output is the existing relay style
func (q *QueryParamsIsAncestor) translate() (types.QueryParamsIsAncestor, error) {
	var query types.QueryParamsIsAncestor;

	// Do any parsing/translation work
	digest, err := bufToH256(q.DigestLE)
	if err != nil {
		return query, err
	}

	prospect, err := bufToH256(q.ProspectiveAncestor)
	if err != nil {
		return query, err
	}

	query.DigestLE = digest
	query.ProspectiveAncestor = prospect
	query.Limit = q.Limit

	return query, nil
}