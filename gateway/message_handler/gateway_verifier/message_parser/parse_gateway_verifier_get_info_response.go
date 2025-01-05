package message_parser

import (
	"errors"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
)

// Deprecated: Because of redundancy this func is no longer needed
func ParseGatewayVerifierGetInfoResponse(msgBytes []byte, senderIp string, senderPort string) (pkg.MessageData, error) {
	cfg, err := config.ReadYaml()

	protoUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		return pkg.MessageData{}, err
	}

	message, err := protoUtil.ConvertByteToMessage(msgBytes)
	if err != nil {
		zap.L().Error("Error while converting byte to message", zap.Error(err))
		return pkg.MessageData{}, err
	}
	if message.IsEncrypted {
		msgData := pkg.MessageData{}
		vDA := data_access.VerifierDA{}
		senderVerfier, err := vDA.GetVerifierByIpAndPort(senderIp, senderPort)
		if err != nil {
			zap.L().Error("Error while getting verifier", zap.Error(err))
			return pkg.MessageData{}, err
		}
		msgData, err = protoUtil.DecryptMessageData(message.Data, senderVerfier.SymmetricKey)
		if err != nil {
			return pkg.MessageData{}, err
		}
		if msgData.Hmac == "" {
			zap.L().Error("HMAC is missing")
			return pkg.MessageData{}, errors.New("HMAC is missing")
		}
		res, err := protoUtil.VerifyHmac(msgData, senderVerfier.SymmetricKey)
		if err != nil {
			zap.L().Error("Error while verifying HMAC", zap.Error(err))
			return pkg.MessageData{}, err
		}
		if !res {
			zap.L().Error("HMAC is not valid")
			return pkg.MessageData{}, errors.New("HMAC is not valid")
		}
		return msgData, err
	} else {
		zap.L().Error("Message is not encrypted")
		return pkg.MessageData{}, errors.New("Message is not encrypted")
	}

}
