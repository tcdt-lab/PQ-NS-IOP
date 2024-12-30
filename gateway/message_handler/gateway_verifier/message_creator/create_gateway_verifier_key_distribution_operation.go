package message_creator

import (
	b64 "encoding/base64"
	"fmt"
	"gateway/config"
	"go.uber.org/zap"
	"test.org/protocol/pkg"

	"gateway/message_handler/util"
	"test.org/protocol/pkg/gateway_verifier"
)

func CreateGatewayVerifierKeyDistributionMessage(c *config.Config) []byte {
	msg := pkg.Message{}
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	var nonce string
	currentUSer, err := util.GetCurrentGatewayUser(c)
	if err != nil {
		fmt.Println(err)
		zap.L().Error("Error while getting current gateway user", zap.Error(err))
		return nil
	}
	// Get current gateway user
	gvKeyDistributionReq := gateway_verifier.GatewayVerifierKeyDistributionRequest{}
	gvKeyDistributionReq.GatewayIP = c.Server.Ip
	gvKeyDistributionReq.GatewayPort = c.Server.Port
	gvKeyDistributionReq.GatewaySignatureScheme = c.Security.DSAScheme
	gvKeyDistributionReq.GatewayPublicKeySignature = currentUSer.PublicKeyDsa
	gvKeyDistributionReq.GatewayPublicKeyKem = currentUSer.PublicKeyKem
	gvKeyDistributionReq.GatewayKemScheme = c.Security.KEMScheme

	nonce, err = util.GenerateNonce()
	if err != nil {
		zap.L().Error("Error while generating nonce", zap.Error(err))
		return nil
	}
	msgInfo.Nonce = nonce
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID
	msgInfo.Params = gvKeyDistributionReq
	msgData.MsgInfo = msgInfo
	err = protocolUtil.SignMessageInfo(&msgData, currentUSer.SecretKeyDsa, c.Security.DSAScheme)
	if err != nil {
		zap.L().Error("Error while signing message info", zap.Error(err))
		return nil
	}
	msgDataByte, err := protocolUtil.ConvertMessageDataToByte(msgData)
	if err != nil {
		zap.L().Error("Error while converting message data to byte", zap.Error(err))
		return nil
	}
	msg.Data = b64.StdEncoding.EncodeToString(msgDataByte)
	msg.IsEncrypted = false
	msgByte, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		zap.L().Error("Error while converting message to byte", zap.Error(err))
		return nil
	}
	return msgByte
}
