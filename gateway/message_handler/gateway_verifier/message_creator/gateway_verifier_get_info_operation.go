package message_creator

import (
	"gateway/config"
)

func CreateGatewayVerifierGetInfoOperationMessage(c *config.Config) []byte {

	//msg := pkg.Message{}
	//msgData := pkg.MessageData{}
	//msgInfo := pkg.MessageInfo{}
	//protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	//destinationVerifier, err := util.GetVerifierByPublicSigKey(c, c.BootstrapNode.PubKeySig)
	//nonce, err := util.GenerateNonce()
	//if err != nil {
	//	zap.L().Error("Error while generating nonce", zap.Error(err))
	//	return nil
	//
	//}
	//msgInfo.Nonce = nonce
	//msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_GET_INFO_OPERATION_REQEST
	//gatewayVerifierInitInfoOperationRequest := gateway_verifier.GatewayVerifierInitInfoOperationRequest{}
	//gatewayVerifierInitInfoOperationRequest.RequestId = 1
	//msgInfo.Params = gatewayVerifierInitInfoOperationRequest
	//msgData.MsgInfo = msgInfo
	//err = protocolUtil.GenerateHmacMsgInfo(&msgData, destinationVerifier.SymmetricKey)
	//if err != nil {
	//	zap.L().Error("Error while generating HMAC", zap.Error(err))
	//	return nil
	//}
	//msgDataEnc, err := protocolUtil.EncryptMessageData(msgData, destinationVerifier.SymmetricKey)
	//if err != nil {
	//	zap.L().Error("Error while encrypting message data", zap.Error(err))
	//	return nil
	//}
	//msg.IsEncrypted = true
	//msg.Data = msgDataEnc
	//msgByte, err := protocolUtil.ConvertMessageToByte(msg)
	//if err != nil {
	//	zap.L().Error("Error while converting message to byte", zap.Error(err))
	//	return nil
	//}
	return nil
}
