package message_creator

import (
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
)

func CreateGatewayVerifierGetInfoOperationMessage(c *config.Config, verifierIp string, verifierPort string) ([]byte, error) {

	msg := pkg.Message{}
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	vDA := data_access.VerifierDA{}
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	destinationVerifier, err := vDA.GetVerifierByIpAndPort(verifierIp, verifierPort)
	nonce, err := util.GenerateNonce()
	if err != nil {
		zap.L().Error("Error while generating nonce", zap.Error(err))
		return nil, err

	}
	msgInfo.Nonce = nonce
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_GET_INFO_OPERATION_REQEST
	gatewayVerifierInitInfoOperationRequest := gateway_verifier.GatewayVerifierInitInfoOperationRequest{}
	gatewayVerifierInitInfoOperationRequest.RequestId = 1
	msgInfo.Params = gatewayVerifierInitInfoOperationRequest
	msgData.MsgInfo = msgInfo
	err = protocolUtil.GenerateHmacMsgInfo(&msgData, destinationVerifier.SymmetricKey)
	if err != nil {
		zap.L().Error("Error while generating HMAC", zap.Error(err))
		return nil, err
	}
	msgDataEnc, err := protocolUtil.EncryptMessageData(msgData, destinationVerifier.SymmetricKey)
	if err != nil {
		zap.L().Error("Error while encrypting message data", zap.Error(err))
		return nil, err
	}
	msg.IsEncrypted = true
	msg.Data = msgDataEnc
	msgByte, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		zap.L().Error("Error while converting message to byte", zap.Error(err))
		return nil, err
	}
	return msgByte, nil
}
