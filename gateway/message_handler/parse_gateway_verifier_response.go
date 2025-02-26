package message_handler

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
)

func ParseGatewayVerifierResponse(msgBytes []byte, senderIp string, senderPort string, db *sql.DB) (pkg.MessageInfo, error) {
	cfg, err := config.ReadYaml()
	var msgInfo = pkg.MessageInfo{}
	vDA := data_access.GenerateVerifierDA(db)
	var decMsgInfoBytes []byte

	protoUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		return pkg.MessageInfo{}, err
	}

	message, err := protoUtil.ConvertByteToMessage(msgBytes)
	if err != nil {
		zap.L().Error("Error while converting byte to message", zap.Error(err))
		return pkg.MessageInfo{}, err
	}

	senderVerfier, err := vDA.GetVerifierByIpAndPort(senderIp, senderPort)
	if err != nil {
		zap.L().Error("Error while getting verifier", zap.Error(err))
		return pkg.MessageInfo{}, err
	}
	if message.IsEncrypted {
		msgInfo, decMsgInfoBytes, err = protoUtil.DecryptMessageInfo(message.MsgInfo, senderVerfier.SymmetricKey)
		if err != nil {
			return pkg.MessageInfo{}, err
		}
	} else {
		msgInfo, decMsgInfoBytes, err = protoUtil.ConvertPlainStrToMessageInfo(message.MsgInfo)
		if err != nil {
			zap.L().Error("Error while converting b64 to message data", zap.Error(err))
			return pkg.MessageInfo{}, err
		}
	}
	if message.Hmac != "" {
		if msgInfo.OperationTypeId == pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID {
			return msgInfo, nil
		}
		res, err := protoUtil.VerifyHmac(message.Hmac, decMsgInfoBytes, senderVerfier.SymmetricKey)

		if err != nil {
			zap.L().Error("Error while verifying HMAC", zap.Error(err))
			return pkg.MessageInfo{}, err
		}
		if !res {
			zap.L().Error("HMAC is not valid")
			return pkg.MessageInfo{}, errors.New("HMAC is not valid")
		}
		return msgInfo, err
	} else if message.Signature != "" {
		msgSignitureBytes, err := b64.StdEncoding.DecodeString(message.Signature)
		if err != nil {
			zap.L().Error("Error while decoding signature", zap.Error(err))
			return pkg.MessageInfo{}, err
		}
		res, err := protoUtil.VerifyMessageSignature(msgSignitureBytes, decMsgInfoBytes, senderVerfier.PublicKey, cfg.Security.DSAScheme)
		if err != nil {
			zap.L().Error("Error while verifying signature", zap.Error(err))
			return pkg.MessageInfo{}, err
		}
		if !res {
			zap.L().Error("Signature is not valid")
			return pkg.MessageInfo{}, errors.New("Signature is not valid")
		}
		return msgInfo, nil
	} else {
		zap.L().Error("HMAC and Signature is missing")
		return pkg.MessageInfo{}, errors.New("HMAC and Signature is missing")
	}

}
