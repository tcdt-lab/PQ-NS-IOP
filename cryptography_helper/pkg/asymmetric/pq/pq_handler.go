package pq

import b64 "encoding/base64"

type Pq_handler struct {
	Mldsa *MLDSA
	Mlkem *MLKEM
}

func (p *Pq_handler) DSKeyGen(schemeName string) (string, string, error) {
	pubKey, secKey, err := p.Mldsa.KeyGen(schemeName)
	if err != nil {
		return "", "", err
	}
	pkStr := p.Mldsa.ConvertPubKeyToBase64String(pubKey)
	skStr := p.Mldsa.ConvertSecKeyToBase64String(secKey)

	return skStr, pkStr, nil
}

func (p *Pq_handler) Sign(sKey string, message []byte, schemeName string) ([]byte, error) {
	secretKey := p.Mldsa.ConvertBase64StringToSecKey(sKey, schemeName)

	result := p.Mldsa.Sign(message, secretKey, schemeName)
	return result, nil
}

func (p *Pq_handler) Verify(pkey string, message []byte, signature []byte, schemeName string) (bool, error) {
	publicKey := p.Mldsa.ConvertBase64StringToPubKey(pkey, schemeName)

	result, err := p.Mldsa.Verify(publicKey, message, signature, schemeName)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (p *Pq_handler) KEMKeyGen(schemeName string) (string, string, error) {
	pubKey, secKey, err := p.Mlkem.KeyGen(schemeName)

	if err != nil {
		return "", "", err
	}
	pkStr := p.Mlkem.ConvertPubKeyToBase64String(pubKey)
	skStr := p.Mlkem.ConvertSecKeyToBase64String(secKey)
	return skStr, pkStr, nil
}

func (p *Pq_handler) KemGenerateSecretKey(secKey string, pubKey string, ciphertext string, schemeName string) ([]byte, []byte, error) {

	if ciphertext == "" {
		//we go for encapsulation
		publicKey := p.Mlkem.ConvertBase64StringToPubKey(pubKey, schemeName)
		cipherText, sharedKey, err := p.Mlkem.EncapsulateDeterministically(publicKey, schemeName)
		if err != nil {
			return nil, nil, err
		}
		return cipherText, sharedKey, nil
	} else {
		//we go for decapsulation
		secretKey := p.Mlkem.ConvertBase64StringToSecKey(secKey, schemeName)
		cipherTextBytes, err := b64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return nil, nil, err
		}
		sharedKey, err := p.Mlkem.Decapsulate(secretKey, cipherTextBytes, schemeName)
		if err != nil {
			return nil, nil, err
		}
		return nil, sharedKey, nil
	}

}
