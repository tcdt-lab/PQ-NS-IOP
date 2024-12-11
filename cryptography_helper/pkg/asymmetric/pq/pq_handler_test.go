package pq

import (
	b64 "encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPq_MLDSA(t *testing.T) {
	handler := Pq_handler{&MLDSA{}, &MLKEM{}}

	t.Run("mldsa", func(t *testing.T) {

		schemeName := "ML-DSA-65"
		msg := []byte("Hello World")
		sk, pk, err := handler.DSKeyGen(schemeName)
		assert.NoError(t, err, "Error while generating key")

		signedMsg, err := handler.Sign(sk, msg, schemeName)
		assert.NoError(t, err, "Error while signing message")

		result, err := handler.Verify(pk, msg, signedMsg, schemeName)
		assert.NoError(t, err, "Error while verifying message")
		assert.True(t, result, "Verification failed")
	})

	t.Run("mlkem", func(t *testing.T) {

		schemeName := "ML-KEM-768"

		sk, pk, err := handler.KEMKeyGen(schemeName)
		assert.NoError(t, err, "Error while generating key")

		//encapsulation
		cipherText, sharedKey, err := handler.KemGenerateSecretKey("", pk, "", schemeName)
		assert.NoError(t, err, "Error while generating secret key")
		assert.NotEmpty(t, cipherText, "Ciphertext is empty")
		assert.NotEmpty(t, sharedKey, "Shared key is empty")
		//decapsulation
		b64CipherText := b64.StdEncoding.EncodeToString(cipherText)
		_, sharedKeyDecap, err := handler.KemGenerateSecretKey(sk, "", b64CipherText, schemeName)
		assert.NoError(t, err, "Error while generating secret key")
		assert.NotEmpty(t, sharedKeyDecap, "Shared key is empty")
		assert.Equal(t, sharedKey, sharedKeyDecap, "Shared keys do not match")
	})

}
