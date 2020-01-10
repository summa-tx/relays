package types

type ProofRequest struct {
	Spends      Hash256Digest `json:"spends"`
	Pays        Hash256Digest `json:"pays"`
	PaysValue   uint64        `json:"paysValue"`
	ActiveState bool          `json:"activeState"`
	NumConfs    uint8         `json:"numConfs"`
}
