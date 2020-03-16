package types

// FilledRequestInfo contains information about what input and/or output satisfied the request
type FilledRequestInfo struct {
	InputIndex  uint32    `json:"inputIndex"`
	OutputIndex uint32    `json:"outputIndex"`
	ID          RequestID `json:"id"`
}

// FilledRequests contains a proof that satisfies one or more requests
type FilledRequests struct {
	Proof  SPVProof            `json:"proof"`
	Filled []FilledRequestInfo `json:"requests"`
}

// NewFilledRequests instantiates a FilledRequests
func NewFilledRequests(proof SPVProof, filled []FilledRequestInfo) FilledRequests {
	return FilledRequests{
		proof,
		filled,
	}
}
