package message_parser

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"go.uber.org/zap"
	"os"
	"test.org/protocol/pkg"
	"verifier/logic"
	"verifier/message_parser/gateway_verifier/message_parser"
	"verifier/utility"

	"verifier/message_parser/util"

	"verifier/config"
	"verifier/data"
)

type MessageHandler struct {
}

func (mp *MessageHandler) HandleRequests(message []byte, senderIp string, senderPort string, c config.Config) ([]byte, error) {
	var response []byte
	cfg, err := config.ReadYaml()
	responseLogic := logic.ResponseLogic{}

	if err != nil {
		return nil, err
	}
	db, err := utility.GetDBConnection(c)
	if err != nil {
		return nil, err
	}
	msgData, err := message_parser.ParseRequest(message, senderIp, senderPort)
	if err != nil {
		return mp.GenerateGeneralErrorResponse(err, *cfg, db), err
	}

	switch msgData.MsgInfo.OperationTypeId {
	case pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID:
		response, err = responseLogic.HandleKeyDistributionResponse(message, senderIp, senderPort)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, *cfg, db), err
		}
		return response, nil

	case pkg.GATEWAY_VERIFIER_GET_INFO_OPERATION_REQEST:
		response, err = responseLogic.HandleGetInfoResponse(message, senderIp, senderPort)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, *cfg, db), err
		}
		return response, nil
	}

	errOperation := errors.New("Operation type not found")
	return mp.GenerateGeneralErrorResponse(errOperation, *cfg, db), errOperation
}

// go to util
func (mp *MessageHandler) GenerateGeneralErrorResponse(err error, c config.Config, db *sql.DB) []byte {
	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)

	msg := pkg.Message{}
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	errorParams := pkg.ErrorParams{}

	verifeirUser, err := data.GetVerifierUserByPassword(db, os.Getenv("PQ_NS_IOP_VU_PASS"))
	if err != nil {
		zap.L().Error("Error while getting verifier_verifier user", zap.Error(err))
	}
	errorParams.ErrorCode = pkg.GENERAL_ERROR
	errorParams.ErrorMessage = err.Error()
	msgInfo.Params = errorParams
	msgData.MsgInfo = msgInfo
	msgInfo.OperationTypeId = pkg.GENERAL_ERROR
	msgUtil.SignMessageInfo(&msgData, verifeirUser.SecretKeySig, c.Security.DSAScheme)
	msgDataByte, err := msgUtil.ConvertMessageDataToByte(msgData)
	if err != nil {
		zap.L().Error("Error while converting message data to byte", zap.Error(err))
	}
	msg.Data = b64.StdEncoding.EncodeToString(msgDataByte)
	msg.IsEncrypted = false
	msg.MsgTicket = ""
	msgByte, err := msgUtil.ConvertMessageToByte(msg)
	return msgByte
}
