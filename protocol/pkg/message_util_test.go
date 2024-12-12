package pkg

import (
	"test.org/cryptography/pkg/asymmetric"

	"github.com/stretchr/testify/assert"

	"test.org/cryptography/pkg/symmetric"
	"testing"
	"time"
)

func generateSymmetricKeyStr() string {
	keyDerivation := symmetric.PBKDF2{}
	salt := []byte("salt")
	iterations := 4096

	aes := symmetric.AesGcm{}
	key := keyDerivation.KeyDerivation([]byte("password"), salt, iterations)
	keyB64 := aes.ConvertKeyBytesToStr64(key)
	return keyB64
}

func generateTicketSymmetricKeyStr() string {
	keyDerivation := symmetric.PBKDF2{}
	salt := []byte("salt")
	iterations := 4096
	aes := symmetric.AesGcm{}
	key := keyDerivation.KeyDerivation([]byte("passwordForTicket"), salt, iterations)
	keyB64 := aes.ConvertKeyBytesToStr64(key)
	return keyB64
}

func messageInfoGenerator() MessageInfo {

	msgData := MessageInfo{OperationTypeId: 1, Params: "params", SourceId: "sourceId", DestinationId: "DestinationId", Nonce: 4441}
	return msgData
}

func messageDataGenerator() MessageData {
	msgData := MessageData{MsgInfo: messageInfoGenerator(), Signature: "", Hmac: ""}
	return msgData
}

func ticketGenerator() Ticket {
	ticket := Ticket{SharedKey: generateTicketSymmetricKeyStr(), Deadline: time.Now(), SourceIp: "sourceIp", DestinationIp: "destinationIp"}
	return ticket
}

func messageUtilGenerator() MessageUtil {
	var util MessageUtil
	util.AesHandler = symmetric.AesGcm{}
	util.AsymmetricHandler = asymmetric.NewAsymmetricHandler("PQ")
	util.HmacHandler = symmetric.HMAC{}
	return util
}

func TestMessageUtil_generate_Encrypt_Decrypt_Message(t *testing.T) {
	// Arrange
	util := messageUtilGenerator()
	secKeyStr, pubKeyStr, err := util.AsymmetricHandler.DSKeyGen("ML-DSA-65")
	assert.NoError(t, err, "Error in DSKeyGen")
	msg := messageDataGenerator()
	key := generateSymmetricKeyStr()
	ticket := ticketGenerator()
	finalMsg := Message{}
	util.RegisterInterfacesInGob()
	var encryptedMsgStr string
	t.Run("Test Encrypt Message", func(t *testing.T) {
		err := util.SignMessageInfo(&msg, secKeyStr, "ML-DSA-65")
		assert.NoError(t, err, "Error in SignMessageInfo")
		err = util.GenerateHmac(&msg, key)
		assert.NoError(t, err, "Error in GenerateHmac")
		t.Log("Message Signed and HMAC generated")
		t.Log(msg)
		encryptedMsgStr, err = util.EncryptMessageData(msg, key)
		assert.NoError(t, err, "Error in EncryptMessageData")
		t.Log("Message Encrypted")
		finalMsg.Data = encryptedMsgStr
		encTicket, err := util.EncryptTicket(ticket, key)
		assert.NoError(t, err, "Error in EncryptTicket")
		t.Log("Ticket Encrypted")
		finalMsg.MsgTicket = encTicket
		t.Log(finalMsg)

	})
	t.Run("Test Decrypt Message", func(t *testing.T) {
		decryptedMsg, err := util.DecryptMessageData(finalMsg.Data, key)
		assert.NoError(t, err, "Error in DecryptMessageData")
		t.Log("Message Decrypted")
		t.Log(decryptedMsg)
		decryptedTicket, err := util.DecryptTicket(finalMsg.MsgTicket, key)
		assert.NoError(t, err, "Error in DecryptTicket")
		t.Log("Ticket Decrypted")
		t.Log(decryptedTicket)

		result, err := util.VerifyMessageDataSignature(decryptedMsg, pubKeyStr, "ML-DSA-65")
		assert.NoError(t, err, "Error in VerifyMessageDataSignature")
		assert.True(t, result, "Signature verification failed")
		t.Log("Signature Verified")
		result, err = util.VerifyHmac(decryptedMsg, key)
		assert.NoError(t, err, "Error in VerifyHmac")
		assert.True(t, result, "HMAC verification failed")
		t.Log("HMAC Verified")
		t.Log(decryptedMsg.MsgInfo.Params)
	})

}
