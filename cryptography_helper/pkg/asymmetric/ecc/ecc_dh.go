package ecc

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	b64 "encoding/base64"
	"errors"
)

type Ecc_dh struct {
}

func (e *Ecc_dh) DHKeyGen(schemeName string) (*ecdh.PrivateKey, *ecdh.PublicKey, error) {
	if schemeName == "x25519" {
		clientCurve := ecdh.X25519()
		clientSecKey, err := clientCurve.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		pubKey := clientSecKey.PublicKey()
		return clientSecKey, pubKey, nil
	} else {
		return nil, nil, errors.New("Unsupported scheme")
	}
}

func (e *Ecc_dh) UnmarshalPublicKey(pubKey []byte) (*ecdh.PublicKey, error) {
	publicKey, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	return publicKey.(*ecdh.PublicKey), nil
}
func (e *Ecc_dh) UnmarshalSecretKey(secKey []byte) (*ecdh.PrivateKey, error) {
	privateKey, err := x509.ParsePKCS8PrivateKey(secKey)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return privateKey.(*ecdh.PrivateKey), nil
}

func (e *Ecc_dh) MarshalSecretKey(secKey *ecdh.PrivateKey) ([]byte, error) {
	return x509.MarshalPKCS8PrivateKey(secKey)
}

func (e *Ecc_dh) MarshalPublicKey(pubKey *ecdh.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pubKey)
}

func (e *Ecc_dh) ConvertSecretKeyToBase64String(secKey *ecdh.PrivateKey) (string, error) {
	secKeyBytes, err := e.MarshalSecretKey(secKey)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(secKeyBytes), nil
}

func (e *Ecc_dh) ConvertPublicKeyToBase64String(pubKey *ecdh.PublicKey) (string, error) {
	pubKeyBytes, err := e.MarshalPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(pubKeyBytes), nil
}

func (e *Ecc_dh) ConvertBase64StringToPublicKey(pubKeyStr string) (*ecdh.PublicKey, error) {
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

func (e *Ecc_dh) ConvertBase64StringToSecretKey(secKeyStr string) (*ecdh.PrivateKey, error) {
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

func (e *Ecc_dh) GenerateSharedSecret(secKey *ecdh.PrivateKey, pubKey *ecdh.PublicKey) ([]byte, []byte, error) {
	secret, err := secKey.ECDH(pubKey)
	if err != nil {
		return nil, nil, err
	}
	h := sha256.New()
	h.Write(secret)
	hash := h.Sum(nil)
	return secret, hash, nil
}
