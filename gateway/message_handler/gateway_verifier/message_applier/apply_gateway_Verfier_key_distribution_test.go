package message_applier

import (
	b64 "encoding/base64"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"github.com/stretchr/testify/assert"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"testing"
)

func createFakeMsgData(t *testing.T, pubKeyKemGateway string) (pkg.MessageData, string) {
	cfg, err := config.ReadYaml()
	assert.NoError(t, err)
	//vDA := data_access.VerifierDA{}
	//verifier, err := vDA.GetVerifierByIpAndPort(cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)
	pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	cipherText, sharedkey, err := pkgUtil.AsymmetricHandler.KemGenerateSecretKey("", pubKeyKemGateway, "", cfg.Security.KEMScheme)
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	gvResponse := gateway_verifier.GatewayVerifierKeyDistributionResponse{}
	gvResponse.PublicKeyKem = pubKeyKemGateway

	gvResponse.CipherText = b64.StdEncoding.EncodeToString(cipherText)
	msgInfo.Params = gvResponse
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID
	msgInfo.Nonce = "123"

	msgData.MsgInfo = msgInfo
	//pkgUtil.SignMessageInfo(&msgData, , cfg.Security.DSAScheme)
	msgDataStr, err := pkgUtil.ConvertMessageDataToB64String(msgData)
	assert.NoError(t, err)
	msg := pkg.Message{}
	msg.Data = msgDataStr
	msg.IsEncrypted = false
	msg.MsgTicket = ""
	return msgData, pkgUtil.AesHandler.ConvertKeyBytesToStr64(sharedkey)
}

func TestApplyGatewayVerifierKeyDistributionResponse(t *testing.T) {

	cfg, err := config.ReadYaml()
	pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)

	sKeyKem, pkKem, err := pkgUtil.AsymmetricHandler.KEMKeyGen(cfg.Security.KEMScheme)
	gud := data_access.GatewayUserDA{}
	gateWayUSer, err := gud.GetGatewayUser(1)
	gateWayUSer.PublicKeyKem = pkKem
	gateWayUSer.SecretKeyKem = sKeyKem
	gud.UpdateGatewayUser(gateWayUSer)

	msgData, sharedKey := createFakeMsgData(t, gateWayUSer.PublicKeyKem)
	t.Log(sharedKey)
	err = ApplyGatewayVerifierKeyDistributionResponse(msgData)
	assert.NoError(t, err)
}
