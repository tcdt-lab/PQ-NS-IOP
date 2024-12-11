package pq

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestMLDSA_KeyGen(t *testing.T) {
	fmt.Println("**********************************************")
	fmt.Println("#Testing ML-DSA KeyGen")
	var mldsa MLDSA
	pubkey, privatekey, err := mldsa.KeyGen("ML-DSA-65")
	pkk, _ := pubkey.MarshalBinary()
	skk, _ := privatekey.MarshalBinary()
	if err != nil {
		t.Errorf("Error in KeyGen: %v", err)
	}
	fmt.Println("Successfully generated keys")
	fmt.Println(hex.EncodeToString(pkk[:32]))
	fmt.Println(hex.EncodeToString(skk[:32]))

	print()
}

func TestMLDSA_Sign_Verify(t *testing.T) {
	fmt.Println("**********************************************")
	fmt.Println("#Testing ML-DSA Sign and Verify")
	var mldsa MLDSA
	pubkey, privatekey, _ := mldsa.KeyGen("ML-DSA-65")
	encryptedMessage := []byte("Hello World")

	signature := mldsa.Sign(encryptedMessage, privatekey, "ML-DSA-65")
	verificationStatus, e := mldsa.Verify(pubkey, encryptedMessage, signature, "ML-DSA-65")
	if verificationStatus && e == nil {
		fmt.Println("Signature verified")
	} else {
		t.Errorf("Failed to verify signature")
	}

	encryptedMessage = []byte("Goodbye World!")
	verificationStatus, e = mldsa.Verify(pubkey, encryptedMessage, signature, "ML-DSA-65")
	if !verificationStatus && e != nil {
		fmt.Println("Wrong Signature Failed ")
	} else {
		t.Errorf("Wrong Signature is Vrified ")
	}
}

func TestMLDSA_Sign_Verify_Unmarshal(t *testing.T) {
	fmt.Println("**********************************************")
	fmt.Println("#Testing ML-DSA Sign and Verify")
	var mldsa MLDSA
	pubkey, privatekey, _ := mldsa.KeyGen("ML-DSA-65")
	encryptedMessage := []byte("Hello World")
	marshalPubKey, _ := pubkey.MarshalBinary()
	marshalSecKey, _ := privatekey.MarshalBinary()
	UmarshaledPubKey, _ := mldsa.UnmarshalPublicKey("ML-DSA-65", marshalPubKey)
	UmarshaledSecKey, _ := mldsa.UnmarshalSecretKey("ML-DSA-65", marshalSecKey)

	signature := mldsa.Sign(encryptedMessage, UmarshaledSecKey, "ML-DSA-65")
	verificationStatus, e := mldsa.Verify(UmarshaledPubKey, encryptedMessage, signature, "ML-DSA-65")
	if verificationStatus && e == nil {
		fmt.Println("Signature verified")
	} else {
		t.Errorf("Failed to verify signature")
	}

	encryptedMessage = []byte("Goodbye World!")
	verificationStatus, e = mldsa.Verify(pubkey, encryptedMessage, signature, "ML-DSA-65")
	if !verificationStatus && e != nil {
		fmt.Println("Wrong Signature Failed ")
	} else {
		t.Errorf("Wrong Signature is Vrified ")
	}
}
