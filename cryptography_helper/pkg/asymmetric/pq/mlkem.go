package pq

import (
	b64 "encoding/base64"
	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/schemes"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type MLKEM struct {
}

func (m *MLKEM) KeyGen(schemeName string) (kem.PublicKey, kem.PrivateKey, error) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	schemeToUse := schemes.ByName(schemeName)
	if schemeToUse == nil {
		zap.L().Error("Scheme unsupported", zap.String("schemeName", schemeName))
		return nil, nil, errors.New(" Scheme unsupported")
	}
	kSeed := make([]byte, schemeToUse.SeedSize())
	for i := 0; i < schemeToUse.SeedSize(); i++ {

		kSeed[i] = byte(rand.Intn(255))
	}
	pKey, sKey := schemeToUse.DeriveKeyPair(kSeed)
	if pKey == nil || sKey == nil {
		zap.L().Error("Failed to generate key pair")
		return nil, nil, errors.New("Failed to generate key pair")
	}
	return pKey, sKey, nil
}

func (m *MLKEM) EncapsulateDeterministically(pKey kem.PublicKey, schemeName string) ([]byte, []byte, error) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	schemeToUse := schemes.ByName(schemeName)
	if schemeToUse == nil {
		zap.L().Error("Scheme unsupported", zap.String("schemeName", schemeName))
		return nil, nil, errors.New(" Scheme unsupported")
	}
	encSeed := make([]byte, schemeToUse.EncapsulationSeedSize())
	for i := 0; i < schemeToUse.EncapsulationSeedSize(); i++ {
		encSeed[i] = byte(rand.Intn(255))
	}
	cipherText, sharedKey, err := schemeToUse.EncapsulateDeterministically(pKey, encSeed)
	if err != nil {
		zap.L().Error("Failed to encapsulate", zap.Error(err))
		return nil, nil, err
	}
	return cipherText, sharedKey, nil
}

func (m *MLKEM) Decapsulate(sKey kem.PrivateKey, cipherText []byte, schemeName string) ([]byte, error) {
	schemeToUse := schemes.ByName(schemeName)
	if schemeToUse == nil {
		zap.L().Error("Scheme unsupported", zap.String("schemeName", schemeName))
		return nil, errors.New(" Scheme unsupported")
	}
	sharedKey, err := schemeToUse.Decapsulate(sKey, cipherText)
	if err != nil {
		zap.L().Error("Failed to decapsulate", zap.Error(err))
		return nil, err
	}
	return sharedKey, nil
}

func (m *MLKEM) UnmarshalPublicKey(schemeName string, pubKey []byte) (kem.PublicKey, error) {
	return schemes.ByName(schemeName).UnmarshalBinaryPublicKey(pubKey)
}

func (m *MLKEM) UnmarshalPrivateKey(schemeName string, secKey []byte) (kem.PrivateKey, error) {
	return schemes.ByName(schemeName).UnmarshalBinaryPrivateKey(secKey)
}

func (m *MLKEM) MarshalPublicKey(pubKey kem.PublicKey) ([]byte, error) {
	return pubKey.MarshalBinary()
}

func (m *MLKEM) MarshalPrivateKey(secKey kem.PrivateKey) ([]byte, error) {
	return secKey.MarshalBinary()
}

func (m *MLKEM) ConvertPubKeyToBase64String(pubKey kem.PublicKey) string {
	pkk, _ := pubKey.MarshalBinary()
	return b64.StdEncoding.EncodeToString(pkk)
}

func (m *MLKEM) ConvertSecKeyToBase64String(secKey kem.PrivateKey) string {
	skk, _ := secKey.MarshalBinary()
	return b64.StdEncoding.EncodeToString(skk)
}

func (m *MLKEM) ConvertBase64StringToPubKey(pubKeyStr string, schemeName string) kem.PublicKey {
	pkk, _ := b64.StdEncoding.DecodeString(pubKeyStr)
	pubKey, _ := m.UnmarshalPublicKey(schemeName, pkk)

	return pubKey
}

func (m *MLKEM) ConvertBase64StringToSecKey(secKeyStr string, schemeName string) kem.PrivateKey {
	skk, _ := b64.StdEncoding.DecodeString(secKeyStr)
	secKey, _ := m.UnmarshalPrivateKey(schemeName, skk)

	return secKey
}
