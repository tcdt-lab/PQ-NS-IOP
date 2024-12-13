package gateway_verifier

import (
	b64 "encoding/base64"
	"github.com/stretchr/testify/assert"
	"test.org/protocol/pkg/gateway_verifier"

	"test.org/cryptography/pkg/asymmetric"
	"test.org/cryptography/pkg/symmetric"
	"test.org/protocol/pkg"
	"testing"
	"verifier/config"
)

const (
	PQ_DSA_SCHEME  = "ML-DSA-65"
	PQ_KEM_SCHEME  = "ML-KEM-768"
	ECC_DSA_SCHEME = "p256"
	ECC_KEM_SCHEME = "x25519"
)

func MessageUtilGenerator() pkg.MessageUtil {
	var util pkg.MessageUtil
	util.AesHandler = symmetric.AesGcm{}
	util.AsymmetricHandler = asymmetric.NewAsymmetricHandler("ECC")
	util.HmacHandler = symmetric.HMAC{}
	util.RegisterInterfacesInGob()
	return util
}
func generatePublicKeyDSA() (string, string) {
	msgUti := MessageUtilGenerator()
	secKey, pubKey, _ := msgUti.AsymmetricHandler.DSKeyGen(ECC_DSA_SCHEME)
	return secKey, pubKey
}

func generatePublicKeyKEM() (string, string) {
	msgUti := MessageUtilGenerator()
	secKey, pubKey, _ := msgUti.AsymmetricHandler.KEMKeyGen(ECC_KEM_SCHEME)
	return secKey, pubKey
}
func FakeUnEncryptedMessageDataGenerator() (pkg.MessageData, string) {
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	msgParams := gateway_verifier.GatewayVerifierKeyDistributionRequest{}
	msgutil := MessageUtilGenerator()
	dsaSecKey, dsaPubKey := generatePublicKeyDSA()
	kemSecKey, kemPubKey := generatePublicKeyKEM()

	msgParams.RequestId = 1
	msgParams.GatewayCompanyName = "company"
	msgParams.GatewayPort = "8000"
	msgParams.GatewayIP = "127.0.0.1"
	msgParams.GatewaySignatureScheme = ECC_DSA_SCHEME
	msgParams.GatewayKemScheme = ECC_KEM_SCHEME
	msgParams.GatewayPublicKeySignature = dsaPubKey
	msgParams.GatewayPublicKeyKem = kemPubKey

	msgInfo.Params = msgParams
	msgInfo.Nonce = 123
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID
	msgData.MsgInfo = msgInfo
	msgutil.SignMessageInfo(&msgData, dsaSecKey, ECC_DSA_SCHEME)
	return msgData, kemSecKey
}
func TestGatewayVerifierKeyDistributionHandler(t *testing.T) {
	var responseBytes []byte
	var kemSecKey string
	var msgDataFake pkg.MessageData
	t.Run("TestGatewayVerifierKeyDistributionHandler_check_functionality", func(t *testing.T) {
		c, err := config.ReadYaml()
		assert.NoError(t, err, "Error while reading config file")
		msgDataFake, kemSecKey = FakeUnEncryptedMessageDataGenerator()
		responseBytes, err = GatewayVerifierKeyDistributionHandler(msgDataFake, *c)
		assert.NoError(t, err, "Error while generating response")
		assert.NotNil(t, responseBytes, "Response is nil")

	})
	t.Run("TestGatewayVerifierKeyDistributionHandler_check_response", func(t *testing.T) {
		msgUtil := MessageUtilGenerator()
		msg, err := msgUtil.ConvertByteToMessage(responseBytes)
		assert.NoError(t, err, "Error while converting byte to message")
		assert.NotNil(t, msg, "Message is nil")
		assert.False(t, msg.IsEncrypted, "Message is encrypted ")
		msgDataStr := msg.Data
		msgDataByte, err := b64.StdEncoding.DecodeString(msgDataStr)
		assert.NoError(t, err, "Error while decoding message data")
		msgData, err := msgUtil.ConvertByteToMessageData(msgDataByte)
		assert.NoError(t, err, "Error while converting byte to message data")
		t.Log(msgData.MsgInfo)
		if msgData.MsgInfo.OperationTypeId == pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID {
			responseParam := msgData.MsgInfo.Params.(gateway_verifier.GatewayVerifierKeyDistributionResponse)
			_, secKey, err := msgUtil.AsymmetricHandler.KemGenerateSecretKey(kemSecKey, "", responseParam.CipherText, "ML-KEM-768")
			assert.NoError(t, err, "Error while generating shared symmetric key")
			assert.NotNil(t, secKey, "Shared symmetric key is nil")
			secKeyStr64 := msgUtil.AesHandler.ConvertKeyBytesToStr64(secKey)
			resHamc, err := msgUtil.VerifyHmac(msgData, secKeyStr64)
			assert.NoError(t, err, "Error while verifying hmac")
			assert.True(t, resHamc, "Hmac verification failed")

			t.Log(responseParam)
		} else if msgData.MsgInfo.OperationTypeId == pkg.GATEWAY_VERIFIER_OPERATION_ERROR_ID {
			responseParam := msgData.MsgInfo.Params.(gateway_verifier.GatewayVerifierKeyDistributionResponse)
			t.Log(msgData.MsgInfo.Params)
			t.Log(responseParam)
			t.Log(responseParam.OperationError)
		}
	})
}
