package gv_get_info

import (
	"database/sql"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func CreateGateVerifierGetInfoResponse(pubKeySender string, reqId int64, db *sql.DB) ([]byte, error) {
	cfg, err := config.ReadYaml()
	if err != nil {
		return nil, err
	}
	msg := pkg.Message{}

	msgInfo := pkg.MessageInfo{}
	gDa := data_access.GenerateGatewayDA(db)
	protoUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	gvGetInfoResponse := gateway_verifier.GatewayVerifierInitInfoOperationResponse{}
	gvVerifier := []gateway_verifier.GatewayVerifierInitInfoStructureVerifier{}
	gvGateways := []gateway_verifier.GatewayVerifierInitInfoStructureGateway{}
	currentVerifierInfo := gateway_verifier.GatewayVerifierInitInfoStructureVerifier{}
	senderGt, err := gDa.GetGatewayByPublicKeySig(pubKeySender)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	gvGateways, err = fillGateways(gvGateways, db)
	if err != nil {
		return nil, err
	}
	gvGetInfoResponse.GatewaysList = gvGateways

	gvVerifier, err = fillVerifiers(gvVerifier, db)
	if err != nil {
		return nil, err
	}

	currentVerifierInfo, err = fillCurrentVerifierInfo(db, *cfg)
	if err != nil {
		return nil, err
	}
	gvGetInfoResponse.VerifiersList = gvVerifier
	gvGetInfoResponse.CurrentVerifierInfo = currentVerifierInfo
	gvGetInfoResponse.RequestId = reqId
	msgInfo.Params = gvGetInfoResponse
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_GET_INFO_OPERATION_RESPONSE
	msgInfo.Nonce = "123"
	msgInfoBytes, _ := protoUtil.ConvertMessageInfoToByte(msgInfo)

	msg.Signature = ""
	hmacStr, _, _ := protoUtil.GenerateHmacMsgInfo(msgInfoBytes, senderGt.SymmetricKey)
	msgInfoStrEnc, _, err := protoUtil.EncryptMessageInfo(msgInfoBytes, senderGt.SymmetricKey)
	if err != nil {
		return nil, err
	}
	msg.Hmac = hmacStr
	msg.MsgInfo = msgInfoStrEnc
	msg.IsEncrypted = true
	msg.MsgTicket = ""
	msgBytes, err := protoUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}

func fillCurrentVerifierInfo(db *sql.DB, c config.Config) (gateway_verifier.GatewayVerifierInitInfoStructureVerifier, error) {
	vuda := data_access.GenerateVerifierUserDA(db)
	cacheHandler := data_access.GenerateCacheHandlerDA()
	adminId, err := cacheHandler.GetUserAdminId()
	if err != nil {
		return gateway_verifier.GatewayVerifierInitInfoStructureVerifier{}, err
	}
	verifierUser, err := vuda.GetVerifierUser(adminId)
	if err != nil {
		return gateway_verifier.GatewayVerifierInitInfoStructureVerifier{}, err
	}
	currentVerifierInfo := gateway_verifier.GatewayVerifierInitInfoStructureVerifier{
		VerifierIpAddress:          c.Server.Ip,
		VerifierPort:               c.Server.Port,
		VerifierPublicKeySignature: verifierUser.PublicKeySig,
		SigScheme:                  c.Security.DSAScheme,
		VerifierPublicKeyKem:       verifierUser.PublicKeyKem,
		TrustScore:                 0,
		IsInCommittee:              false,
	}
	return currentVerifierInfo, nil
}

func fillGateways(gatewaysStruct []gateway_verifier.GatewayVerifierInitInfoStructureGateway, db *sql.DB) ([]gateway_verifier.GatewayVerifierInitInfoStructureGateway, error) {
	gDA := data_access.GenerateGatewayDA(db)

	allGateways, err := gDA.GetGateways()

	for _, gateway := range allGateways {
		gatewaysStruct = append(gatewaysStruct, gateway_verifier.GatewayVerifierInitInfoStructureGateway{
			GatewayIpAddress:          gateway.Ip,
			GatewayPort:               gateway.Port,
			GatewayPublicKeyKem:       gateway.PublicKeyKem,
			GatewayPublicKeySignature: gateway.PublicKeySig,
			KemScheme:                 gateway.KemScheme,
			SigScheme:                 gateway.SigScheme,
		})
	}
	return gatewaysStruct, err
}

func fillVerifiers(verifiersStruct []gateway_verifier.GatewayVerifierInitInfoStructureVerifier, db *sql.DB) ([]gateway_verifier.GatewayVerifierInitInfoStructureVerifier, error) {
	vDA := data_access.GenerateVerifierDA(db)
	allVerifiers, err := vDA.GetVerifiers()

	for _, verifier := range allVerifiers {
		verifiersStruct = append(verifiersStruct, gateway_verifier.GatewayVerifierInitInfoStructureVerifier{
			VerifierIpAddress:          verifier.Ip,
			VerifierPort:               verifier.Port,
			VerifierPublicKeySignature: verifier.PublicKeySig,
			SigScheme:                  verifier.SigScheme,
			IsInCommittee:              verifier.IsInCommittee,
			TrustScore:                 verifier.TrustScore,
		})
	}
	return verifiersStruct, err
}
