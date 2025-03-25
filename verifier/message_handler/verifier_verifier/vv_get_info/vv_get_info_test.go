package vv_get_info

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"verifier/config"
	"verifier/message_handler"
	"verifier/utility"
)

func Test_vv_get_info_process(t *testing.T) {
	c, err := config.ReadYaml()
	assert.NoError(t, err, "Error in ReadYaml")
	db, err := utility.GetDBConnection(*c)
	assert.NoError(t, err, "Error in GetDBConnection")
	assert.NotNilf(t, db, "Error in GetDBConnection")
	var reqId int64 = 12345
	var reqBytes []byte
	var resBytes []byte

	t.Run("Create Get Info Request", func(t *testing.T) {
		reqBytes, err = CreateGetInfoRequest(c, reqId, db)
		assert.NoError(t, err, "Error in CreateGetInfoRequest")
		assert.NotNilf(t, reqBytes, "Error in CreateGetInfoRequest")
	})
	t.Run("Create Get Info Response", func(t *testing.T) {
		resBytes, err = CreateGetInfoResponse(c, reqId, db, "127.0.0.1", "50051")
		assert.NoError(t, err, "Error in CreateGetInfoResponse")
		assert.NotNilf(t, resBytes, "Error in CreateGetInfoResponse")
	})
	t.Run("Apply Get Info Response", func(t *testing.T) {
		msgInfo, _, err := message_handler.ParseRequest(resBytes, "", "", db)
		assert.NoError(t, err, "Error in ParseRequest")
		err = ApplyGetInfoResponse(msgInfo, db, c)
		assert.NoError(t, err, "Error in ApplyGetInfoResponse")
	})

}
