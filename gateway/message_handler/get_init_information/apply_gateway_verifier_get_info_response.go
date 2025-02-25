package message_applier

import (
	"database/sql"
	"gateway/data"
	"gateway/data_access"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
)

func ApplyGatewayVerifierGetInfoResponse(msgData pkg.MessageData, db *sql.DB) error {
	vDA := data_access.GenerateVerifierDA(db)
	gDA := data_access.GenerateGatewayDA(db)
	gvGetInfoRes := msgData.MsgInfo.Params.(gateway_verifier.GatewayVerifierInitInfoOperationResponse)
	gatewaysList := extractGatewaysInfo(gvGetInfoRes)
	verifiersList := extractVerifiersInfo(gvGetInfoRes)
	err := gDA.AddUpdateGateways(gatewaysList)
	if err != nil {
		return err
	}
	err = vDA.AddUpdateVerifiers(verifiersList)
	if err != nil {
		return err
	}
	return nil
}

func extractGatewaysInfo(gvGetInfoRes gateway_verifier.GatewayVerifierInitInfoOperationResponse) []data.Gateway {
	var gatewaysInfo []data.Gateway
	for _, gatewayInfo := range gvGetInfoRes.GatewaysList {
		gateway := data.Gateway{}
		gateway.Ip = gatewayInfo.GatewayIpAddress
		gateway.Port = gatewayInfo.GatewayPort
		gateway.PublicKey = gatewayInfo.GatewayPublicKeyKem
		gatewaysInfo = append(gatewaysInfo, gateway)
	}
	return gatewaysInfo
}

func extractVerifiersInfo(gvGetInfoRes gateway_verifier.GatewayVerifierInitInfoOperationResponse) []data.Verifier {
	var verifiersInfo []data.Verifier
	for _, verifierInfo := range gvGetInfoRes.VerifiersList {
		verifier := data.Verifier{}
		verifier.Ip = verifierInfo.VerifierIpAddress
		verifier.Port = verifierInfo.VerifierPort
		verifier.PublicKey = verifierInfo.VerifierPublicKeyKem
		verifiersInfo = append(verifiersInfo, verifier)
	}
	return verifiersInfo
}
