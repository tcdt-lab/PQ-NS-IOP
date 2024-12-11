package ecc

const ()

type Ecc_handler struct {
	eccdh  *ecc_dh
	eccdsa *ecc_dsa
}

func (e *Ecc_handler) DSKeyGen(schemeName string) (string, string, error) {
	sk, pk, err := e.eccdsa.DSKeyGen(schemeName)
	if err != nil {
		return "", "", err
	}
	pkStr := e.eccdsa.ConvertPublicKeyToBase64String(pk)

	skStr := e.eccdsa.ConvertSecretKeyToBase64String(sk)
	return skStr, pkStr, nil
}

func (e *Ecc_handler) Sign(sKey string, message []byte, schemeName string) ([]byte, error) {
	secretKey, err := e.eccdsa.ConvertBase64StringToSecretKey(sKey)
	if err != nil {
		return nil, err
	}
	result, err := e.eccdsa.Sign(message, secretKey)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Ecc_handler) Verify(pkey string, message []byte, signature []byte, schemeName string) (bool, error) {
	publicKey, err := e.eccdsa.ConvertBase64StringToPublicKey(pkey)
	if err != nil {
		return false, err
	}
	result, err := e.eccdsa.Verify(publicKey, message, signature)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (e *Ecc_handler) KEMKeyGen(schemeName string) (string, string, error) {
	sk, pk, err := e.eccdh.DHKeyGen(schemeName)
	if err != nil {
		return "", "", err
	}
	pkStr, err := e.eccdh.ConvertPublicKeyToBase64String(pk)
	if err != nil {
		return "", "", err
	}
	skStr, err := e.eccdh.ConvertSecretKeyToBase64String(sk)
	if err != nil {
		return "", "", err
	}

	return skStr, pkStr, nil
}

func (e *Ecc_handler) KemGenerateSecretKey(secKey string, pubKey string, ciphertext string, schemeName string) ([]byte, []byte, error) {
	secretKey, err := e.eccdh.ConvertBase64StringToSecretKey(secKey)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := e.eccdh.ConvertBase64StringToPublicKey(pubKey)
	if err != nil {
		return nil, nil, err
	}

	sharedSecret, sharedSecretShaw256, err := e.eccdh.GenerateSharedSecret(secretKey, publicKey)
	return sharedSecret, sharedSecretShaw256, err
}
