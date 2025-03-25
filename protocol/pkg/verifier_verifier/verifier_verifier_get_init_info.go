package verifier_verifier

type VVInitInfoOperationRequest struct {
	RequestId int64 `json:"requestId"`
}
type VVInitInfoOperationResponse struct {
	RequestId           int64                         `json:"requestId"`
	OperationError      string                        `json:"operationError"`
	CurrentVerifierInfo VVInitInfoStructureVerifier   `json:"currentVerifierInfo"`
	VerifiersList       []VVInitInfoStructureVerifier `json:"verifierList"`
	GatewaysList        []VVInitInfoStructureGateway  `json:"gatewayList"`
}
type VVInitInfoStructureVerifier struct {
	VerifierPublicKeySignature string  `json:"verifierPublicKeySigniture"`
	VerifierPublicKeyKem       string  `json:"verifierPublicKeyKem"`
	VerifierIpAddress          string  `json:"verifierIpAddress"`
	VerifierPort               string  `json:"verifierPort"`
	SigScheme                  string  `json:"sigScheme"`
	TrustScore                 float64 `json:"trustScore"`
	IsInCommittee              bool    `json:"isInCommittee"`
}

type VVInitInfoStructureGateway struct {
	GatewayPublicKeySignature string `json:"gatewayPublicKeySigniture"`
	GatewayPublicKeyKem       string `json:"gatewayPublicKeyKem"`
	GatewayIpAddress          string `json:"gatewayIpAddress"`
	GatewayPort               string `json:"gatewayPort"`
	KemScheme                 string `json:"kemScheme"`
	SigScheme                 string `json:"sigScheme"`
}
