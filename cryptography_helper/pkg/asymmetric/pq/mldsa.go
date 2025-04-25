package pq

import (
	b64 "encoding/base64"
	"encoding/hex"
	"github.com/cloudflare/circl/sign"
	"github.com/cloudflare/circl/sign/schemes"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type MLDSA struct {
}

func (m *MLDSA) UnmarshalPublicKey(schemeName string, pubKey []byte) (sign.PublicKey, error) {
	return schemes.ByName(schemeName).UnmarshalBinaryPublicKey(pubKey)
}

func (m *MLDSA) UnmarshalSecretKey(schemeName string, secKey []byte) (sign.PrivateKey, error) {
	return schemes.ByName(schemeName).UnmarshalBinaryPrivateKey(secKey)
}

func (m *MLDSA) KeyGen(schemeName string) (sign.PublicKey, sign.PrivateKey, error) {
	scheme := schemes.ByName(schemeName)
	pKey, sKey, err := scheme.GenerateKey()
	if err != nil {
		zap.L().Error("Error while generating ML-DSA keys", zap.String("scheme", schemeName))
		return nil, nil, err
	}
	pkb, _ := pKey.MarshalBinary()
	zap.L().Info("Successful Key Generation", zap.String("Pub Key (first 32 bytes)", hex.EncodeToString(pkb[:32])), zap.String("Scheme", schemeName))
	return pKey, sKey, nil
}

func (m *MLDSA) Sign(encryptedMsg []byte, sKey sign.PrivateKey, schemeName string) []byte {
	opts := &sign.SignatureOpts{}
	scheme := schemes.ByName(schemeName)
	signature := scheme.Sign(sKey, encryptedMsg, opts)
	zap.L().Info("Successful Signing", zap.String("Signature", hex.EncodeToString(signature[:32])))
	return signature
}

func (m *MLDSA) Verify(pk sign.PublicKey, message []byte, signature []byte, schemeName string) (bool, error) {
	//pkk, _ := pk.MarshalBinary()
	opts := &sign.SignatureOpts{}
	scheme := schemes.ByName(schemeName)
	if !scheme.Verify(pk, message, signature, opts) {
		zap.L().Error("Failed signature verification", zap.String("Signature", hex.EncodeToString(signature[:32])), zap.String("Message", hex.EncodeToString(message)), zap.String("Signature", hex.EncodeToString(signature[:32])))
		return false, errors.New("Failed Signature")
	} else {
		zap.L().Info("Successful verification", zap.String("Signature", hex.EncodeToString(signature[:32])), zap.String("Signature", hex.EncodeToString(signature[:32])))
		return true, nil
	}

}

func (m *MLDSA) ConvertPubKeyToBase64String(pubKey sign.PublicKey) string {
	pkk, _ := pubKey.MarshalBinary()
	return b64.StdEncoding.EncodeToString(pkk)
}

func (m *MLDSA) ConvertBase64StringToPubKey(pubKeyStr string, schemeName string) sign.PublicKey {
	pkk, _ := b64.StdEncoding.DecodeString(pubKeyStr)
	pubKey, _ := m.UnmarshalPublicKey(schemeName, pkk)

	return pubKey
}

func (m *MLDSA) ConvertSecKeyToBase64String(secKey sign.PrivateKey) string {
	skk, _ := secKey.MarshalBinary()
	return b64.StdEncoding.EncodeToString(skk)
}

func (m *MLDSA) ConvertBase64StringToSecKey(secKeyStr string, schemeName string) sign.PrivateKey {
	skk, _ := b64.StdEncoding.DecodeString(secKeyStr)
	secKey, _ := m.UnmarshalSecretKey(schemeName, skk)

	return secKey
}
