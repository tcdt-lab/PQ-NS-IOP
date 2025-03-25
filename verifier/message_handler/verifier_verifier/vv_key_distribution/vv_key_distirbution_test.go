package vv_key_distribution

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"verifier/config"
	"verifier/message_handler"
	"verifier/utility"
)

func Test_KeyDistribution(t *testing.T) {
	c, err := config.ReadYaml()
	assert.NoError(t, err, "Error in ReadYaml")
	db, err := utility.GetDBConnection(*c)
	assert.NoError(t, err, "Error in GetDBConnection")
	assert.NotNilf(t, db, "Error in GetDBConnection")
	var reqId int64 = 12345
	var reqBytes []byte
	var resBytes []byte
	var cipherTextStr string

	t.Run("Test Create KeyDistribution Request", func(t *testing.T) {
		reqBytes, err = CreateKeyDistributionRequest(c, reqId, db)
		assert.NoError(t, err, "Error in CreateKeyDistributionRequest")
		assert.NotNilf(t, reqBytes, "Error in CreateKeyDistributionRequest")
	})

	t.Run("Test Apply KeyDistribution Request", func(t *testing.T) {
		msgInfo, _, err := message_handler.ParseRequest(reqBytes, "", "", db)
		assert.NoError(t, err, "Error in ParseRequest")
		cipherTextStr, err = ApplyKeyDistributionRequest(db, c, msgInfo)
		assert.NoError(t, err, "Error in ApplyKeyDistributionRequest")
		assert.NotNilf(t, cipherTextStr, "Error in ApplyKeyDistributionRequest")
	})

	t.Run("Test Create Key Distribution Response", func(t *testing.T) {
		resBytes, err = CreateKeyDistributionResponse(db, c, cipherTextStr, reqId)
		assert.NoError(t, err, "Error in CreateKeyDistributionResponse")
		assert.NotNilf(t, resBytes, "Error in CreateKeyDistributionResponse")
	})

	t.Run("Test Apply Key Distribution Response", func(t *testing.T) {
		msgInfo, _, err := message_handler.ParseRequest(resBytes, "", "", db)
		assert.NoError(t, err, "Error in ParseRequest")
		err = ApplyKeyDistributionResponse(msgInfo, db, c)
		assert.NoError(t, err, "Error in ApplyKeyDistributionResponse")
	})

}
