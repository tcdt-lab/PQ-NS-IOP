package message_handler

import (
	"database/sql"
	"gateway/config"
	"gateway/message_handler/balance_check"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
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
	msgInfo, _, err, ticketKey := ParseMessage(message, senderIp, senderPort, mp.db)
	switch msgInfo.OperationTypeId {
	case pkg.GATEWAY_GATEWAY_BALANCE_CHECK_REQUEST_ID:
		response, err = balance_check.CreateBalanceCheckResponse(ticketKey, msgInfo.RequestId, cfg, mp.db)
		if err != nil {
			zap.L().Error("Error while creating balance check response", zap.Error(err))
			return nil, err
		}

	}
	return response, nil
}
