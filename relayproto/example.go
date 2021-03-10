package relayproto

import "github.com/summa-tx/relays/golang/x/relay/types"

var _ QueryServer;

func (q *QueryParamsIsAncestor) Translate() (*types.QueryParamsIsAncestor) {
	return nil
}