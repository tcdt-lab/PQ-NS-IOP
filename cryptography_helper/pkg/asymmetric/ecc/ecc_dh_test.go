package ecc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEcc_dh_DHKeyGen(t *testing.T) {
	var e ecc_dh
	t.Run("Test DHKeyGen_x", func(t *testing.T) {
		secKEy, pubKey, err := e.DHKeyGen("x25519")
		assert.NoError(t, err, "Error in DHKeyGen")

		secKeyBytes, err := e.MarshalSecretKey(secKEy)
		assert.NoError(t, err, "Error in MarshalSecretKey")
		UnmarshalledKey, err := e.UnmarshalSecretKey(secKeyBytes)
		assert.NoError(t, err, "Error in UnmarshalSecretKey")
		//assert.Equal(t, secKEy, UnmarshalledKey, "Error in UnmarshalSecretKey")

		pubKeyBytes, err := e.MarshalPublicKey(pubKey)
		assert.NoError(t, err, "Error in MarshalPublicKey")
		UnmarshalledPubKey, err := e.UnmarshalPublicKey(pubKeyBytes)
		assert.NoError(t, err, "Error in UnmarshalPublicKey")
		//assert.Equal(t, pubKey, UnmarshalledPubKey, "Error in UnmarshalPublicKey")

		secret1, err := secKEy.ECDH(pubKey)
		assert.NoError(t, err, "Error in ECDH")

		secret2, err := UnmarshalledKey.ECDH(UnmarshalledPubKey)

		assert.NoError(t, err, "Error in ECDH")
		assert.Equal(t, secret1, secret2, "Error in ECDH")

	})

}

func TestEcc_dhString(t *testing.T) {
	var e ecc_dh
	t.Run("Test ConvertSecretKeyToBase64String", func(t *testing.T) {
		secKey, _, err := e.DHKeyGen("x25519")
		assert.NoError(t, err, "Error in DHKeyGen")
		_, err = e.ConvertSecretKeyToBase64String(secKey)
		assert.NoError(t, err, "Error in ConvertSecretKeyToBase64String")
	})
	t.Run("Test ConvertPublicKeyToBase64String", func(t *testing.T) {
		_, pubKey, err := e.DHKeyGen("x25519")
		assert.NoError(t, err, "Error in DHKeyGen")
		_, err = e.ConvertPublicKeyToBase64String(pubKey)
		assert.NoError(t, err, "Error in ConvertPublicKeyToBase64String")
	})
	t.Run("Test ConvertBase64StringToPublicKey", func(t *testing.T) {
		_, pubKey, err := e.DHKeyGen("x25519")
		assert.NoError(t, err, "Error in DHKeyGen")
		pubKeyStr, err := e.ConvertPublicKeyToBase64String(pubKey)
		assert.NoError(t, err, "Error in ConvertPublicKeyToBase64String")
		_, err = e.ConvertBase64StringToPublicKey(pubKeyStr)
		assert.NoError(t, err, "Error in ConvertBase64StringToPublicKey")
	})
	t.Run("Test ConvertBase64StringToSecretKey", func(t *testing.T) {
		secKey, _, err := e.DHKeyGen("x25519")
		assert.NoError(t, err, "Error in DHKeyGen")
		secKeyStr, err := e.ConvertSecretKeyToBase64String(secKey)
		assert.NoError(t, err, "Error in ConvertSecretKeyToBase64String")
		_, err = e.ConvertBase64StringToSecretKey(secKeyStr)
		assert.NoError(t, err, "Error in ConvertBase64StringToSecretKey")
	})
	t.Run("Test Full Operation", func(t *testing.T) {

		secKey, pubKey, err := e.DHKeyGen("x25519")
		assert.NoError(t, err, "Error in DHKeyGen")
		pubKeyStr, err := e.ConvertPublicKeyToBase64String(pubKey)
		assert.NoError(t, err, "Error in ConvertPublicKeyToBase64String")
		secKeyStr, err := e.ConvertSecretKeyToBase64String(secKey)
		assert.NoError(t, err, "Error in ConvertSecretKeyToBase64String")
		pubKey2, err := e.ConvertBase64StringToPublicKey(pubKeyStr)
		assert.NoError(t, err, "Error in ConvertBase64StringToPublicKey")
		secKey2, err := e.ConvertBase64StringToSecretKey(secKeyStr)
		assert.NoError(t, err, "Error in ConvertBase64StringToSecretKey")

		secret1, err := secKey2.ECDH(pubKey)
		assert.NoError(t, err, "Error in ECDH")
		secret2, err := secKey.ECDH(pubKey2)
		assert.NoError(t, err, "Error in ECDH")
		assert.Equal(t, secret1, secret2, "Error in ECDH secrets are not equal")

		secret3, err := secKey.ECDH(pubKey)
		assert.NoError(t, err, "Error in ECDH")
		secret4, err := secKey2.ECDH(pubKey2)
		assert.NoError(t, err, "Error in ECDH")
		assert.Equal(t, secret3, secret4, "Error in ECDH secrets are not equal")
	})
}
