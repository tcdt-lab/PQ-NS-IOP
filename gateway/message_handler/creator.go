package message_handler

import (
	"gateway/config"
	"gateway/message_handler/gateway_verifier"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
)

type MessageCreator struct {
}

func (msgCreator *MessageCreator) CreateMessage(operationCode int, c *config.Config) []byte {
	switch operationCode {
	case pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID:
		return gateway_verifier.CreateGatewayVerifierKeyDistributionMessage(c)
	}
	zap.L().Error("Operation code not found")
	return nil
}
