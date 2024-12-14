package error

import (
	"database/sql"
	"gateway/config"
	"gateway/data"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
)

func GenerateUnencryptedGeneralErrorResponse(err error, c config.Config, db *sql.DB, nonce string, errorCode int, currentGatewayUser data.GatewayUser) []byte {
	msg := pkg.Message{}
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	msgParams := pkg.ErrorParams{}

	msgParams.ErrorCode = errorCode
	msgParams.ErrorMessage = err.Error()
	msgInfo.Nonce = nonce
	msgInfo.Params = msgParams
	msgData.MsgInfo = msgInfo
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	protocolUtil.SignMessageInfo(&msgData, currentGatewayUser.SecretKeyDsa, c.Security.DSAScheme)
	msgDataStr, err := protocolUtil.ConvertMessageDataToB64String(msgData)
	if err != nil {
		zap.L().Error("ErrorParams while converting message data to byte", zap.Error(err))
		return nil
	}
	msg.Data = msgDataStr
	msg.IsEncrypted = false
	msg.MsgTicket = ""

	msgByte, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		zap.L().Error("ErrorParams while converting message to byte", zap.Error(err))
		return nil
	}
	return msgByte
}
