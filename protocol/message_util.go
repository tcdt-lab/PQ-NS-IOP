package protocol

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/gob"
	"go.uber.org/zap"

	"log"
	"test.org/cryptography/pkg/asymmetric"
	"test.org/cryptography/pkg/symmetric"
)

type MessageUtil struct {
	asymmetricHandler asymmetric.AsymmetricHandler
	aesHandler        symmetric.AesGcm
	hmacHandler       symmetric.HMAC
}

// VerifyMessageDataSignature Gets a message and verifies the signature of the message
func (mp *MessageUtil) VerifyMessageDataSignature(msg MessageData, pubKeyStr string, schemeName string) (bool, error) {

	signatureBytes, err := b64.StdEncoding.DecodeString(msg.Signature)
	if err != nil {
		return false, err
	}
	infoBytes := mp.ConvertMessageInfoToByte(msg.MsgInfo)

	response, err := mp.asymmetricHandler.Verify(pubKeyStr, infoBytes, signatureBytes, schemeName)
	if err != nil {
		return false, err
	}
	return response, nil
}

// SignMessageInfo Gets a message and signs the message
// It returns a message with signature field filled
func (mp *MessageUtil) SignMessageInfo(msg *MessageData, secKeyStr string, schemeName string) error {
	dataMsgBytes := mp.ConvertMessageInfoToByte(msg.MsgInfo)
	response, err := mp.asymmetricHandler.Sign(secKeyStr, dataMsgBytes, schemeName)
	if err != nil {
		zap.L().Error("Error while signing the message", zap.Error(err))
		return err
	}
	msg.Signature = b64.StdEncoding.EncodeToString(response)
	return nil
}

// DecryptMessageData Gets a message as byte string and decrypts the message
func (mp *MessageUtil) DecryptMessageData(msg string, symmetricKey string) (MessageData, error) {

	msgBytes, err := b64.StdEncoding.DecodeString(msg)
	if err != nil {
		return MessageData{}, err
	}
	keyByte, err := mp.aesHandler.ConvertKeyStr64ToBytes(symmetricKey)
	if err != nil {
		return MessageData{}, err
	}
	decText, err := mp.aesHandler.Decrypt(msgBytes, keyByte)
	if err != nil {
		return MessageData{}, err
	}
	convertedMessage := mp.ConvertByteToMessageData(decText)
	return convertedMessage, nil
}

// It gets a message and encrypts the message and retun it as a byte array
func (mp *MessageUtil) EncryptMessageData(msg MessageData, symmetricKey string) (string, error) {

	keyByte, err := mp.aesHandler.ConvertKeyStr64ToBytes(symmetricKey)
	if err != nil {
		return "", err
	}
	msgBytes := mp.ConvertMessageDataToByte(msg)
	encBytes, err := mp.aesHandler.Encrypt(msgBytes, keyByte)
	if err != nil {
		return "", err
	}
	encText := b64.StdEncoding.EncodeToString(encBytes)
	return encText, nil
}

// It gets a message data produces HMAC for the message
// It returns a message with filled hmac
func (mp *MessageUtil) GenerateHmac(msg *MessageData, key string) error {

	byteKey, err := mp.aesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return err
	}

	hmacBytes, err := mp.hmacHandler.GenerateMessageMac(byteKey, mp.ConvertMessageInfoToByte(msg.MsgInfo))
	if err != nil {
		return err
	}
	msg.Hmac = mp.hmacHandler.ConvertHMacMsgToBase64(hmacBytes)
	return nil
}

// It gets a message data and HMAC and verifies the HMAC
func (mp *MessageUtil) VerifyHmac(msg MessageData, key string) (bool, error) {

	byteKey, err := mp.aesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return false, err
	}

	msgHmac := mp.hmacHandler.ConvertBase64ToHMacMsg(msg.Hmac)
	return mp.hmacHandler.VerifyMessageMac(byteKey, mp.ConvertMessageInfoToByte(msg.MsgInfo), msgHmac), nil
}

// It gets a message and returna ticket struct  Ticket
func (mp *MessageUtil) DecryptTicket(ticket string, key string) (Ticket, error) {

	byteKey, err := mp.aesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return Ticket{}, err
	}
	ticketBytes, err := b64.StdEncoding.DecodeString(ticket)
	if err != nil {
		return Ticket{}, err
	}

	decTicket, err := mp.aesHandler.Decrypt(ticketBytes, byteKey)
	if err != nil {
		return Ticket{}, err
	}
	return mp.ConvertByteToTicket(decTicket), nil
}

func (mp *MessageUtil) EncryptTicket(ticket Ticket, key string) (string, error) {

	byteKey, err := mp.aesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return "", err
	}
	ticketBytes := mp.ConvertTicketToByte(ticket)
	encTicket, err := mp.aesHandler.Encrypt(ticketBytes, byteKey)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(encTicket), nil
}

// It gets a ticket and encrypts the ticket. it retruns both byte  array and base64 encoded version of the byte array
func (mp *MessageUtil) ConvertTicketToByte(ticket Ticket) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(ticket)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return buf.Bytes()
}

// It gets a ticket and dencrypts the ticket. it retruns both byte  array and base64 encoded version of the byte array
func (mp *MessageUtil) ConvertByteToTicket(data []byte) Ticket {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var ticket Ticket
	err := dec.Decode(&ticket)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return ticket
}

func (mp *MessageUtil) ConvertMessageDataToByte(msg MessageData) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(msg)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return buf.Bytes()
}

func (mp *MessageUtil) ConvertByteToMessageData(data []byte) MessageData {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var msg MessageData
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return msg
}

func (mp *MessageUtil) ConvertMessageInfoToByte(data MessageInfo) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(data)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return buf.Bytes()
}

func (mp *MessageUtil) ConvertByteToMessageInfo(data []byte) MessageInfo {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var msg MessageInfo
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return msg
}

func (mp *MessageUtil) ConvertMessageToByte(data Message) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(data)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return buf.Bytes()
}

func (mp *MessageUtil) ConvertByteToMessage(data []byte) Message {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var msg Message
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return msg
}
