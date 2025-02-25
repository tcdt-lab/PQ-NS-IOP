package message_handler

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func ParseRequest(msgBytes []byte, senderIp string, senderPort string, db *sql.DB) (pkg.MessageData, string, error) {
	cfg, err := config.ReadYaml()

	protoUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		return pkg.MessageData{}, "", err
	}

	message, err := protoUtil.ConvertByteToMessage(msgBytes)
	if err != nil {
		zap.L().Error("Error while converting byte to message", zap.Error(err))
		return pkg.MessageData{}, "", err
	}

	msgData := pkg.MessageData{}

	gDA := data_access.GenerateGatewayDA(db)

	gatewayExists, err := gDA.IfGatewayExistByPublicKeySig(message.PublicKeySig)
	if gatewayExists {
		sourceGateway, err := gDA.GetGatewayByPublicKeySig(message.PublicKeySig)
		if err != nil {
			zap.L().Error("Error while getting gateway", zap.Error(err))
			return pkg.MessageData{}, "", err
		}
		if message.IsEncrypted {

			msgData, err = protoUtil.DecryptMessageData(message.Data, sourceGateway.SymmetricKey)
			if err != nil {
				zap.L().Error("Error while decrypting message data", zap.Error(err))
				return pkg.MessageData{}, "", err
			}
		} else {
			msgData, err = protoUtil.ConvertB64ToMessageData(message.Data)
			msgInfoByte, _ := protoUtil.ConvertMessageInfoToByte(msgData.MsgInfo)
			messageInfoStr := b64.StdEncoding.EncodeToString(msgInfoByte)
			zap.L().Info("message info Str", zap.String("message info", messageInfoStr))
			if err != nil {
				zap.L().Error("Error while converting b64 to message data", zap.Error(err))
				return pkg.MessageData{}, "", err
			}
		}
		if msgData.Hmac != "" {
			protoUtil.GenerateHmacMsgInfo(&msgData, sourceGateway.SymmetricKey)

			protoUtil.GenerateHmacMsgInfo(&msgData, sourceGateway.SymmetricKey)

			res, err := protoUtil.VerifyHmac(msgData, sourceGateway.SymmetricKey)

			if err != nil {
				zap.L().Error("Error while verifying HMAC", zap.Error(err))
				return pkg.MessageData{}, "", err
			}
			if !res {
				zap.L().Error("HMAC is not valid")
				return pkg.MessageData{}, "", errors.New("HMAC is not valid")
			}
			return pkg.MessageData{}, "", err
		} else if msgData.Signature != "" {
			res, err := protoUtil.VerifyMessageDataSignature(msgData, sourceGateway.PublicKeySig, cfg.Security.DSAScheme)
			if err != nil {
				zap.L().Error("Error while verifying signature", zap.Error(err))
				return pkg.MessageData{}, "", err
			}
			if !res {
				zap.L().Error("Signature is not valid")
				return pkg.MessageData{}, "", errors.New("Signature is not valid")
			}
		} else {
			msgData, err = protoUtil.ConvertB64ToMessageData(message.Data)
		}
		return msgData, message.PublicKeySig, nil
	} else {
		msgData, err = protoUtil.ConvertB64ToMessageData(message.Data)
		if err != nil {
			zap.L().Error("Error while converting b64 to message data", zap.Error(err))
			return pkg.MessageData{}, "", err
		}
		return msgData, message.PublicKeySig, nil
	}

}
