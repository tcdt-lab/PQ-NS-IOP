package pkg

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/gob"
	"go.uber.org/zap"
	"test.org/protocol/pkg/gateway_verifier"

	"test.org/cryptography/pkg/asymmetric"
	"test.org/cryptography/pkg/symmetric"
)

type ProtocolUtil struct {
	AsymmetricHandler asymmetric.AsymmetricHandler
	AesHandler        symmetric.AesGcm
	HmacHandler       symmetric.HMAC
}

func (mp *ProtocolUtil) RegisterInterfacesInGob() {
	gob.Register(ErrorParams{})
	gob.Register(gateway_verifier.GatewayVerifierKeyDistributionRequest{})
	gob.Register(gateway_verifier.GatewayVerifierKeyDistributionResponse{})
}

// VerifyMessageDataSignature Gets A message and verifies the signature of the message
func (mp *ProtocolUtil) VerifyMessageDataSignature(msg MessageData, pubKeyStr string, schemeName string) (bool, error) {

	signatureBytes, err := b64.StdEncoding.DecodeString(msg.Signature)
	if err != nil {
		return false, err
	}
	infoBytes, err := mp.ConvertMessageInfoToByte(msg.MsgInfo)
	if err != nil {
		return false, err
	}
	response, err := mp.AsymmetricHandler.Verify(pubKeyStr, infoBytes, signatureBytes, schemeName)
	if err != nil {
		return false, err
	}
	return response, nil
}

// SignMessageInfo Gets A message and signs the message
// It returns A message with signature field filled
func (mp *ProtocolUtil) SignMessageInfo(msg *MessageData, secKeyStr string, schemeName string) error {
	dataMsgBytes, err := mp.ConvertMessageInfoToByte(msg.MsgInfo)
	if err != nil {
		return err
	}
	response, err := mp.AsymmetricHandler.Sign(secKeyStr, dataMsgBytes, schemeName)
	if err != nil {
		zap.L().Error("ErrorParams while signing the message", zap.Error(err))
		return err
	}
	msg.Signature = b64.StdEncoding.EncodeToString(response)
	return nil
}

// DecryptMessageData Gets A message as byte string and decrypts the message
func (mp *ProtocolUtil) DecryptMessageData(msg string, symmetricKey string) (MessageData, error) {

	msgBytes, err := b64.StdEncoding.DecodeString(msg)
	if err != nil {
		return MessageData{}, err
	}
	keyByte, err := mp.AesHandler.ConvertKeyStr64ToBytes(symmetricKey)
	if err != nil {
		return MessageData{}, err
	}
	decText, err := mp.AesHandler.Decrypt(msgBytes, keyByte)
	if err != nil {
		return MessageData{}, err
	}
	convertedMessage, err := mp.ConvertByteToMessageData(decText)
	if err != nil {
		return MessageData{}, err
	}
	return convertedMessage, nil
}

// It gets A message and encrypts the message and retun it as A byte array
func (mp *ProtocolUtil) EncryptMessageData(msg MessageData, symmetricKey string) (string, error) {

	keyByte, err := mp.AesHandler.ConvertKeyStr64ToBytes(symmetricKey)
	if err != nil {
		return "", err
	}
	msgBytes, err := mp.ConvertMessageDataToByte(msg)
	if err != nil {
		return "", err
	}
	encBytes, err := mp.AesHandler.Encrypt(msgBytes, keyByte)
	if err != nil {
		return "", err
	}
	encText := b64.StdEncoding.EncodeToString(encBytes)
	return encText, nil
}

// It gets A message data produces HMAC for the message
// It returns A message with filled hmac
func (mp *ProtocolUtil) GenerateHmac(msg *MessageData, key string) error {

	byteKey, err := mp.AesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return err
	}
	byteMsgInfo, err := mp.ConvertMessageInfoToByte(msg.MsgInfo)
	if err != nil {
		return err
	}
	hmacBytes, err := mp.HmacHandler.GenerateMessageMac(byteKey, byteMsgInfo)
	if err != nil {
		return err
	}
	msg.Hmac = mp.HmacHandler.ConvertHMacMsgToBase64(hmacBytes)
	return nil
}

// It gets A message data and HMAC and verifies the HMAC
func (mp *ProtocolUtil) VerifyHmac(msg MessageData, key string) (bool, error) {

	byteKey, err := mp.AesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return false, err
	}

	msgHmac := mp.HmacHandler.ConvertBase64ToHMacMsg(msg.Hmac)
	byteMsgInfo, err := mp.ConvertMessageInfoToByte(msg.MsgInfo)
	if err != nil {
		return false, err
	}
	return mp.HmacHandler.VerifyMessageMac(byteKey, byteMsgInfo, msgHmac), nil
}

// It gets A message and returna ticket struct  Ticket
func (mp *ProtocolUtil) DecryptTicket(ticket string, key string) (Ticket, error) {

	byteKey, err := mp.AesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return Ticket{}, err
	}
	ticketBytes, err := b64.StdEncoding.DecodeString(ticket)
	if err != nil {
		return Ticket{}, err
	}

	decTicket, err := mp.AesHandler.Decrypt(ticketBytes, byteKey)
	if err != nil {
		return Ticket{}, err
	}
	convertedTicket, err := mp.ConvertByteToTicket(decTicket)
	if err != nil {
		return Ticket{}, err
	}
	return convertedTicket, nil
}

func (mp *ProtocolUtil) EncryptTicket(ticket Ticket, key string) (string, error) {

	byteKey, err := mp.AesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return "", err
	}
	ticketBytes, err := mp.ConvertTicketToByte(ticket)
	if err != nil {
		return "", err
	}
	encTicket, err := mp.AesHandler.Encrypt(ticketBytes, byteKey)
	if err != nil {
		return "", err
	}
	return b64.StdEncoding.EncodeToString(encTicket), nil
}

// It gets A ticket and encrypts the ticket. it retruns both byte  array and base64 encoded version of the byte array
func (mp *ProtocolUtil) ConvertTicketToByte(ticket Ticket) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(ticket)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// It gets A ticket and dencrypts the ticket. it retruns both byte  array and base64 encoded version of the byte array
func (mp *ProtocolUtil) ConvertByteToTicket(data []byte) (Ticket, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var ticket Ticket
	err := dec.Decode(&ticket)
	if err != nil {
		return Ticket{}, err
	}
	return ticket, nil
}

func (mp *ProtocolUtil) ConvertMessageDataToByte(msg MessageData) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(msg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func (mp *ProtocolUtil) ConvertByteToMessageData(data []byte) (MessageData, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var msg MessageData
	err := dec.Decode(&msg)
	if err != nil {
		return MessageData{}, err
	}
	return msg, nil
}

func (mp *ProtocolUtil) ConvertMessageInfoToByte(data MessageInfo) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func (mp *ProtocolUtil) ConvertByteToMessageInfo(data []byte) (MessageInfo, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var msg MessageInfo
	err := dec.Decode(&msg)
	if err != nil {
		return MessageInfo{}, err
	}
	return msg, nil
}

func (mp *ProtocolUtil) ConvertMessageToByte(data Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (mp *ProtocolUtil) ConvertByteToMessage(data []byte) (Message, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var msg Message
	err := dec.Decode(&msg)
	if err != nil {
		return Message{}, err
	}
	return msg, nil
}
