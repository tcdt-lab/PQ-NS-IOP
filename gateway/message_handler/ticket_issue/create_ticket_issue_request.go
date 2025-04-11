package ticket_issue

import (
	"database/sql"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
)

func CreateTicketIssueRequest(requestId int64, destinationServerIP string, destinationServerPort string, db *sql.DB, c config.Config) ([]byte, error) {
	protoclUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	guDa := data_access.GenerateGatewayUserDA(db)
	vDa := data_access.GenerateVerifierDA(db)
	cacheHandler := data_access.NewCacheHandlerDA()
	adminId, err := cacheHandler.GetUserAdminId()
	if err != nil {
		return nil, err
	}
	gatewayUser, err := guDa.GetGatewayUser(adminId)

	if err != nil {
		return nil, err
	}
	boottrapVerfier, err := vDa.GetVerifierByIpAndPort(c.BootstrapNode.Ip, c.BootstrapNode.Port)
	if err != nil {
		return nil, err
	}

	ticketReqParam := gateway_verifier.GatewayVerifierTicketRequest{
		RequestId:             requestId,
		SourceServerIP:        c.Server.Ip,
		SourceServerPort:      c.Server.Port,
		DestinationServerIP:   destinationServerIP,
		DestinationServerPort: destinationServerPort,
	}

	msgInfo := pkg.MessageInfo{}
	message := pkg.Message{}
	msgInfo.RequestId = requestId
	msgInfo.OperationTypeId = pkg.GATEWAT_VERIFIER_TICKET_ISSUE_REQUEST_ID
	msgInfo.Params = ticketReqParam

	msgInfoByte, err := protoclUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err
	}
	protoclUtil.SignMessageInfo(&message, msgInfoByte, gatewayUser.SecretKeyDsa, c.Security.DSAScheme)
	message.PublicKeySig = gatewayUser.PublicKeyDsa
	encmsgInfo, _, err := protoclUtil.EncryptMessageInfo(msgInfoByte, boottrapVerfier.SymmetricKey)
	if err != nil {
		return nil, err
	}
	message.MsgInfo = encmsgInfo
	message.IsEncrypted = true

	messageBytes, err := protoclUtil.ConvertMessageToByte(message)
	if err != nil {
		return nil, err
	}
	return messageBytes, nil
}
