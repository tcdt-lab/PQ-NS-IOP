package message_creator

import (
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_parser/util"
)

func CreateGatewayVerifierKeyDistributionResponse(cipherTextStr string) ([]byte, error) {
	cfg, err := config.ReadYaml()
	if err != nil {
		return nil, err
	}

	vuda := data_access.VerifierUserDA{}
	vuda.GetVerifierUser(1)
	protoUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	msg := pkg.Message{}
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	adminVerifierUSer, err := vuda.GetAdminVerifierUser()
	if err != nil {
		return nil, err
	}
	gvKeyDistributionResponse := gateway_verifier.GatewayVerifierKeyDistributionResponse{}
	gvKeyDistributionResponse.CipherText = cipherTextStr
	gvKeyDistributionResponse.PublicKeyKem = adminVerifierUSer.SecretKeyKem
	msgInfo.Params = gvKeyDistributionResponse
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID
	msgInfo.Nonce = "123"
	msgData.MsgInfo = msgInfo

	protoUtil.GenerateHmacMsgInfo(&msgData, adminVerifierUSer.SymmetricKey)
	msgDataStr, err := protoUtil.ConvertMessageDataToB64String(msgData)
	if err != nil {
		return nil, err
	}
	msg.Data = msgDataStr
	msg.IsEncrypted = false
	msg.MsgTicket = ""
	msgBytes, err := protoUtil.ConvertMessageToByte(msg)
	return msgBytes, nil
}
