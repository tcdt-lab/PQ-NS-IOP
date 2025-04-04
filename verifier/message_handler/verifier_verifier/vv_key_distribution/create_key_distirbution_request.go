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

func CreateKeyDistributionRequest(c *config.Config, requestId int64, db *sql.DB) ([]byte, error) {

	msg := pkg.Message{}
	msgInfo := pkg.MessageInfo{}
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	cacheHandlerDa := data_access.GenerateCacheHandlerDA()
	verifierUerDa := data_access.GenerateVerifierUserDA(db)

	adminId, err := cacheHandlerDa.GetUserAdminId()
	if err != nil {
		return nil, err
	}
	verifierUer, err := verifierUerDa.GetVerifierUser(adminId)
	if err != nil {
		return nil, err
	}

	vvKeyDistributionReq := verifier_verifier.VVerifierKeyDistributionRequest{}
	vvKeyDistributionReq.RequestId = requestId
	vvKeyDistributionReq.VerifierIP = c.Server.Ip
	vvKeyDistributionReq.VerifierPort = c.Server.Port
	vvKeyDistributionReq.VerifierSignatureScheme = c.Security.DSAScheme
	vvKeyDistributionReq.VerifierPublicKeySignature = verifierUer.PublicKeySig
	vvKeyDistributionReq.VerifierPublicKeyKem = verifierUer.PublicKeyKem
	vvKeyDistributionReq.VerifierKemScheme = c.Security.KEMScheme

	msgInfo.OperationTypeId = pkg.VERIFIER_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID
	msgInfo.Nonce = "123"
	msgInfo.Params = vvKeyDistributionReq
	msgInfo.RequestId = requestId
	msgInfoBytes, err := protocolUtil.ConvertMessageInfoToByte(msgInfo)

	msg.PublicKeySig = verifierUer.PublicKeySig
	msg.IsEncrypted = false
	msg.Signature = ""
	msg.MsgInfo = b64.StdEncoding.EncodeToString(msgInfoBytes)
	msgByte, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err
	}

	return msgByte, nil
}
