package vv_sync_info

import (
	"database/sql"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/verifier_verifier"
	"verifier/config"
	"verifier/data"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func CreateSyncInfoRequest(c *config.Config, requestId int64, destinationIp string, destinationPort string, db *sql.DB) ([]byte, error) {
	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	var msgInfo = pkg.MessageInfo{}
	var msg = pkg.Message{}

	var resp verifier_verifier.VVSyncInfoOperationRequest
	resp.RequestId = requestId
	veriferDataAccess := data_access.GenerateVerifierDA(db)
	desitantionVerifier, err := veriferDataAccess.GetVerifierByIpAndPort(destinationIp, destinationPort)
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
	msgInfo.OperationTypeId = pkg.VERIFIER_VERIFIER_SYNC_INFO_OPERATION_REQEST_ID
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

func formatVerifier(verifiers data.Verifier) verifier_verifier.VVSyncInfoStructureVerifier {
	var syncInfoVerifier = verifier_verifier.VVSyncInfoStructureVerifier{}
	syncInfoVerifier.VerifierPort = verifiers.Port
	syncInfoVerifier.VerifierIpAddress = verifiers.Ip
	syncInfoVerifier.VerifierPublicKeySignature = verifiers.PublicKeySig
	syncInfoVerifier.SigScheme = verifiers.SigScheme
	syncInfoVerifier.VerifierPublicKeyKem = verifiers.PublicKeyKem
	syncInfoVerifier.TrustScore = verifiers.TrustScore
	syncInfoVerifier.IsInCommittee = verifiers.IsInCommittee
	syncInfoVerifier.SymmetricKey = verifiers.SymmetricKey
	return syncInfoVerifier

}
func formatVerifiers(verifiers []data.Verifier) []verifier_verifier.VVSyncInfoStructureVerifier {
	var syncInfoVerifiers []verifier_verifier.VVSyncInfoStructureVerifier
	for _, verifier := range verifiers {
		syncInfoVerifiers = append(syncInfoVerifiers, formatVerifier(verifier))

	}
	return syncInfoVerifiers
}

func formatGateways(gateways []data.Gateway) []verifier_verifier.VVSyncInfoStructureGateway {
	var syncInfoGateways []verifier_verifier.VVSyncInfoStructureGateway
	for _, gateway := range gateways {
		var inintInfoGateway = verifier_verifier.VVSyncInfoStructureGateway{}
		inintInfoGateway.GatewayIpAddress = gateway.Ip
		inintInfoGateway.GatewayPort = gateway.Port
		inintInfoGateway.GatewayPublicKeySignature = gateway.PublicKeySig
		inintInfoGateway.SigScheme = gateway.SigScheme
		inintInfoGateway.GatewayPublicKeyKem = gateway.PublicKeyKem
		inintInfoGateway.KemScheme = gateway.KemScheme
		inintInfoGateway.SymmetricKey = gateway.SymmetricKey
		syncInfoGateways = append(syncInfoGateways, inintInfoGateway)
	}
	return syncInfoGateways
}
