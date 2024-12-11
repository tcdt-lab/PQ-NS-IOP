package symmetric

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_KeyDerivation(t *testing.T) {
	keyDrivation := PBKDF2{}
	resKey := keyDrivation.KeyDerivation([]byte("password"), []byte("salt"), 4096)
	if len(resKey) != 32 {
		t.Errorf("Expected key length 32, got %d", len(resKey))
	}
	fmt.Println("Derived key:", hex.EncodeToString(resKey))
}
