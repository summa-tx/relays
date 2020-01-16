package types

type FilledRequestInfo struct {
	InputIndex  uint8  `json:"inputIndex"`
	OutputIndex uint8  `json:"outputIndex"`
	ID          uint64 `json:"id"`
}

type FilledRequests struct {
	Proof    SPVProof            `json:"proof"`
	Requests []FilledRequestInfo `json:"requests"`
}
