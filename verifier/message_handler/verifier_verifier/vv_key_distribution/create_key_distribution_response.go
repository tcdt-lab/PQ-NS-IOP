package vv_key_distribution

import (
	"database/sql"
	b64 "encoding/base64"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/verifier_verifier"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func CreateKeyDistributionResponse(db *sql.DB, c *config.Config, cipherText string, requestId int64) ([]byte, error) {
	msg := pkg.Message{}
	msgInfo := pkg.MessageInfo{}
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	cacheHandlerDa := data_access.NewCacheHandlerDA()
	userDa := data_access.GenerateVerifierUserDA(db)
	adminId, err := cacheHandlerDa.GetUserAdminId()
	if err != nil {
		return nil, err
	}
	user, err := userDa.GetVerifierUser(adminId)
	if err != nil {
		return nil, err
	}
	response := verifier_verifier.VVerifierKeyDistributionResponse{}
	response.CipherText = cipherText
	response.PublicKeyKem = user.PublicKeyKem
	msgInfo.OperationTypeId = pkg.VERIFIER_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID
	msgInfo.Nonce = "123"
	msgInfo.Params = response
	msgInfo.RequestId = requestId
	msgInfoBytes, err := protocolUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err
	}
	msg.PublicKeySig = user.PublicKeySig
	msg.IsEncrypted = false
	protocolUtil.SignMessageInfo(&msg, msgInfoBytes, user.SecretKeySig, c.Security.DSAScheme)
	msg.MsgInfo = b64.StdEncoding.EncodeToString(msgInfoBytes)
	msg.Hmac = ""
	if err != nil {
		return nil, err
	}
	msgBytes, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}
