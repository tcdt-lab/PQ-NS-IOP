package get_init_information

import (
	"database/sql"
	b64 "encoding/base64"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
)

func CreateGatewayVerifierGetInfoOperationMessage(adminId int64, requestId int64, c *config.Config, verifierIp string, verifierPort string, db *sql.DB) ([]byte, error) {

	msg := pkg.Message{}
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	vDA := data_access.GenerateVerifierDA(db)
	gtewayUser := data_access.GenerateGatewayUserDA(db)

	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	destinationVerifier, err := vDA.GetVerifierByIpAndPort(verifierIp, verifierPort)
	nonce, err := util.GenerateNonce()
	if err != nil {
		zap.L().Error("Error while generating nonce", zap.Error(err))
		return nil, err

	}
	admin, _ := gtewayUser.GetGatewayUser(adminId)
	msgInfo.Nonce = nonce
	msgInfo.RequestId = requestId
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_GET_INFO_OPERATION_REQEST
	gatewayVerifierInitInfoOperationRequest := gateway_verifier.GatewayVerifierInitInfoOperationRequest{}
	gatewayVerifierInitInfoOperationRequest.RequestId = requestId
	msgInfo.Params = gatewayVerifierInitInfoOperationRequest
	msgData.MsgInfo = msgInfo
	err = protocolUtil.GenerateHmacMsgInfo(&msgData, destinationVerifier.SymmetricKey)
	zap.L().Info("HMAC generated", zap.String("HMAC", msgData.Hmac))
	zap.L().Info("Request ID is genrated", zap.Int64("Request ID", requestId))
	zap.L().Info("Nonce is generated", zap.String("Nonce", nonce))
	zap.L().Info("Operation Type ID", zap.Int("Operation Type ID", pkg.GATEWAY_VERIFIER_GET_INFO_OPERATION_REQEST))
	zap.L().Info("Symmetric Key ", zap.String("Symmetric Key", destinationVerifier.SymmetricKey))
	msgInfoByte, _ := protocolUtil.ConvertMessageInfoToByte(msgData.MsgInfo)
	messageInfoStr := b64.StdEncoding.EncodeToString(msgInfoByte)
	zap.L().Info("message info Str", zap.String("message info", messageInfoStr))
	if err != nil {
		zap.L().Error("Error while generating HMAC", zap.Error(err))
		return nil, err
	}
	//msgDataEnc, err := protocolUtil.EncryptMessageData(msgData, destinationVerifier.SymmetricKey)
	if err != nil {
		zap.L().Error("Error while encrypting message data", zap.Error(err))
		return nil, err
	}
	msg.IsEncrypted = false
	msg.Data, _ = protocolUtil.ConvertMessageDataToB64String(msgData)
	msg.PublicKeySig = admin.PublicKeyDsa
	msgByte, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		zap.L().Error("Error while converting message to byte", zap.Error(err))
		return nil, err
	}
	zap.L().Info("string msg data", zap.String("data", msg.Data))
	zap.L().Info("msg bytes", zap.String("byte", string(msgByte)))
	return msgByte, nil
}
