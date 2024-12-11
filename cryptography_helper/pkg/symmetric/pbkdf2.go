package symmetric

import (
	"crypto/rand"
	"crypto/sha256"
	"go.uber.org/zap"
	"golang.org/x/crypto/pbkdf2"
)

type PBKDF2 struct {
}

func (pb *PBKDF2) KeyDerivation(password []byte, salt []byte, iterations int) []byte {
	dk := pbkdf2.Key(password, salt, iterations, 32, sha256.New)
	return dk
}

func (pb *PBKDF2) GeneratingSalt(length int) ([]byte, error) {
	secret := make([]byte, length)
	if _, err := rand.Read(secret); err != nil {
		zap.L().Error("Failed to generate salt", zap.Error(err))
		return nil, err
	}
	return secret, nil
}
