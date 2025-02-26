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
	msgInfoBytes, err := protocolUtil.ConvertMessageInfoToByte(msgInfo)
	hmacStr, _, err := protocolUtil.GenerateHmacMsgInfo(msgInfoBytes, destinationVerifier.SymmetricKey)

	//msgDataEnc, err := protocolUtil.EncryptMessageData(msgData, destinationVerifier.SymmetricKey)
	if err != nil {
		zap.L().Error("Error while encrypting message data", zap.Error(err))
		return nil, err
	}
	msg.IsEncrypted = false
	msg.Hmac = hmacStr
	msg.MsgInfo = b64.StdEncoding.EncodeToString(msgInfoBytes)
	msg.PublicKeySig = admin.PublicKeyDsa
	msgByte, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		zap.L().Error("Error while converting message to byte", zap.Error(err))
		return nil, err
	}
	return msgByte, nil
}
