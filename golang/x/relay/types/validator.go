package types

import "github.com/summa-tx/relays/proto"

// FilledRequestInfo contains information about what input and/or output satisfied the request
type FilledRequestInfo struct {
	InputIndex  uint32
	OutputIndex uint32
	ID          RequestID
}

// FromProto populates a FilledRequestInfo from a protobuf
func filledRequestInfoFromProto(m *proto.FilledRequestInfo) (FilledRequestInfo, error) {
	var filledRequest FilledRequestInfo

	id, err := bufToRequestID(q.ID)
	if err != nil {
		return filledRequest, err
	}

	filledRequest.InputIndex = m.InputIndex
	filledRequest.OutputIndex = m.OutputIndex
	filledRequest.ID = id

	return filledRequest, nil
}

// FilledRequests contains a proof that satisfies one or more requests
type FilledRequests struct {
	Proof  SPVProof
	Filled []FilledRequestInfo
}

// FromProto populates FilledRequests from a protobuf
func filledRequestsFromProto(m *proto.FilledRequests) (FilledRequests, error) {
	filledRequests := make([]FilledRequests, len(m))
	for	i, f := range m {
		filledRequest, err := filledRequestInfoFromProto(f)
		filledRequests[i] = filledRequest
		if err != nil {
			return nil, err
		}
	}

	return filledRequests, nil
}