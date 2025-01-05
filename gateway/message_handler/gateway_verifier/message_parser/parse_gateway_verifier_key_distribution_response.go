package message_parser

import (
	"errors"
	"gateway/config"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
)

// Deprecated: Because of redundancy this func is no longer needed
func ParseGatewayVerifierKeyDistributionResponse(msgBytes []byte, sourceIp string, sourcePort string) (pkg.MessageData, error) {
	cfg, err := config.ReadYaml()
	pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		return pkg.MessageData{}, err
	}
	protocolUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	message, err := protocolUtil.ConvertByteToMessage(msgBytes)
	if err != nil {
		zap.L().Error("Error while converting byte to message", zap.Error(err))
		return pkg.MessageData{}, err
	}

	msgData, err := pkgUtil.ConvertB64ToMessageData(message.Data)
	if err != nil {
		zap.L().Error("Error while converting b64 to message data", zap.Error(err))
		return pkg.MessageData{}, err
	}
	if msgData.Signature != "" {
		util.CheckMessageSignature(&msgData, sourceIp, sourcePort, true, pkgUtil, cfg)
	} else {
		return pkg.MessageData{}, errors.New("Signature is missing")
	}

	return msgData, nil
}
