package gateway_verifier

type GatewayVerifierInitInfoOperation struct {
	RequestId int64 `json:"requestId"`
}
type GatewayVerifierInitInfoResponseOperation struct {
	OperationError      string                             `json:"operationError"`
	CurrentVerifierInfo GatewayVerifierInitInfoStructure   `json:"currentVerifierInfo"`
	VerifiersList       []GatewayVerifierInitInfoStructure `json:"verifierList"`
}
type GatewayVerifierInitInfoStructure struct {
	VerifierPublicKeySignature string  `json:"verifierPublicKeySigniture"`
	VerifierPublicKeyKem       string  `json:"verifierPublicKeyKem"`
	VerifierIpAddress          string  `json:"verifierIpAddress"`
	VerifierPort               string  `json:"verifierPort"`
	SigScheme                  string  `json:"sigScheme"`
	TrustScore                 float64 `json:"trustScore"`
	IsInCommittee              bool    `json:"isInCommittee"`
}
