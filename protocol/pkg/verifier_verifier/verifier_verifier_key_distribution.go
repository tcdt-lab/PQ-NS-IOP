package verifier_verifier

type VVerifierKeyDistributionRequest struct {
	RequestId int64 `json:"requestId"`

	VerifierPort               string `json:"gatewayPort"`
	VerifierIP                 string `json:"gatewayIP"`
	VerifierSignatureScheme    string `json:"gatewaySignitureScheme"`
	VerifierKemScheme          string `json:"gatewayKemScheme"`
	VerifierPublicKeySignature string `json:"gatewayPublicKeySigniture"`
	VerifierPublicKeyKem       string `json:"gatewayPublicKeyKem"`
}

type VVerifierKeyDistributionResponse struct {
	OperationError string `json:"operationError"`
	CipherText     string `json:"cipherText"`
	PublicKeyKem   string `json:"publicKeyKem"`
	RequestId      int64  `json:"requestId"`
}
