package balance_check

import (
	"database/sql"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_gateway"
)

func CreateBalanceCheckRequest(requestId int64, threshold int64, destinationIp string, destinationPort string, sharedNewKey string, ticket string, db *sql.DB, c config.Config) ([]byte, error) {
	guDa := data_access.GenerateGatewayUserDA(db)
	var msg = pkg.Message{}
	var msgInfo = pkg.MessageInfo{}
	var balanceCheckRequest = gateway_gateway.BalanceCheckRequest{}
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	cacheHandler := data_access.NewCacheHandlerDA()
	adminId, err := cacheHandler.GetUserAdminId()
	gatewayUser, err := guDa.GetGatewayUser(adminId)
	if err != nil {
		return nil, err
	}

	balanceCheckRequest.RequestId = requestId
	balanceCheckRequest.BalanceThreshold = threshold

	msgInfo.Params = balanceCheckRequest
	msgInfo.RequestId = requestId
	msgInfo.OperationTypeId = pkg.GATEWAY_GATEWAY_BALANCE_CHECK_REQUEST_ID
	msgInfo.SourceId = c.Server.Ip + ":" + c.Server.Port
	msgInfo.DestinationId = destinationIp + ":" + destinationPort

	msgInfoByte, err := protocolUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err
	}

	encMsgInfoStr, _, err := protocolUtil.EncryptMessageInfo(msgInfoByte, sharedNewKey)
	protocolUtil.SignMessageInfo(&msg, msgInfoByte, gatewayUser.SecretKeyDsa, c.Security.DSAScheme)
	msg.PublicKeySig = gatewayUser.PublicKeyDsa
	msg.MsgInfo = encMsgInfoStr
	msg.MsgTicket = ticket
	msg.IsEncrypted = true
	msgBytes, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}
