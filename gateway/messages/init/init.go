package init

// These structs will be used to parse the JSON messages (params)
type Init struct {
	// Init operation
	sourcePublicKey string "json:sourcePublicKey"
	sourceHostName  string "json:sourceHostName"
	sourceIP        string "json:sourceIP"
	kem_pub_key     string "json:kem_pub_key"
	kem_scheme      string "json kem_scheme"
}

type InitResponse struct {
	kem_cipher_text string         "json:kem_cipher_text"
	verifiersInfo   []VerifierInfo "json:verifiers"
}

type VerifierInfo struct {
	ip        string "json:ip"
	port      int    "json:port"
	publicKey string "json:publicKey"
}
type InitResponseError struct {
	errorMessage string "json:errorMessage"
}
