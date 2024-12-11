package symmetric

import (
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"go.uber.org/zap"
)

type HMAC struct {
}

func (h *HMAC) GenerateMessageMac(key []byte, message []byte) ([]byte, error) {

	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	generatedMac := mac.Sum(nil)
	zap.L().Info("Generated MAC", zap.String("mac", hex.EncodeToString(generatedMac)))
	return generatedMac, nil
}

func (h *HMAC) VerifyMessageMac(key []byte, message []byte, messageMAC []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	zap.L().Info("Expected MAC", zap.String("mac", string(expectedMAC)))
	res := hmac.Equal(messageMAC, expectedMAC)
	if !res {
		zap.L().Error("MAC verification failed", zap.String("expectedMAC", hex.EncodeToString(expectedMAC)), zap.String("messageMAC", hex.EncodeToString(messageMAC)))
		return res
	}

	zap.L().Info("MAC verification successful")
	return res
}

func (h *HMAC) ConvertHMacMsgToBase64(mac []byte) string {
	return b64.StdEncoding.EncodeToString(mac)
}

func (h *HMAC) ConvertBase64ToHMacMsg(mac string) []byte {
	macByte, _ := b64.StdEncoding.DecodeString(mac)
	return macByte
}
