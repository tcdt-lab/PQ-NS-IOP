package symmetric

import (
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
)

type AesGcm struct {
}

func (a *AesGcm) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		zap.L().Error("Failed to create AES cipher", zap.Error(err), zap.String("key", hex.EncodeToString(key)))
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesCipher)

	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	for i := 0; i < gcm.NonceSize(); i++ {

		nonce[i] = byte(rand.Intn(255))
	}
	encText := gcm.Seal(nonce, nonce, plaintext, nil)
	zap.L().Info("Text is encrypted successfully")
	return encText, nil
}

func (a *AesGcm) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		zap.L().Error("Failed to create AES cipher", zap.Error(err), zap.String("key", hex.EncodeToString(key)))
		return nil, err
	}
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, nil
}

func (a *AesGcm) ConvertKeyBytesToStr64(key []byte) string {
	return b64.StdEncoding.EncodeToString(key)
}

func (a *AesGcm) ConvertKeyStr64ToBytes(key string) ([]byte, error) {
	return b64.StdEncoding.DecodeString(key)
}
