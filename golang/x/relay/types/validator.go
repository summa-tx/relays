package types

import "github.com/summa-tx/relays/proto"

// FilledRequestInfo contains information about what input and/or output satisfied the request
type FilledRequestInfo struct {
	InputIndex  uint32
	OutputIndex uint32
	ID          RequestID
}

func (f *FilledRequestInfo) FromProto(m *proto.FilledRequestInfo) (error) {
	// TODO
	return nil
}

// FilledRequests contains a proof that satisfies one or more requests
type FilledRequests struct {
	Proof  SPVProof
	Filled []FilledRequestInfo
}

func (f *FilledRequests) FromProto(m *proto.FilledRequests) (error) {
	// TODO
	return nil
}