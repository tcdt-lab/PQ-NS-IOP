package message_parser

import (
	b64 "encoding/base64"
	"errors"
	"test.org/protocol/pkg"

	"verifier/message_parser/gateway_verifier"
	"verifier/message_parser/util"

	"verifier/config"
	"verifier/data"
)

type MessageParser struct {
}

// Get a Message, parse it and then and generate an appropriate response
func (mp *MessageParser) ParseMessage(msg []byte, senderIp string, senderPort string, c config.Config) ([]byte, error) {
	var msgUtil pkg.MessageUtil
	message, err := msgUtil.ConvertByteToMessage(msg)
	var messageData pkg.MessageData
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
			return nil, err
		}
		if messageData.Signature != "" {
			res, err := msgUtil.VerifyMessageDataSignature(messageData, pubKeySig, c.Security.MlDSAScheme)
			if err != nil {
				return nil, err
			}
			if !res {
				return nil, errors.New("Signature verification failed")
			}
		} else if messageData.Hmac != "" {
			res, err := msgUtil.VerifyHmac(messageData, symmetricKey)
			if err != nil {
				return nil, err
			}
			if !res {
				return nil, errors.New("Hmac verification failed")
			}
		}

	} else { //if message is not encrypted it is in an init message, so keys are in the message body. In the handler we will check signature/hmac
		msgDataBytes, err := b64.StdEncoding.DecodeString(message.Data)
		if err != nil {
			return nil, err
		}
		messageData, err = msgUtil.ConvertByteToMessageData(msgDataBytes)
		if err != nil {
			return nil, err
		}
	}
	return mp.generateResponse(messageData, senderIp, senderPort, c)
}

func (mp *MessageParser) decryptMessage(msgData string, symmetricKey string, c config.Config) (pkg.MessageData, error) {

	msgUtil := util.MessageUtilGenerator(c.Security.CryptographyScheme)
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
