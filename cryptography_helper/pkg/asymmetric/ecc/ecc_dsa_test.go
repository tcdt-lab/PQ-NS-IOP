package ecc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MarshalingUnmarshaling(t *testing.T) {
	eccdsa := ecc_dsa{}
	var marshalledPubKey []byte
	var marshalledSecKey []byte
	t.Run("Test MarshalUnMarshalPublicKey", func(t *testing.T) {
		// Act
		_, pubKey, err := eccdsa.DSKeyGen("p256")
		if err != nil {
			t.Errorf("Error in KeyGen: %v", err)
		}
		marshalledPubKey, err = eccdsa.MarshalPublicKey(pubKey)
		assert.NoError(t, err, "Error in marshalling public key")
		unmarshalledPubKey, err := eccdsa.UnmarshalPublicKey(marshalledPubKey)
		assert.NoError(t, err, "Error in unmarshalling public key")
		assert.Equal(t, pubKey, unmarshalledPubKey, "Public keys do not match")
		// Assert
	})
	t.Run("Test MarshalUnMarshalSecretKey", func(t *testing.T) {
		// Act
		secKey, _, err := eccdsa.DSKeyGen("p256")
		if err != nil {
			t.Errorf("Error in KeyGen: %v", err)
		}
		marshalledSecKey, err = eccdsa.MarshalSecretKey(secKey)
		assert.NoError(t, err, "Error in marshalling secret key")
		unmarshalledSecKey, err := eccdsa.UnmarshalSecretKey(marshalledSecKey)
		assert.NoError(t, err, "Error in unmarshalling secret key")
		assert.Equal(t, secKey, unmarshalledSecKey, "Secret keys do not match")
		// Assert
	})
}

func Test_SignVerify(t *testing.T) {
	eccdsa := ecc_dsa{}
	t.Run("Test SignVerify", func(t *testing.T) {
		// Act
		secKey, pubKey, err := eccdsa.DSKeyGen("p256")
		if err != nil {
			t.Errorf("Error in KeyGen: %v", err)
		}

		message := []byte("Hello World")
		signature, err := eccdsa.Sign(message, secKey)
		assert.NoError(t, err, "Error in signing")
		res, err := eccdsa.Verify(pubKey, message, signature)
		assert.NoError(t, err, "Error in verifying")
		assert.True(t, res, "Signature verification failed")
		// Assert
	})
}

func TestStringKey(t *testing.T) {
	eccdsa := ecc_dsa{}
	t.Run("Test StringKey", func(t *testing.T) {
		// Act
		secKey, pubKey, err := eccdsa.DSKeyGen("p256")
		assert.NoError(t, err, "Error in KeyGen")
		signiture1, err := eccdsa.Sign([]byte("Hello World"), secKey)
		assert.NoError(t, err, "Error in Sign")
		pubKeyStr := eccdsa.ConvertPublicKeyToBase64String(pubKey)
		privKeyStr := eccdsa.ConvertSecretKeyToBase64String(secKey)

		pubKey2, err := eccdsa.ConvertBase64StringToPublicKey(pubKeyStr)
		assert.NoError(t, err, "Error in ConvertBase64StringToPublicKey")
		privKey2, err := eccdsa.ConvertBase64StringToSecretKey(privKeyStr)
		assert.NoError(t, err, "Error in ConvertBase64StringToSecretKey")

		signiture2, err := eccdsa.Sign([]byte("Hello World"), privKey2)
		assert.NoError(t, err, "Error in Sign")

		res, err := eccdsa.Verify(pubKey2, []byte("Hello World"), signiture1)
		assert.NoError(t, err, "Error in Verify")
		assert.True(t, res, "Signature verification failed")

		res, err = eccdsa.Verify(pubKey, []byte("Hello World"), signiture2)
		assert.NoError(t, err, "Error in Verify")
		assert.True(t, res, "Signature verification failed")
	})
}
