package symmetric

import (
	"encoding/hex"
	"math/rand"
	"testing"
)

func TestAesGcm_aes(t *testing.T) {

	var a AesGcm
	key := make([]byte, 32)
	plaintext := []byte("sample text for encryption")
	for j := 0; j < 32; j++ {
		key[j] = byte(rand.Intn(255))
	}

	ciphertext, err := a.Encrypt(plaintext, key)
	if err != nil {
		t.Errorf("Error in Encrypt: %v", err)
	}
	decryptedText, err := a.Decrypt(ciphertext, key)
	if err != nil {
		t.Errorf("Error in Decrypt: %v", err)
	}
	if hex.EncodeToString(decryptedText) != hex.EncodeToString(plaintext) {
		t.Errorf("Decrypted text is not equal to plaintext")
	}

}
