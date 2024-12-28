package pq

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyGen(t *testing.T) {
	fmt.Println("**********************************************")
	fmt.Println("#Testing ML-KEM KeyGen")

	var mlkem MLKEM
	pk, sk, err := mlkem.KeyGen("ML-KEM-768")
	if err != nil {
		t.Errorf("Error in KeyGen: %v", err)
	}
	if pk == nil || sk == nil {
		t.Errorf("Failed to generate key pair")
	}

	fmt.Println("Successfully generated keys")
	// marshaling the kesys and printing the first 32 bytes of keys
	pkk, _ := pk.MarshalBinary()
	skk, _ := sk.MarshalBinary()
	fmt.Println(hex.EncodeToString(pkk[:32]))
	fmt.Println(hex.EncodeToString(skk[:32]))

}

func TestEncapsulateDeterministically(t *testing.T) {
	fmt.Println("**********************************************")
	fmt.Println("#Testing ML-KEM EncapsulateDeterministically")

	var mlkem MLKEM
	pk, _, _ := mlkem.KeyGen("ML-KEM-768")
	t.Run("Test EncapsulateDeterministically", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			cipherText, sharedKey, err := mlkem.EncapsulateDeterministically(pk, "ML-KEM-768")
			if err != nil {
				t.Errorf("Error in EncapsulateDeterministically: %v", err)
			}
			if cipherText == nil || sharedKey == nil {
				t.Errorf("Failed to encapsulate")
			}
			fmt.Println("Successfully encapsulated")
			fmt.Println("cipherText: ", hex.EncodeToString(cipherText))
			fmt.Println("sharedKey: ", hex.EncodeToString(sharedKey[:32]))
		}
	})
	cipherText, sharedKey, err := mlkem.EncapsulateDeterministically(pk, "ML-KEM-768")
	if err != nil {
		t.Errorf("Error in EncapsulateDeterministically: %v", err)
	}
	if cipherText == nil || sharedKey == nil {
		t.Errorf("Failed to encapsulate")
	}

	fmt.Println("Successfully encapsulated")
	fmt.Println("cipherText: ", hex.EncodeToString(cipherText))
	fmt.Println("sharedKey: ", hex.EncodeToString(sharedKey[:32]))
}

func TestDecapsulate(t *testing.T) {
	fmt.Println("**********************************************")
	fmt.Println("#Testing ML-KEM Decapsulate")

	var mlkem MLKEM
	_, sk, _ := mlkem.KeyGen("ML-KEM-768")
	cipherText, encSharedKey, _ := mlkem.EncapsulateDeterministically(sk.Public(), "ML-KEM-768")
	decSharedKey, err := mlkem.Decapsulate(sk, cipherText, "ML-KEM-768")
	if err != nil {
		t.Errorf("Error in Decapsulate: %v", err)
	}
	if decSharedKey == nil {
		t.Errorf("Failed to decapsulate")
	}

	if hex.EncodeToString(decSharedKey) != hex.EncodeToString(encSharedKey) {
		t.Errorf("Decapsulated shared key does not match the encapsulation shared key")
	}

	fmt.Println("Successfully decapsulated")
	fmt.Println("decSharedKey: ", hex.EncodeToString(decSharedKey[:32]))
}

func TestMlKEm_Str(t *testing.T) {
	fmt.Println("**********************************************")
	fmt.Println("#Testing ML-KEM ConvertPubKeyToBase64String")

	var mlkem MLKEM
	pk, sk, err := mlkem.KeyGen("ML-KEM-768")
	assert.NoError(t, err, "Error in KeyGen")
	ciphertext, SharedKey, err := mlkem.EncapsulateDeterministically(pk, "ML-KEM-768")
	assert.NoError(t, err, "Error in EncapsulateDeterministically")
	sharedKey2, err := mlkem.Decapsulate(sk, ciphertext, "ML-KEM-768")
	assert.NoError(t, err, "Error in Decapsulate")
	assert.Equal(t, SharedKey, sharedKey2, "Decapsulated shared key does not match the encapsulation shared key")

	pkStr := mlkem.ConvertPubKeyToBase64String(pk)
	if pkStr == "" {
		t.Errorf("Failed to convert public key to base64 string")
	}
	skStr := mlkem.ConvertSecKeyToBase64String(sk)

	newSk := mlkem.ConvertBase64StringToSecKey(skStr, "ML-KEM-768")
	assert.NotNil(t, newSk, "Error in ConvertBase64StringToSecKey")
	newPk := mlkem.ConvertBase64StringToPubKey(pkStr, "ML-KEM-768")
	assert.NotNil(t, newPk, "Error in ConvertBase64StringToPubKey")

	assert.Equal(t, newSk, sk, "Private keys are not equal")
	assert.Equal(t, newPk, pk, "Public keys are not equal")

	newSharedkey2, err := mlkem.Decapsulate(newSk, ciphertext, "ML-KEM-768")
	assert.NoError(t, err, "Error in Decapsulate")

	assert.Equal(t, SharedKey, newSharedkey2, "Decapsulated shared key does not match the encapsulation shared key")

}
