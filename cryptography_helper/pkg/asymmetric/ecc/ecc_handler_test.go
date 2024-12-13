package ecc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEDSA(t *testing.T) {
	ecc := Ecc_handler{&Ecc_dh{}, &Ecc_dsa{}}

	t.Run("p256", func(t *testing.T) {
		schemeName := "p256"
		msg := []byte("Hello World")
		sk, pk, err := ecc.DSKeyGen(schemeName)
		assert.NoError(t, err, "Error while generating key")
		assert.NotEmpty(t, sk, "Secret key is empty")
		assert.NotEmpty(t, pk, "Public key is empty")

		signedMsg, err := ecc.Sign(sk, msg, schemeName)
		assert.NoError(t, err, "Error while signing message")
		assert.NotEmpty(t, signedMsg, "Signed message is empty")

		result, err := ecc.Verify(pk, msg, signedMsg, schemeName)
		assert.NoError(t, err, "Error while verifying message")
		assert.True(t, result, "Verification failed")
	})

	t.Run("p384", func(t *testing.T) {
		schemeName := "p384"
		msg := []byte("Hello World")
		sk, pk, err := ecc.DSKeyGen(schemeName)
		assert.NoError(t, err, "Error while generating key")
		assert.NotEmpty(t, sk, "Secret key is empty")
		assert.NotEmpty(t, pk, "Public key is empty")

		signedMsg, err := ecc.Sign(sk, msg, schemeName)
		assert.NoError(t, err, "Error while signing message")
		assert.NotEmpty(t, signedMsg, "Signed message is empty")

		result, err := ecc.Verify(pk, msg, signedMsg, schemeName)
		assert.NoError(t, err, "Error while verifying message")
		assert.True(t, result, "Verification failed")
	})

	t.Run("p521", func(t *testing.T) {
		schemeName := "p521"
		msg := []byte("Hello World")
		sk, pk, err := ecc.DSKeyGen(schemeName)
		assert.NoError(t, err, "Error while generating key")
		assert.NotEmpty(t, sk, "Secret key is empty")
		assert.NotEmpty(t, pk, "Public key is empty")

		signedMsg, err := ecc.Sign(sk, msg, schemeName)
		assert.NoError(t, err, "Error while signing message")
		assert.NotEmpty(t, signedMsg, "Signed message is empty")

		result, err := ecc.Verify(pk, msg, signedMsg, schemeName)
		assert.NoError(t, err, "Error while verifying message")
		assert.True(t, result, "Verification failed")
	})
}

func TestECDH(t *testing.T) {
	ecc := Ecc_handler{&Ecc_dh{}, &Ecc_dsa{}}

	t.Run("SharedSecretGenerating", func(t *testing.T) {
		sk, pk, err := ecc.KEMKeyGen("x25519")
		assert.NoError(t, err, "Error while generating key")
		assert.NotEmpty(t, sk, "Secret key is empty")
		assert.NotEmpty(t, pk, "Public key is empty")

		sharedSecret, sharedSecretHash, err := ecc.KemGenerateSecretKey(sk, pk, "", "x25519")

		assert.NoError(t, err, "Error while generating shared secret")
		assert.NotEmpty(t, sharedSecret, "Shared secret is empty")
		assert.NotEmpty(t, sharedSecretHash, "Shared secret hash is empty")
		t.Log(sharedSecretHash)
	})
}
