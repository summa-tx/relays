package types

type FilledRequestInfo struct {
	InputIndex  uint32     `json:"inputIndex"`
	OutputIndex uint32     `json:"outputIndex"`
	ID          RequestID  `json:"id"`
}

type FilledRequests struct {
	Proof  SPVProof            `json:"proof"`
	Filled []FilledRequestInfo `json:"requests"`
}
