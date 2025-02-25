package gateway_verifier

type GatewayVerifierInitInfoOperationRequest struct {
	RequestId int64 `json:"requestId"`
}
type GatewayVerifierInitInfoOperationResponse struct {
	RequestId           int64                                      `json:"requestId"`
	OperationError      string                                     `json:"operationError"`
	CurrentVerifierInfo GatewayVerifierInitInfoStructureVerifier   `json:"currentVerifierInfo"`
	VerifiersList       []GatewayVerifierInitInfoStructureVerifier `json:"verifierList"`
	GatewaysList        []GatewayVerifierInitInfoStructureGateway  `json:"gatewayList"`
}
type GatewayVerifierInitInfoStructureVerifier struct {
	VerifierPublicKeySignature string  `json:"verifierPublicKeySigniture"`
	VerifierPublicKeyKem       string  `json:"verifierPublicKeyKem"`
	VerifierIpAddress          string  `json:"verifierIpAddress"`
	VerifierPort               string  `json:"verifierPort"`
	SigScheme                  string  `json:"sigScheme"`
	TrustScore                 float64 `json:"trustScore"`
	IsInCommittee              bool    `json:"isInCommittee"`
}

type GatewayVerifierInitInfoStructureGateway struct {
	GatewayPublicKeySignature string `json:"gatewayPublicKeySigniture"`
	GatewayPublicKeyKem       string `json:"gatewayPublicKeyKem"`
	GatewayIpAddress          string `json:"gatewayIpAddress"`
	GatewayPort               string `json:"gatewayPort"`
	KemScheme                 string `json:"kemScheme"`
	SigScheme                 string `json:"sigScheme"`
}
