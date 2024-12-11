package pq

import b64 "encoding/base64"

type Pq_handler struct {
	mldsa *MLDSA
	mlkem *MLKEM
}

func (p *Pq_handler) DSKeyGen(schemeName string) (string, string, error) {
	pubKey, secKey, err := p.mldsa.KeyGen(schemeName)
	if err != nil {
		return "", "", err
	}
	pkStr := p.mldsa.ConvertPubKeyToBase64String(pubKey)
	skStr := p.mldsa.ConvertSecKeyToBase64String(secKey)

	return skStr, pkStr, nil
}

func (p *Pq_handler) Sign(sKey string, message []byte, schemeName string) ([]byte, error) {
	secretKey := p.mldsa.ConvertBase64StringToSecKey(sKey, schemeName)

	result := p.mldsa.Sign(message, secretKey, schemeName)
	return result, nil
}

func (p *Pq_handler) Verify(pkey string, message []byte, signature []byte, schemeName string) (bool, error) {
	publicKey := p.mldsa.ConvertBase64StringToPubKey(pkey, schemeName)

	result, err := p.mldsa.Verify(publicKey, message, signature, schemeName)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (p *Pq_handler) KEMKeyGen(schemeName string) (string, string, error) {
	pubKey, secKey, err := p.mlkem.KeyGen(schemeName)

	if err != nil {
		return "", "", err
	}
	pkStr := p.mlkem.ConvertPubKeyToBase64String(pubKey)
	skStr := p.mlkem.ConvertSecKeyToBase64String(secKey)
	return skStr, pkStr, nil
}

func (p *Pq_handler) KemGenerateSecretKey(secKey string, pubKey string, ciphertext string, schemeName string) ([]byte, []byte, error) {

	if secKey == "" && ciphertext == "" {
		//we go for encapsulation
		publicKey := p.mlkem.ConvertBase64StringToPubKey(pubKey, schemeName)
		cipherText, sharedKey, err := p.mlkem.EncapsulateDeterministically(publicKey, schemeName)
		if err != nil {
			return nil, nil, err
		}
		return cipherText, sharedKey, nil
	} else {
		//we go for decapsulation
		secretKey := p.mlkem.ConvertBase64StringToSecKey(secKey, schemeName)
		cipherTextBytes, err := b64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return nil, nil, err
		}
		sharedKey, err := p.mlkem.Decapsulate(secretKey, cipherTextBytes, schemeName)
		if err != nil {
			return nil, nil, err
		}
		return nil, sharedKey, nil
	}

}
