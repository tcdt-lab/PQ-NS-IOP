package message_parser

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"go.uber.org/zap"
	"os"
	"test.org/protocol/pkg"

	"verifier/message_parser/gateway_verifier"
	"verifier/message_parser/util"

	"verifier/config"
	"verifier/data"
)

type MessageParser struct {
}

// Get a Message, parse it and then and generate an appropriate response
// all signiture checking and decryption should be done in parser
func (mp *MessageParser) ParseMessage(msg []byte, senderIp string, senderPort string, c config.Config) ([]byte, error) {
	var msgUtil pkg.ProtocolUtil
	message, err := msgUtil.ConvertByteToMessage(msg)
	var messageData pkg.MessageData
	db, err := util.GetDBConnection(c)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if message.IsEncrypted {
		symmetricKey, pubKeySig, err := mp.getSenderKeys(senderIp, c)
		if err != nil {
			return nil, err
		}
		messageData, err = mp.decryptMessage(message.Data, symmetricKey, c)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, c, messageData, db), err
		}
		if messageData.Signature != "" {
			res, err := msgUtil.VerifyMessageDataSignature(messageData, pubKeySig, c.Security.MlDSAScheme)
			if err != nil {
				return nil, err
			}
			if !res {
				err = errors.New("Signature verification failed")
				return mp.GenerateGeneralErrorResponse(err, c, messageData, db), err
			}
		} else if messageData.Hmac != "" {
			res, err := msgUtil.VerifyHmac(messageData, symmetricKey)
			if err != nil {
				return mp.GenerateGeneralErrorResponse(err, c, messageData, db), err
			}
			if !res {
				err := errors.New("Hmac verification failed")
				return mp.GenerateGeneralErrorResponse(err, c, messageData, db), err
			}
		}

	} else { //if message is not encrypted it is in an init message, so keys are in the message body. In the handler we will check signature/hmac
		msgDataBytes, err := b64.StdEncoding.DecodeString(message.Data)
		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, c, messageData, db), err
		}
		messageData, err = msgUtil.ConvertByteToMessageData(msgDataBytes)

		if err != nil {
			return mp.GenerateGeneralErrorResponse(err, c, messageData, db), err
		}
	}
	return mp.generateResponse(messageData, senderIp, senderPort, c)
}

func (mp *MessageParser) decryptMessage(msgData string, symmetricKey string, c config.Config) (pkg.MessageData, error) {

	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	decryptedMsg, err := msgUtil.DecryptMessageData(msgData, symmetricKey)
	if err != nil {
		return pkg.MessageData{}, err
	}
	return decryptedMsg, nil
}

func (mp *MessageParser) getSenderKeys(senderIp string, c config.Config) (string, string, error) {
	db, err := util.GetDBConnection(c)
	if err != nil {
		return "", "", err
	}
	var gateway data.Gateway
	var verifier data.Verifier

	gateway, err = data.GetGatewayByIp(db, senderIp)
	if err == nil {
		return gateway.SymmetricKey, gateway.PublicKeySig, nil
	}
	verifier, err = data.GetVerifierByIp(db, senderIp)
	if err != nil {
		return "", "", err
	}
	return verifier.SymmetricKey, verifier.PublicKeySig, nil
}
func (mp *MessageParser) generateResponse(msgData pkg.MessageData, senderIp string, senderPort string, c config.Config) ([]byte, error) {
	switch msgData.MsgInfo.OperationTypeId {
	case pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID:
		return gateway_verifier.GatewayVerifierKeyDistributionHandler(msgData, c)
	}
	return nil, errors.New("Operation type not found")
}

// go to util
func (mp *MessageParser) GenerateGeneralErrorResponse(err error, c config.Config, incomingMsg pkg.MessageData, db *sql.DB) []byte {
	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	msg := pkg.Message{}
	msgData := pkg.MessageData{}
	msgInfo := pkg.MessageInfo{}
	errorParams := pkg.ErrorParams{}

	verifeirUser, err := data.GetVerifierUserByPassword(db, os.Getenv("PQ_NS_IOP_VU_PASS"))
	if err != nil {
		zap.L().Error("Error while getting verifier_verifier user", zap.Error(err))
	}
	errorParams.ErrorCode = pkg.GENERAL_ERROR
	errorParams.ErrorMessage = err.Error()
	msgInfo.Nonce = incomingMsg.MsgInfo.Nonce
	msgInfo.Params = errorParams
	msgData.MsgInfo = msgInfo
	msgUtil.SignMessageInfo(&msgData, verifeirUser.SecretKeySig, c.Security.MlDSAScheme)
	msgDataByte, err := msgUtil.ConvertMessageDataToByte(msgData)
	if err != nil {
		zap.L().Error("Error while converting message data to byte", zap.Error(err))
	}
	msg.Data = b64.StdEncoding.EncodeToString(msgDataByte)
	msg.IsEncrypted = false
	msg.MsgTicket = ""
	msgByte, err := msgUtil.ConvertMessageToByte(msg)
	return msgByte
}