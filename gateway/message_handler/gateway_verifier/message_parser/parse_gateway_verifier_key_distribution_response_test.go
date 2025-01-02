package message_parser

import (
	"gateway/config"
	"gateway/data_access"
	"github.com/stretchr/testify/assert"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"testing"

	"gateway/message_handler/util"
)

func createMockGVKeyDistrobutionResponse(t *testing.T, privateKeyDSAVerifier string) []byte {
	cfg, err := config.ReadYaml()
	assert.NoError(t, err)
	pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	gvResponse := gateway_verifier.GatewayVerifierKeyDistributionResponse{}
	_, pkKem, err := pkgUtil.AsymmetricHandler.KEMKeyGen(cfg.Security.KEMScheme)
	cipher, _, err := pkgUtil.AsymmetricHandler.KemGenerateSecretKey("", pkKem, "", cfg.Security.KEMScheme)
	assert.NoError(t, err)
	gvResponse.PublicKeyKem = pkKem
	gvResponse.CipherText = pkgUtil.AesHandler.ConvertKeyBytesToStr64(cipher)
	gvResponse.RequestId = 1
	msgInfo := pkg.MessageInfo{}
	msgInfo.Params = gvResponse
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID
	msgInfo.Nonce = "123"
	msgData := pkg.MessageData{}
	msgData.MsgInfo = msgInfo
	pkgUtil.SignMessageInfo(&msgData, privateKeyDSAVerifier, cfg.Security.DSAScheme)
	msgDataStr, err := pkgUtil.ConvertMessageDataToB64String(msgData)
	assert.NoError(t, err)
	msg := pkg.Message{}
	msg.Data = msgDataStr
	msg.IsEncrypted = false
	msg.MsgTicket = ""
	msgBytes, err := pkgUtil.ConvertMessageToByte(msg)
	assert.NoError(t, err)
	return msgBytes
}

func Test_ParseGatewayVerifierKeyDistributionResponse(t *testing.T) {
	t.Run("Test_Parse_Gateway_Verifier_KeyDistribution_Response_check_functionality", func(t *testing.T) {
		cfg, err := config.ReadYaml()
		vDA := data_access.VerifierDA{}
		assert.NoError(t, err)
		pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
		privateKeyDSAVerifier, publicKeyDSAVerifier, err := pkgUtil.AsymmetricHandler.DSKeyGen(cfg.Security.DSAScheme)
		assert.NoError(t, err)
		msgBytes := createMockGVKeyDistrobutionResponse(t, privateKeyDSAVerifier)
		verfiier, err := vDA.GetVerifier(1)
		verfiier.PublicKey = publicKeyDSAVerifier
		vDA.UpdateVerifier(verfiier)
		msgData, err := ParseGatewayVerifierKeyDistributionResponse(msgBytes, verfiier.Ip, verfiier.Port)
		assert.NoError(t, err)
		assert.NotNil(t, msgData)
		t.Log(msgData)
	})
}
