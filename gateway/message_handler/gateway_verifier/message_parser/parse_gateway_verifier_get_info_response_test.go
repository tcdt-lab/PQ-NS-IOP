package message_parser

import (
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"github.com/stretchr/testify/assert"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"testing"
)

func generateGatewayVerifierGetInfoResponseMessage(cfg config.Config) ([]byte, error) {
	msg := pkg.Message{}
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	var gatewaysInfo []gateway_verifier.GatewayVerifierInitInfoStructureGateway
	var verifiersInfo []gateway_verifier.GatewayVerifierInitInfoStructureVerifier
	var initReponse gateway_verifier.GatewayVerifierInitInfoOperationResponse
	var gatewayInfo1 gateway_verifier.GatewayVerifierInitInfoStructureGateway
	var gatewayInfo2 gateway_verifier.GatewayVerifierInitInfoStructureGateway
	var verifierInfo1 gateway_verifier.GatewayVerifierInitInfoStructureVerifier
	var verifierInfo2 gateway_verifier.GatewayVerifierInitInfoStructureVerifier
	vDA := data_access.VerifierDA{}

	gatewayInfo1.GatewayIpAddress = "127.0.0.1"
	gatewayInfo1.GatewayPort = "8080"
	gatewayInfo1.GatewayPublicKeySignature = ""
	gatewayInfo1.KemScheme = cfg.Security.KEMScheme
	gatewayInfo1.SigScheme = cfg.Security.DSAScheme
	gatewayInfo1.GatewayPublicKeyKem = ""

	gatewayInfo2.GatewayIpAddress = "127.0.0.1"
	gatewayInfo2.GatewayPort = "8081"
	gatewayInfo2.GatewayPublicKeySignature = ""
	gatewayInfo2.KemScheme = cfg.Security.KEMScheme
	gatewayInfo2.SigScheme = cfg.Security.DSAScheme
	gatewayInfo2.GatewayPublicKeyKem = ""

	gatewaysInfo = append(gatewaysInfo, gatewayInfo1)
	gatewaysInfo = append(gatewaysInfo, gatewayInfo2)

	verifierInfo1.VerifierIpAddress = "127.0.0.1"
	verifierInfo1.VerifierPort = "8070"
	verifierInfo1.VerifierPublicKeySignature = ""
	verifierInfo1.SigScheme = cfg.Security.DSAScheme
	verifierInfo1.TrustScore = 0.5
	verifierInfo1.VerifierPublicKeyKem = ""
	verifierInfo1.IsInCommittee = true

	verifierInfo2.VerifierIpAddress = "127.0.0.1"
	verifierInfo2.VerifierPort = "8071"
	verifierInfo2.VerifierPublicKeySignature = ""
	verifierInfo2.SigScheme = cfg.Security.DSAScheme
	verifierInfo2.TrustScore = 0.5
	verifierInfo2.VerifierPublicKeyKem = ""
	verifierInfo2.IsInCommittee = true

	verifiersInfo = append(verifiersInfo, verifierInfo2)

	initReponse.CurrentVerifierInfo = verifierInfo1
	initReponse.GatewaysList = gatewaysInfo
	initReponse.VerifiersList = verifiersInfo
	initReponse.OperationError = ""

	msgInfo.Params = initReponse
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_GET_INFO_OPERATION_RESPONSE
	msgInfo.Nonce = "123"
	msgData.MsgInfo = msgInfo
	pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)

	bootsrapVerifier, err := vDA.GetVerifierByIpAndPort(cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)
	if err != nil {
		return nil, err
	}
	symKey := bootsrapVerifier.SymmetricKey
	err = pkgUtil.GenerateHmacMsgInfo(&msgData, symKey)
	if err != nil {
		return nil, err
	}
	msgDataStr, err := pkgUtil.EncryptMessageData(msgData, symKey)
	if err != nil {
		return nil, err
	}
	msg.Data = msgDataStr
	msg.IsEncrypted = true
	msg.MsgTicket = ""
	msgBytes, err := pkgUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil

}

func TestParseGatewayVerifierGetInfoResponse(t *testing.T) {
	t.Run("Test_Parse_Gateway_Verifier_Get_Info_Response_check_functionality", func(t *testing.T) {
		cfg, err := config.ReadYaml()
		assert.NoError(t, err)
		msgBytes, err := generateGatewayVerifierGetInfoResponseMessage(*cfg)
		msgData, err := ParseGatewayVerifierGetInfoResponse(msgBytes, cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)
		assert.NoError(t, err)
		assert.NotNil(t, msgData)
		t.Log(msgData)
	})
}
