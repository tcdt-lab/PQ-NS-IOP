package vv_sync_info

import (
	"database/sql"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/verifier_verifier"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func CreateSyncInfoResponse(result bool, destinationIp string, destinationPort string, db *sql.DB, c config.Config) ([]byte, error) {
	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	msgInfo := pkg.MessageInfo{}
	msg := pkg.Message{}
	vDA := data_access.GenerateVerifierDA(db)
	var resp verifier_verifier.VVSyncInfoOperationResponse
	resp.Result = result
	msgInfo.OperationTypeId = pkg.VERIFIER_VERIFIER_SYNC_INFO_OPERATION_RESPONSE_ID
	msgInfo.Params = resp
	msgInfoByte, err := msgUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err
	}
	destinationVerifier, err := vDA.GetVerifierByIpAndPort(destinationIp, destinationPort)
	if err != nil {
		return nil, err
	}
	key := destinationVerifier.SymmetricKey
	msg.Signature = ""
	hmacStr, _, _ := msgUtil.GenerateHmacMsgInfo(msgInfoByte, key)
	msgInfoStrEnc, _, err := msgUtil.EncryptMessageInfo(msgInfoByte, key)
	msg.MsgInfo = msgInfoStrEnc
	msg.Hmac = hmacStr
	msg.IsEncrypted = true
	msgBytes, err := msgUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}
