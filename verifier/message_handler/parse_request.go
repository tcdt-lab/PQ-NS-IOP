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

func ParseRequest(msgBytes []byte, senderIp string, senderPort string, db *sql.DB) (pkg.MessageInfo, string, error) {
	cfg, err := config.ReadYaml()

	protoUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		return pkg.MessageInfo{}, "", err
	}

	message, err := protoUtil.ConvertByteToMessage(msgBytes)
	if err != nil {
		zap.L().Error("Error while converting byte to message", zap.Error(err))
		return pkg.MessageInfo{}, "", err
	}

	msgInfo := pkg.MessageInfo{}
	var decMsgInfoBytes []byte
	gDA := data_access.GenerateGatewayDA(db)

	gatewayExists, err := gDA.IfGatewayExistByPublicKeySig(message.PublicKeySig)
	if gatewayExists {
		sourceGateway, err := gDA.GetGatewayByPublicKeySig(message.PublicKeySig)
		if err != nil {
			zap.L().Error("Error while getting gateway", zap.Error(err))
			return pkg.MessageInfo{}, "", err
		}
		if message.IsEncrypted {

			msgInfo, decMsgInfoBytes, err = protoUtil.DecryptMessageInfo(message.MsgInfo, sourceGateway.SymmetricKey)
			if err != nil {
				zap.L().Error("Error while decrypting message data", zap.Error(err))
				return pkg.MessageInfo{}, "", err
			}
		} else {
			msgInfo, decMsgInfoBytes, err = protoUtil.ConvertPlainStrToMessageInfo(message.MsgInfo)

			if err != nil {
				zap.L().Error("Error while converting b64 to message data", zap.Error(err))
				return pkg.MessageInfo{}, "", err
			}
		}
		if message.Hmac != "" {

			res, err := protoUtil.VerifyHmac(message.Hmac, decMsgInfoBytes, sourceGateway.SymmetricKey)

			if err != nil {
				zap.L().Error("Error while verifying HMAC", zap.Error(err))
				return pkg.MessageInfo{}, "", err
			}
			if !res {
				zap.L().Error("HMAC is not valid")
				return pkg.MessageInfo{}, "", errors.New("HMAC is not valid")
			}

		}
		if message.Signature != "" {
			signatureBytes, _ := b64.StdEncoding.DecodeString(message.Signature)
			res, err := protoUtil.VerifyMessageSignature(signatureBytes, decMsgInfoBytes, sourceGateway.PublicKeySig, cfg.Security.DSAScheme)
			if err != nil {
				zap.L().Error("Error while verifying signature", zap.Error(err))
				return pkg.MessageInfo{}, "", err
			}
			if !res {
				zap.L().Error("Signature is not valid")
				return pkg.MessageInfo{}, "", errors.New("Signature is not valid")
			}
		}
		return msgInfo, message.PublicKeySig, nil
	} else {
		msgInfoBytes, _ := b64.StdEncoding.DecodeString(message.MsgInfo)
		msgInfo, err = protoUtil.ConvertByteToMessageInfo(msgInfoBytes)
		if err != nil {
			zap.L().Error("Error while converting b64 to message data", zap.Error(err))
			return pkg.MessageInfo{}, "", err
		}
		return msgInfo, message.PublicKeySig, nil
	}

}
