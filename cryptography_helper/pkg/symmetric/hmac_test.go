package symmetric

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
)

func TestHMAC_GenerateMessageMac(t *testing.T) {
	var h HMAC
	key := make([]byte, 32)
	for i := 0; i < 32; i++ {
		key[i] = byte(rand.Intn(255))
	}
	message := []byte("message")
	generatedMac, err := h.GenerateMessageMac(key, message)

	if err != nil {
		t.Errorf("Error in GenerateMessageMac: %v", err)
	}
	fmt.Println(hex.EncodeToString(generatedMac))
}

func TestHMAC_VerifyMessageMac(t *testing.T) {
	var h HMAC
	key := make([]byte, 32)
	for i := 0; i < 32; i++ {
		key[i] = byte(rand.Intn(255))
	}
	message := []byte("message")
	generatedMac, err := h.GenerateMessageMac(key, message)

	if err != nil {
		t.Errorf("Error in GenerateMessageMac: %v", err)
	}

	if !h.VerifyMessageMac(key, message, generatedMac) {
		t.Errorf("Failed to verify message MAC")
	}
}
