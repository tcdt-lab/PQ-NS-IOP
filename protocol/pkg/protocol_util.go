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
	PBKDF2Handler     symmetric.PBKDF2
}

func (mp *ProtocolUtil) RegisterInterfacesInGob() {
	gob.Register(ErrorParams{})
	gob.Register(gateway_verifier.GatewayVerifierKeyDistributionRequest{})
	gob.Register(gateway_verifier.GatewayVerifierKeyDistributionResponse{})
	gob.Register(gateway_verifier.GatewayVerifierInitInfoStructureVerifier{})
	gob.Register(gateway_verifier.GatewayVerifierInitInfoStructureGateway{})
	gob.Register(gateway_verifier.GatewayVerifierInitInfoOperationResponse{})
	gob.Register(gateway_verifier.GatewayVerifierInitInfoOperationRequest{})
}

// VerifyMessageDataSignature Gets A message and verifies the signature of the message
func (mp *ProtocolUtil) VerifyMessageSignature(signatureBytes []byte, infoBytes []byte, pubKeyStr string, schemeName string) (bool, error) {
	response, err := mp.AsymmetricHandler.Verify(pubKeyStr, infoBytes, signatureBytes, schemeName)
	if err != nil {
		return false, err
	}
	return response, nil
}

// SignMessageInfo Gets A message and signs the message
// It returns A message with signature field filled
func (mp *ProtocolUtil) SignMessageInfo(msg *Message, messageInfoBytes []byte, secKeyStr string, schemeName string) error {

	response, err := mp.AsymmetricHandler.Sign(secKeyStr, messageInfoBytes, schemeName)
	if err != nil {
		zap.L().Error("ErrorParams while signing the message", zap.Error(err))
		return err
	}
	msg.Signature = b64.StdEncoding.EncodeToString(response)
	return nil
}

// DecryptMessageData Gets A message as byte string and decrypts the message
func (mp *ProtocolUtil) DecryptMessageInfo(msgInfo string, symmetricKey string) (MessageInfo, []byte, error) {

	msgBytes, err := b64.StdEncoding.DecodeString(msgInfo)
	if err != nil {
		return MessageInfo{}, nil, err
	}
	keyByte, err := mp.AesHandler.ConvertKeyStr64ToBytes(symmetricKey)
	if err != nil {
		return MessageInfo{}, nil, err
	}
	decInfoBytes, err := mp.AesHandler.Decrypt(msgBytes, keyByte)
	if err != nil {
		return MessageInfo{}, nil, err
	}
	convertedMessage, err := mp.ConvertByteToMessageInfo(decInfoBytes)
	if err != nil {
		return MessageInfo{}, nil, err
	}
	return convertedMessage, decInfoBytes, nil
}
func (mp *ProtocolUtil) ConvertPlainStrToMessageInfo(data string) (MessageInfo, []byte, error) {
	msgBytes, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return MessageInfo{}, nil, err
	}
	convertedMessage, err := mp.ConvertByteToMessageInfo(msgBytes)
	if err != nil {
		return MessageInfo{}, nil, err
	}
	return convertedMessage, msgBytes, nil
}

// It gets A message and encrypts the message and retun it as A byte array
func (mp *ProtocolUtil) EncryptMessageInfo(msgInfoBytes []byte, symmetricKey string) (string, []byte, error) {

	keyByte, err := mp.AesHandler.ConvertKeyStr64ToBytes(symmetricKey)
	if err != nil {
		return "", nil, err
	}
	encBytes, err := mp.AesHandler.Encrypt(msgInfoBytes, keyByte)
	if err != nil {
		return "", nil, err
	}
	encText := b64.StdEncoding.EncodeToString(encBytes)
	return encText, encBytes, nil
}

// It gets A message data produces HMAC for the message
// It returns A message with filled hmac
func (mp *ProtocolUtil) GenerateHmacMsgInfo(msgInfo []byte, key string) (string, []byte, error) {

	byteKey, err := mp.AesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return "", nil, err
	}

	hmacBytes, err := mp.HmacHandler.GenerateMessageMac(byteKey, msgInfo)
	if err != nil {
		return "", nil, err
	}
	hmacStr := mp.HmacHandler.ConvertHMacMsgToBase64(hmacBytes)
	return hmacStr, hmacBytes, nil
}

// It gets A message data and HMAC and verifies the HMAC

func (mp *ProtocolUtil) VerifyHmacByte(recievedHmac []byte, recievedMsgInfo []byte, key string) (bool, error) {

	byteKey, err := mp.AesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return false, err
	}
	_, generared_hmc, err := mp.GenerateHmacMsgInfo(recievedMsgInfo, key)
	if err != nil {
		return false, err
	}
	return mp.HmacHandler.VerifyMessageMac(byteKey, generared_hmc, recievedHmac), nil
}

func (mp *ProtocolUtil) VerifyHmac(msgHmac string, msgInfoBytes []byte, key string) (bool, error) {
	byteKey, err := mp.AesHandler.ConvertKeyStr64ToBytes(key)
	if err != nil {
		return false, err
	}
	hmacBytes, err := b64.StdEncoding.DecodeString(msgHmac)
	if err != nil {
		return false, err
	}
	if err != nil {
		return false, err
	}
	return mp.HmacHandler.VerifyMessageMac(byteKey, msgInfoBytes, hmacBytes), nil
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
