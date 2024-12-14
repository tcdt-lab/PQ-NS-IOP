package error

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"gateway/config"
	"gateway/data"
	"gateway/message_handler/util"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenerateUnencryptedGeneralErrorResponse(t *testing.T) {
	var res []byte
	c, err := config.ReadYaml()
	var protocolUtil = util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	t.Run("TestGenerateUnencryptedGeneralErrorResponse", func(t *testing.T) {

		assert.NoError(t, err, "ErrorParams while reading config file")
		testError := errors.New("Test Error")
		nonce := "testNonce"
		db, err := util.GetDBConnection(*c)
		assert.NoError(t, err, "ErrorParams while getting db connection")
		currentUser := getCurrentUser(db)
		res = GenerateUnencryptedGeneralErrorResponse(testError, *c, db, nonce, 1, currentUser)
		assert.NotNil(t, res, "ErrorParams while generating response")
	})
	t.Run("Decrypt the errorMessage", func(t *testing.T) {
		msg, err := protocolUtil.ConvertByteToMessage(res)
		assert.NoError(t, err, "ErrorParams while converting byte to message")
		assert.NotNil(t, msg, "ErrorParams while converting byte to message")

		msgData, err := protocolUtil.ConvertB64ToMessageData(msg.Data)
		assert.NoError(t, err, "ErrorParams while converting b64 string to message data")
		t.Log(msgData)
	})

}

func getCurrentUser(db *sql.DB) data.GatewayUser {
	cg, err := data.GetGatewayUserByPassword(db, b64.StdEncoding.EncodeToString([]byte(os.Getenv("PQ_NS_IOP_GU_PASS"))))
	if err != nil {
		panic(err)
	}
	return cg
}
