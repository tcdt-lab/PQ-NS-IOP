package ecc

const ()

type Ecc_handler struct {
	Eccdh  *Ecc_dh
	Eccdsa *Ecc_dsa
}

func (e *Ecc_handler) DSKeyGen(schemeName string) (string, string, error) {
	sk, pk, err := e.Eccdsa.DSKeyGen(schemeName)
	if err != nil {
		return "", "", err
	}
	pkStr := e.Eccdsa.ConvertPublicKeyToBase64String(pk)

	skStr := e.Eccdsa.ConvertSecretKeyToBase64String(sk)
	return skStr, pkStr, nil
}

func (e *Ecc_handler) Sign(sKey string, message []byte, schemeName string) ([]byte, error) {
	secretKey, err := e.Eccdsa.ConvertBase64StringToSecretKey(sKey)
	if err != nil {
		return nil, err
	}
	result, err := e.Eccdsa.Sign(message, secretKey)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Ecc_handler) Verify(pkey string, message []byte, signature []byte, schemeName string) (bool, error) {
	publicKey, err := e.Eccdsa.ConvertBase64StringToPublicKey(pkey)
	if err != nil {
		return false, err
	}
	result, err := e.Eccdsa.Verify(publicKey, message, signature)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (e *Ecc_handler) KEMKeyGen(schemeName string) (string, string, error) {
	sk, pk, err := e.Eccdh.DHKeyGen(schemeName)
	if err != nil {
		return "", "", err
	}
	pkStr, err := e.Eccdh.ConvertPublicKeyToBase64String(pk)
	if err != nil {
		return "", "", err
	}
	skStr, err := e.Eccdh.ConvertSecretKeyToBase64String(sk)
	if err != nil {
		return "", "", err
	}

	return skStr, pkStr, nil
}

func (e *Ecc_handler) KemGenerateSecretKey(secKey string, pubKey string, ciphertext string, schemeName string) ([]byte, []byte, error) {
	secretKey, err := e.Eccdh.ConvertBase64StringToSecretKey(secKey)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := e.Eccdh.ConvertBase64StringToPublicKey(pubKey)
	if err != nil {
		return nil, nil, err
	}

	sharedSecret, sharedSecretShaw256, err := e.Eccdh.GenerateSharedSecret(secretKey, publicKey)
	return sharedSecret, sharedSecretShaw256, err
}
