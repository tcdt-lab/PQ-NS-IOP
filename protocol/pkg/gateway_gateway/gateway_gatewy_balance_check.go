package gateway_gateway

type BalanceCheckRequest struct {
	RequestId        int64 `json:"requestId"`
	BalanceThreshold int64 `json:"balanceThreshould"`
}

type BalanceCheckResponse struct {
	RequestId    int64  `json:"requestId"`
	Proof        string `json:"proof"`
	PublicInputs string `json:"publicInputs"`
}
