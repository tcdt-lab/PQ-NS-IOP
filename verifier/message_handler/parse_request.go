package message_parser

import (
	"errors"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func ParseRequest(msgBytes []byte, senderIp string, senderPort string) (pkg.MessageData, error) {
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

	msgData := pkg.MessageData{}
	gDA := data_access.GenerateGatewayDA()
	defer gDA.CloseGatewayDaConnection()
	gatewayExists, err := gDA.IfGatewayExist(senderIp, senderPort)
	if gatewayExists {
		sourceGateway, err := gDA.GetGatewayByIpAndPort(senderIp, senderPort)
		if err != nil {
			zap.L().Error("Error while getting verifier", zap.Error(err))
			return pkg.MessageData{}, err
		}
		if message.IsEncrypted {
			msgData, err = protoUtil.DecryptMessageData(message.Data, sourceGateway.SymmetricKey)
			if err != nil {
				return pkg.MessageData{}, err
			}
		} else {
			msgData, err = protoUtil.ConvertB64ToMessageData(message.Data)
			if err != nil {
				zap.L().Error("Error while converting b64 to message data", zap.Error(err))
				return pkg.MessageData{}, err
			}
		}
		if msgData.Hmac != "" {
			res, err := protoUtil.VerifyHmac(msgData, sourceGateway.SymmetricKey)

			if err != nil {
				zap.L().Error("Error while verifying HMAC", zap.Error(err))
				return pkg.MessageData{}, err
			}
			if !res {
				zap.L().Error("HMAC is not valid")
				return pkg.MessageData{}, errors.New("HMAC is not valid")
			}
			return msgData, err
		} else if msgData.Signature != "" {
			res, err := protoUtil.VerifyMessageDataSignature(msgData, sourceGateway.PublicKeySig, cfg.Security.DSAScheme)
			if err != nil {
				zap.L().Error("Error while verifying signature", zap.Error(err))
				return pkg.MessageData{}, err
			}
			if !res {
				zap.L().Error("Signature is not valid")
				return pkg.MessageData{}, errors.New("Signature is not valid")
			}
		} else {
			msgData, err = protoUtil.ConvertB64ToMessageData(message.Data)
		}
		return msgData, nil
	} else {
		msgData, err = protoUtil.ConvertB64ToMessageData(message.Data)
		if err != nil {
			zap.L().Error("Error while converting b64 to message data", zap.Error(err))
			return pkg.MessageData{}, err
		}
		return msgData, nil
	}

}
