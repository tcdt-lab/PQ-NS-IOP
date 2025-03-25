package vv_get_info

import (
	"database/sql"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/verifier_verifier"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func CreateGetInfoRequest(c *config.Config, requestId int64, db *sql.DB) ([]byte, error) {
	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	var req verifier_verifier.VVInitInfoOperationRequest
	req.RequestId = requestId
	veriferDataAccess := data_access.GenerateVerifierDA(db)
	boostrapVerifier, err := veriferDataAccess.GetVerifierByIpAndPort(c.BootstrapNode.Ip, c.BootstrapNode.Port)
	if err != nil {
		return nil, err
	}
	msg := pkg.Message{}
	msgInfo := pkg.MessageInfo{}
	msgInfo.OperationTypeId = pkg.VERIFIER_VERIFIER_GET_INFO_OPERATION_REQEST_ID
	msgInfo.RequestId = requestId
	msgInfo.Params = req
	msgInfoByte, err := msgUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err
	}
	msg.Signature = ""
	hmacStr, _, _ := msgUtil.GenerateHmacMsgInfo(msgInfoByte, boostrapVerifier.SymmetricKey)
	msgInfoStrEnc, _, err := msgUtil.EncryptMessageInfo(msgInfoByte, boostrapVerifier.SymmetricKey)
	msg.MsgInfo = msgInfoStrEnc
	msg.Hmac = hmacStr
	msg.IsEncrypted = true
	msgBytes, err := msgUtil.ConvertMessageToByte(msg)
	return msgBytes, nil
}
