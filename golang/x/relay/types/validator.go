package types

type FilledRequestInfo struct {
	InputIndex  uint8     `json:"inputIndex"`
	OutputIndex uint8     `json:"outputIndex"`
	ID          RequestID `json:"id"`
}

type FilledRequests struct {
	Proof  SPVProof            `json:"proof"`
	Filled []FilledRequestInfo `json:"requests"`
}
