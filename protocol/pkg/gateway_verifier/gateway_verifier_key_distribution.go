package gateway_verifier

// Gateway send its' information to verifier_verifier
type GatewayVerifierKeyDistributionRequest struct {
	RequestId                 int64  `json:"requestId"`
	GatewayCompanyName        string `json:"gatewayCompanyName"`
	GatewayPort               string `json:"gatewayPort"`
	GatewayIP                 string `json:"gatewayIP"`
	GatewaySignatureScheme    string `json:"gatewaySignitureScheme"`
	GatewayKemScheme          string `json:"gatewayKemScheme"`
	GatewayPublicKeySignature string `json:"gatewayPublicKeySigniture"`
	GatewayPublicKeyKem       string `json:"gatewayPublicKeyKem"`
}

type GatewayVerifierKeyDistributionResponse struct {
	OperationError string `json:"operationError"`
	CipherText     string `json:"cipherText"`
	PublicKeyKem   string `json:"publicKeyKem"`
	RequestId      int64  `json:"requestId"`
}
