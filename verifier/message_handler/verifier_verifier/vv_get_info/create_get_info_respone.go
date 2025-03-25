package vv_get_info

import (
	"database/sql"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/verifier_verifier"
	"verifier/config"
	"verifier/data"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func CreateGetInfoResponse(c *config.Config, requestId int64, db *sql.DB, senderIp string, senderPort string) ([]byte, error) {
	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	var msgInfo = pkg.MessageInfo{}
	var msg = pkg.Message{}

	var resp verifier_verifier.VVInitInfoOperationResponse
	resp.RequestId = requestId
	veriferDataAccess := data_access.GenerateVerifierDA(db)
	desitantionVerifier, err := veriferDataAccess.GetVerifierByIpAndPort(senderIp, senderPort)
	gatewayDataAccess := data_access.GenerateGatewayDA(db)
	cachHandler := data_access.NewCacheHandlerDA()

	adminId, _ := cachHandler.GetUserAdminId()
	adminVerifier, err := veriferDataAccess.GetVerifier(adminId)
	resp.CurrentVerifierInfo = formatVerifier(adminVerifier)

	gateways, err := gatewayDataAccess.GetGateways()
	if err != nil {
		return nil, err
	}
	resp.GatewaysList = formatGateways(gateways)

	verifiers, err := veriferDataAccess.GetVerifiers()
	if err != nil {
		return nil, err
	}
	resp.VerifiersList = formatVerifiers(verifiers)

	msgInfo.RequestId = requestId
	msgInfo.OperationTypeId = pkg.VERIFIER_VERIFIER_GET_INFO_OPERATION_RESPONSE_ID
	msgInfo.Params = resp
	msgInfoByte, err := msgUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err
	}
	msg.Signature = ""
	hmacStr, _, _ := msgUtil.GenerateHmacMsgInfo(msgInfoByte, desitantionVerifier.SymmetricKey)
	msgInfoStrEnc, _, err := msgUtil.EncryptMessageInfo(msgInfoByte, desitantionVerifier.SymmetricKey)
	msg.MsgInfo = msgInfoStrEnc
	msg.Hmac = hmacStr
	msg.IsEncrypted = true
	msgBytes, err := msgUtil.ConvertMessageToByte(msg)
	return msgBytes, nil
}

func formatVerifier(verifiers data.Verifier) verifier_verifier.VVInitInfoStructureVerifier {
	var inintInfoVerifier = verifier_verifier.VVInitInfoStructureVerifier{}
	inintInfoVerifier.VerifierPort = verifiers.Port
	inintInfoVerifier.VerifierIpAddress = verifiers.Ip
	inintInfoVerifier.VerifierPublicKeySignature = verifiers.PublicKeySig
	inintInfoVerifier.SigScheme = verifiers.SigScheme
	inintInfoVerifier.VerifierPublicKeyKem = verifiers.PublicKeyKem
	inintInfoVerifier.TrustScore = verifiers.TrustScore
	inintInfoVerifier.IsInCommittee = verifiers.IsInCommittee
	return inintInfoVerifier

}
func formatVerifiers(verifiers []data.Verifier) []verifier_verifier.VVInitInfoStructureVerifier {
	var inintInfoVerifiers []verifier_verifier.VVInitInfoStructureVerifier
	for _, verifier := range verifiers {
		inintInfoVerifiers = append(inintInfoVerifiers, formatVerifier(verifier))

	}
	return inintInfoVerifiers
}

func formatGateways(gateways []data.Gateway) []verifier_verifier.VVInitInfoStructureGateway {
	var inintInfoGateways []verifier_verifier.VVInitInfoStructureGateway
	for _, gateway := range gateways {
		var inintInfoGateway = verifier_verifier.VVInitInfoStructureGateway{}
		inintInfoGateway.GatewayIpAddress = gateway.Ip
		inintInfoGateway.GatewayPort = gateway.Port
		inintInfoGateway.GatewayPublicKeySignature = gateway.PublicKeySig
		inintInfoGateway.SigScheme = gateway.SigScheme
		inintInfoGateway.GatewayPublicKeyKem = gateway.PublicKeyKem
		inintInfoGateway.KemScheme = gateway.KemScheme
		inintInfoGateways = append(inintInfoGateways, inintInfoGateway)
	}
	return inintInfoGateways
}
