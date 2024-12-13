package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	b64 "encoding/base64"
	"errors"
)

type Ecc_dsa struct {
}

func (Ecc_dsa) DSKeyGen(schemeName string) (privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey, err error) {
	if schemeName == "p256" {
		secKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}

		return secKey, &secKey.PublicKey, nil
	} else if schemeName == "p384" {
		secKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}

		return secKey, &secKey.PublicKey, nil
	} else if schemeName == "p521" {
		secKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}

		return secKey, &secKey.PublicKey, nil
	} else {
		return nil, nil, errors.New("Unsupported scheme")
	}

}

func (e *Ecc_dsa) UnmarshalPublicKey(pubKey []byte) (*ecdsa.PublicKey, error) {

	publicKey, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	return publicKey.(*ecdsa.PublicKey), nil
}
func (e *Ecc_dsa) UnmarshalSecretKey(secKey []byte) (*ecdsa.PrivateKey, error) {
	privateKey, err := x509.ParseECPrivateKey(secKey)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func (e *Ecc_dsa) MarshalSecretKey(secKey *ecdsa.PrivateKey) ([]byte, error) {
	return x509.MarshalECPrivateKey(secKey)
}
func (e *Ecc_dsa) MarshalPublicKey(pubKey *ecdsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pubKey)
}

func (e *Ecc_dsa) Sign(messages []byte, secKey *ecdsa.PrivateKey) ([]byte, error) {
	signedBytes, err := ecdsa.SignASN1(rand.Reader, secKey, messages)
	if err != nil {
		return nil, err
	}
	return signedBytes, nil
}

func (e *Ecc_dsa) Verify(pubKey *ecdsa.PublicKey, messages []byte, signature []byte) (bool, error) {
	return ecdsa.VerifyASN1(pubKey, messages, signature), nil
}

func (e *Ecc_dsa) ConvertPublicKeyToBase64String(pubKey *ecdsa.PublicKey) string {
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(pubKey)
	return b64.StdEncoding.EncodeToString(pubKeyBytes)
}

func (e *Ecc_dsa) ConvertBase64StringToPublicKey(pubKeyStr string) (*ecdsa.PublicKey, error) {
	pubKeyBytes, err := b64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}
	pubKey, err := e.UnmarshalPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func (e *Ecc_dsa) ConvertSecretKeyToBase64String(secKey *ecdsa.PrivateKey) string {
	secKeyBytes, _ := x509.MarshalECPrivateKey(secKey)
	return b64.StdEncoding.EncodeToString(secKeyBytes)
}

func (e *Ecc_dsa) ConvertBase64StringToSecretKey(secKeyStr string) (*ecdsa.PrivateKey, error) {
	secKeyBytes, err := b64.StdEncoding.DecodeString(secKeyStr)
	if err != nil {
		return nil, err
	}
	secKey, err := e.UnmarshalSecretKey(secKeyBytes)
	if err != nil {
		return nil, err
	}
	return secKey, nil
}
