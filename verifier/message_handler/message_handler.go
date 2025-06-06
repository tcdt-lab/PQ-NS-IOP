package message_handler

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"verifier/message_handler/gateway_verifier/gv_balance_verification"
	"verifier/message_handler/gateway_verifier/gv_get_info"
	key_distribution2 "verifier/message_handler/gateway_verifier/gv_key_distribution"
	"verifier/message_handler/gateway_verifier/gv_ticket_issue"
	"verifier/message_handler/util"
	"verifier/message_handler/verifier_verifier/vv_get_info"
	"verifier/message_handler/verifier_verifier/vv_key_distribution"

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

func (mp *MessageHandler) HandleRequests(message []byte, senderIp string, senderPort string, c config.Config, mutex *sync.Mutex) ([]byte, error) {

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
		response, err = mp.GV_HandleKeyDistributionResponse(msgInfo, senderIp, senderPort)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, *cfg), err
		}
		return response, nil

	case pkg.GATEWAY_VERIFIER_GET_INFO_OPERATION_REQEST:
		zap.L().Info("Handling Get Info Request", zap.String("sender id", senderIp), zap.String("sender port", senderPort), zap.String("req ID", strconv.FormatInt(msgInfo.RequestId, 10)))
		response, err = mp.GV_HandleGetInfoResponse(msgInfo.RequestId, senderPubKey)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, *cfg), err
		}
		return response, nil

	case pkg.VERIFIER_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID:
		zap.L().Info("Handling Key Distribution Request", zap.String("sender id", senderIp), zap.String("sender port", senderPort), zap.String("req ID", strconv.FormatInt(msgInfo.RequestId, 10)))
		response, err = mp.VV_HandleKeyDistributionResponse(msgInfo, senderIp, senderPort, cfg)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, *cfg), err
		}
		return response, nil

	case pkg.VERIFIER_VERIFIER_GET_INFO_OPERATION_REQEST_ID:
		zap.L().Info("Handling Get Info Request", zap.String("sender id", senderIp), zap.String("sender port", senderPort), zap.String("req ID", strconv.FormatInt(msgInfo.RequestId, 10)))
		response, err = mp.VV_HanldeGetInfoResponse(msgInfo.RequestId, senderPubKey, cfg)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, *cfg), err
		}

	case pkg.GATEWAT_VERIFIER_TICKET_ISSUE_REQUEST_ID:
		zap.L().Info("Handling Ticket Issue Request", zap.String("sender id", senderIp), zap.String("sender port", senderPort), zap.String("req ID", strconv.FormatInt(msgInfo.RequestId, 10)))
		response, err = mp.GV_HandleTicketIssueRequest(msgInfo, senderPubKey, cfg)

	case pkg.GATEWAY_VERIFIER_BALANCE_VERIFICATION_REQUEST_ID:

		zap.L().Info("Handling Balance Verification Request", zap.String("sender id", senderIp), zap.String("sender port", senderPort), zap.String("req ID", strconv.FormatInt(msgInfo.RequestId, 10)))
		response, err = mp.GV_HandleBalanceVerificationRequest(msgInfo, senderPubKey, cfg, mutex)

	}

	return response, nil
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

func (mp *MessageHandler) GV_HandleKeyDistributionResponse(msgData pkg.MessageInfo, senderIp string, senderPort string) ([]byte, error) {

	cipherText, err := key_distribution2.ApplyGatewayVerifierKeyDistributionRequest(msgData, mp.db)
	if err != nil {
		zap.L().Error("Error while applying key distribution request", zap.Error(err))
		return nil, err
	}
	res, err := key_distribution2.CreateGatewayVerifierKeyDistributionResponse(cipherText, mp.db)
	if err != nil {
		zap.L().Error("Error while creating key distribution response", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (mp *MessageHandler) GV_HandleGetInfoResponse(reqId int64, senderPubKey string) ([]byte, error) {

	res, err := gv_get_info.CreateGateVerifierGetInfoResponse(senderPubKey, reqId, mp.db)
	if err != nil {
		zap.L().Error("Error while creating get info response", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (mp *MessageHandler) VV_HandleKeyDistributionResponse(msgData pkg.MessageInfo, senderIp string, senderPort string, cfg *config.Config) ([]byte, error) {

	cipherText, err := vv_key_distribution.ApplyKeyDistributionRequest(mp.db, cfg, msgData)
	if err != nil {
		zap.L().Error("Error while applying key distribution request", zap.Error(err))
		return nil, err
	}
	res, err := vv_key_distribution.CreateKeyDistributionResponse(mp.db, cfg, cipherText, msgData.RequestId)
	if err != nil {
		zap.L().Error("Error while creating key distribution response", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (mp *MessageHandler) VV_HanldeGetInfoResponse(reqId int64, senderPubKey string, config *config.Config) ([]byte, error) {

	res, err := vv_get_info.CreateGetInfoResponse(config, reqId, mp.db, senderPubKey)
	if err != nil {
		zap.L().Error("Error while creating get info response", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (mp *MessageHandler) GV_HandleTicketIssueRequest(info pkg.MessageInfo, key string, cfg *config.Config) ([]byte, error) {
	ticketRequestParams := info.Params.(gateway_verifier.GatewayVerifierTicketRequest)
	res, err := gv_ticket_issue.CreateTicketIssueResponse(ticketRequestParams.RequestId, key, ticketRequestParams.DestinationServerIP, ticketRequestParams.DestinationServerPort, mp.db, *cfg)
	if err != nil {
		zap.L().Error("Error while creating ticket issue response", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (mp *MessageHandler) GV_HandleBalanceVerificationRequest(info pkg.MessageInfo, key string, cfg *config.Config, mutex *sync.Mutex) ([]byte, error) {
	verificationReq := info.Params.(gateway_verifier.VerificationRequest)

	res, err := gv_balance_verification.CreateBalanceVerificationResponse(verificationReq.Proof, verificationReq.PublicInputs, info.RequestId, key, mp.db, *cfg, mutex)

	if err != nil {
		zap.L().Error("Error while creating balance verification response", zap.Error(err))
		return nil, err
	}
	return res, nil
}
