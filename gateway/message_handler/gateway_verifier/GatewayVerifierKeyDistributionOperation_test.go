package gateway_verifier

import (
	"fmt"
	"gateway/config"
	"gateway/message_handler/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateGatewayVerifierKeyDistributionMessage(t *testing.T) {

	c, err := config.ReadYaml()
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	var res []byte
	assert.NoError(t, err, "Error while reading config file")
	t.Run("TestCreateGatewayVerifierKeyDistributionMessage", func(t *testing.T) {

		res = CreateGatewayVerifierKeyDistributionMessage(c)
		assert.NotNil(t, res, "Error while generating response")
	})

	t.Run("Decode Response", func(t *testing.T) {
		msg, err := protocolUtil.ConvertByteToMessage(res)
		fmt.Println(err)
		assert.NoError(t, err, "Error while converting byte to message")
		assert.NotNil(t, msg, "Error while converting byte to message")
		msgData, err := protocolUtil.ConvertB64ToMessageData(msg.Data)
		assert.NoError(t, err, "Error while converting b64 string to message data")
		t.Log(msgData)
		t.Log(msgData.MsgInfo.Params)
		t.Log(msgData.MsgInfo)
		t.Log(msgData)
	})
}
