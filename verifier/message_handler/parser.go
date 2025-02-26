package message_handler

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"go.uber.org/zap"
	"os"
	"strconv"
	"test.org/protocol/pkg"
	"verifier/message_handler/gateway_verifier/message_creator"
	"verifier/message_handler/key_distribution"

	"verifier/message_handler/util"

	"verifier/config"
	"verifier/data"
)

type MessageHandler struct {
	db *sql.DB
}

func GenerateNewMessageHandler(database *sql.DB) MessageHandler {
	var msgHandler = MessageHandler{db: database}
	return msgHandler
}

func (mp *MessageHandler) HandleRequests(message []byte, senderIp string, senderPort string, c config.Config) ([]byte, error) {

	zap.L().Info("handleing a request")
	var response []byte
	cfg, err := config.ReadYaml()

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	msgInfo, senderPubKey, err := ParseRequest(message, senderIp, senderPort, mp.db)
	if err != nil {
		zap.L().Error("Error while parsing request", zap.Error(err))
		return mp.GenerateGeneralErrorResponse(err, *cfg), err
	}

	switch msgInfo.OperationTypeId {
	case pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID:
		zap.L().Info("Handling Key Distribution Request", zap.String("sender id", senderIp), zap.String("sender port", senderPort), zap.String("req ID", strconv.FormatInt(msgInfo.RequestId, 10)))
		response, err = mp.HandleKeyDistributionResponse(msgInfo, senderIp, senderPort)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, *cfg), err
		}
		return response, nil

	case pkg.GATEWAY_VERIFIER_GET_INFO_OPERATION_REQEST:
		zap.L().Info("Handling Get Info Request", zap.String("sender id", senderIp), zap.String("sender port", senderPort), zap.String("req ID", strconv.FormatInt(msgInfo.RequestId, 10)))
		response, err = mp.HandleGetInfoResponse(msgInfo.RequestId, senderPubKey)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, *cfg), err
		}
		return response, nil
	}

	errOperation := errors.New("Operation type not found")
	return mp.GenerateGeneralErrorResponse(errOperation, *cfg), errOperation
}

// go to util
func (mp *MessageHandler) GenerateGeneralErrorResponse(err error, c config.Config) []byte {
	zap.L().Error("Generating Error Response", zap.Error(err))
	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)

	msg := pkg.Message{}

	msgInfo := pkg.MessageInfo{}
	errorParams := pkg.ErrorParams{}

	verifeirUser, err := data.GetVerifierUserByPassword(mp.db, os.Getenv("PQ_NS_IOP_VU_PASS"))

	if err != nil {
		zap.L().Error("Error while getting verifier_verifier user", zap.Error(err))
	}
	errorParams.ErrorCode = pkg.GENERAL_ERROR
	errorParams.ErrorMessage = err.Error()
	msgInfo.Params = errorParams
	msgInfo.OperationTypeId = pkg.GENERAL_ERROR
	msgInfoBytes, _ := msgUtil.ConvertMessageInfoToByte(msgInfo)
	msgInfoStr := b64.StdEncoding.EncodeToString(msgInfoBytes)
	msg.MsgInfo = msgInfoStr
	msgUtil.SignMessageInfo(&msg, msgInfoBytes, verifeirUser.SecretKeySig, c.Security.DSAScheme)
	msg.IsEncrypted = false
	msg.MsgTicket = ""
	msgByte, err := msgUtil.ConvertMessageToByte(msg)
	return msgByte
}

func (mp *MessageHandler) HandleKeyDistributionResponse(msgData pkg.MessageInfo, senderIp string, senderPort string) ([]byte, error) {

	cipherText, err := key_distribution.ApplyGatewayVerifierKeyDistributionRequest(msgData, mp.db)
	if err != nil {
		zap.L().Error("Error while applying key distribution request", zap.Error(err))
		return nil, err
	}
	res, err := key_distribution.CreateGatewayVerifierKeyDistributionResponse(cipherText, mp.db)
	if err != nil {
		zap.L().Error("Error while creating key distribution response", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (mp *MessageHandler) HandleGetInfoResponse(reqId int64, senderPubKey string) ([]byte, error) {

	res, err := message_creator.CreateGateVerifierGetInfoResponse(senderPubKey, reqId, mp.db)
	if err != nil {
		zap.L().Error("Error while creating get info response", zap.Error(err))
		return nil, err
	}
	return res, nil
}
