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

	id, err := bufToRequestID(m.ID)
	if err != nil {
		return filledRequest, err
	}

	filledRequest.InputIndex = m.InputIndex
	filledRequest.OutputIndex = m.OutputIndex
	filledRequest.ID = id

	return filledRequest, nil
}

func filledRequestInfoToProto(m *FilledRequestInfo) (proto.FilledRequestInfo) {
	var filledRequest FilledRequestInfo

	filledRequest.InputIndex = m.InputIndex
	filledRequest.OutputIndex = m.OutputIndex
	filledRequest.ID = []byte{m.ID}

	return filledRequest
}

// FilledRequests contains a proof that satisfies one or more requests
type FilledRequests struct {
	Proof  SPVProof
	Filled []FilledRequestInfo
}

// FromProto populates FilledRequests from a protobuf
func filledRequestsFromProto(m *proto.FilledRequests) (FilledRequests, error) {
	var filledRequests FilledRequests

	proof, err := spvProofFromProto(m.Proof)
	if err != nil {
		return filledRequests, err
	}

	filled := make([]FilledRequestInfo, len(m.Filled))
	for	i, f := range m.Filled {
		filledRequest, err := filledRequestInfoFromProto(f)
		filled[i] = filledRequest
		if err != nil {
			return filledRequests, err
		}
	}

	filledRequests.Proof = proof
	filledRequests.Filled = filled

	return filledRequests, nil
}

func filledRequestsToProto(m *FilledRequests) (FilledRequests) {
	var filledRequests FilledRequests

	filled := make([]proto.FilledRequestInfo, len(m.Filled))
	for	i, f := range m.Filled {
		filledRequest := filledRequestInfoToProto(f)
		filled[i] = filledRequest
	}

	filledRequests.Proof = spvProofToProto(m.Proof)
	filledRequests.Filled = filled

	return filledRequests
}