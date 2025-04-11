package gv_ticket_issue

import (
	"crypto/rand"
	"database/sql"
	b64 "encoding/base64"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func CreateTicketIssueResponse(requestId int64, senderPubKey string, destinationIp string, destinatinPort string, db *sql.DB, c config.Config) ([]byte, error) {

	gDa := data_access.GenerateGatewayDA(db)
	vuDa := data_access.GenerateVerifierUserDA(db)
	protoclUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	var resp = gateway_verifier.GatewayVerifierTicketResponse{
		RequestId: requestId,
	}
	cacheHandelr := data_access.GenerateCacheHandlerDA()
	adminId, _ := cacheHandelr.GetUserAdminId()
	verfierUser, _ := vuDa.GetVerifierUser(adminId)
	msg := pkg.Message{}
	msgInfo := pkg.MessageInfo{}
	ticket := pkg.Ticket{}
	destGateway, err := gDa.GetGatewayByIpAndPort(destinationIp, destinatinPort)
	if err != nil {
		return nil, err
	}
	sourceGateway, err := gDa.GetGatewayByPublicKeySig(senderPubKey)
	if err != nil {
		return nil, err
	}

	generatedKeyBytes := generateKey(protoclUtil)
	generatedKeyStr := protoclUtil.AesHandler.ConvertKeyBytesToStr64(generatedKeyBytes)
	ticket.SharedKey = generatedKeyStr
	ticket.DestinationIp = destinationIp
	ticket.SourceIp = c.Server.Ip
	ticketBytes, err := protoclUtil.ConvertTicketToByte(ticket)
	if err != nil {
		return nil, err
	}

	destGtKeyByte, err := protoclUtil.AesHandler.ConvertKeyStr64ToBytes(destGateway.SymmetricKey)
	if err != nil {
		return nil, err
	}

	finalTicket, err := protoclUtil.AesHandler.Encrypt(ticketBytes, destGtKeyByte)
	if err != nil {
		return nil, err
	}

	resp.TicketKey = generatedKeyStr
	resp.TicketString = b64.StdEncoding.EncodeToString(finalTicket)

	msgInfo.RequestId = requestId
	msgInfo.OperationTypeId = pkg.GATEWAT_VERIFIER_TICKET_ISSUE_RESPONSE_ID
	msgInfo.Params = resp
	msgInfoBytes, err := protoclUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err
	}
	protoclUtil.SignMessageInfo(&msg, msgInfoBytes, verfierUser.SecretKeySig, c.Security.DSAScheme)
	encMsgInfoStr, _, err := protoclUtil.EncryptMessageInfo(msgInfoBytes, sourceGateway.SymmetricKey)
	msg.MsgInfo = encMsgInfoStr

	msg.IsEncrypted = true
	msg.PublicKeySig = verfierUser.PublicKeySig
	msgBytes, err := protoclUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}

func generateKey(protocolUtil pkg.ProtocolUtil) []byte {
	randomByteArray := make([]byte, 32)
	randomSalt := make([]byte, 32)
	rand.Read(randomByteArray)
	rand.Read(randomSalt)
	key := protocolUtil.PBKDF2Handler.KeyDerivation(randomByteArray, randomSalt, 4096)
	return key
}
