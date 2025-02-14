package message_handler

import (
	"database/sql"
	"errors"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
)

func ParseGatewayVerifierResponse(msgBytes []byte, senderIp string, senderPort string, db *sql.DB) (pkg.MessageData, error) {
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
	vDA := data_access.GenerateVerifierDA(db)
	senderVerfier, err := vDA.GetVerifierByIpAndPort(senderIp, senderPort)
	if err != nil {
		zap.L().Error("Error while getting verifier", zap.Error(err))
		return pkg.MessageData{}, err
	}
	if message.IsEncrypted {
		msgData, err = protoUtil.DecryptMessageData(message.Data, senderVerfier.SymmetricKey)
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
		if msgData.MsgInfo.OperationTypeId == pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID {
			return msgData, nil
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
	} else if msgData.Signature != "" {
		res, err := protoUtil.VerifyMessageDataSignature(msgData, senderVerfier.PublicKey, cfg.Security.DSAScheme)
		if err != nil {
			zap.L().Error("Error while verifying signature", zap.Error(err))
			return pkg.MessageData{}, err
		}
		if !res {
			zap.L().Error("Signature is not valid")
			return pkg.MessageData{}, errors.New("Signature is not valid")
		}
		return msgData, nil
	} else {
		zap.L().Error("HMAC and Signature is missing")
		return pkg.MessageData{}, errors.New("HMAC and Signature is missing")
	}

}
