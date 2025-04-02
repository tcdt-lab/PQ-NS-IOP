package verifier_verifier

type VVSyncInfoOperationResponse struct {
	RequestId int64 `json:"requestId"`
	Result    bool  `json:"result"`
}
type VVSyncInfoOperationRequest struct {
	RequestId           int64                         `json:"requestId"`
	OperationError      string                        `json:"operationError"`
	CurrentVerifierInfo VVSyncInfoStructureVerifier   `json:"currentVerifierInfo"`
	VerifiersList       []VVSyncInfoStructureVerifier `json:"verifierList"`
	GatewaysList        []VVSyncInfoStructureGateway  `json:"gatewayList"`
}
type VVSyncInfoStructureVerifier struct {
	VerifierPublicKeySignature string  `json:"verifierPublicKeySigniture"`
	VerifierPublicKeyKem       string  `json:"verifierPublicKeyKem"`
	VerifierIpAddress          string  `json:"verifierIpAddress"`
	VerifierPort               string  `json:"verifierPort"`
	SigScheme                  string  `json:"sigScheme"`
	TrustScore                 float64 `json:"trustScore"`
	IsInCommittee              bool    `json:"isInCommittee"`
	SymmetricKey               string  `json:"symmetricKey"`
}

type VVSyncInfoStructureGateway struct {
	GatewayPublicKeySignature string `json:"gatewayPublicKeySigniture"`
	GatewayPublicKeyKem       string `json:"gatewayPublicKeyKem"`
	GatewayIpAddress          string `json:"gatewayIpAddress"`
	GatewayPort               string `json:"gatewayPort"`
	KemScheme                 string `json:"kemScheme"`
	SigScheme                 string `json:"sigScheme"`
	SymmetricKey              string `json:"symmetricKey"`
}
