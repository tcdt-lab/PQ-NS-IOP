package gateway_verifier

type VerificationRequest struct {
	RequestId    int64  `json:"requestId"`
	PublicInputs string `json:"publicInputs"`
	Proof        string `json:"proof"`
}

type VerificationResponse struct {
	RequestId          int64 `json:"requestId"`
	VerificationResult bool  `json:"verificationResult"`
}
