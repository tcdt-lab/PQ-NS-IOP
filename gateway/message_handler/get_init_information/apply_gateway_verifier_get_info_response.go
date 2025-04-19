package get_init_information

import (
	"database/sql"

	"gateway/config"
	"gateway/data"
	"gateway/data_access"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
)

func ApplyGatewayVerifierGetInfoResponse(msgInfo pkg.MessageInfo, db *sql.DB) error {
	vDA := data_access.GenerateVerifierDA(db)
	gDA := data_access.GenerateGatewayDA(db)
	config, err := config.ReadYaml()
	bootstrapVerifier, err := vDA.GetVerifierByIpAndPort(config.BootstrapNode.Ip, config.BootstrapNode.Port)

	gvGetInfoRes := msgInfo.Params.(gateway_verifier.GatewayVerifierInitInfoOperationResponse)
	gatewaysList := extractGatewaysInfo(gvGetInfoRes)
	verifiersList := extractVerifiersInfo(gvGetInfoRes, bootstrapVerifier.SymmetricKey)
	senderVerifier := extractSenderVerifier(gvGetInfoRes, bootstrapVerifier.SymmetricKey)
	err = gDA.AddUpdateGateways(gatewaysList)
	if err != nil {
		return err
	}
	err = vDA.AddUpdateVerifiers(verifiersList)
	if err != nil {
		return err
	}
	err = vDA.AddUpdateVerifier(senderVerifier)
	return nil
}

func extractSenderVerifier(res gateway_verifier.GatewayVerifierInitInfoOperationResponse, symmetricKey string) data.Verifier {
	verifier := data.Verifier{
		Ip:           res.CurrentVerifierInfo.VerifierIpAddress,
		Port:         res.CurrentVerifierInfo.VerifierPort,
		PublicKey:    res.CurrentVerifierInfo.VerifierPublicKeySignature,
		SymmetricKey: symmetricKey,
	}
	return verifier
}

func extractGatewaysInfo(gvGetInfoRes gateway_verifier.GatewayVerifierInitInfoOperationResponse) []data.Gateway {
	var gatewaysInfo []data.Gateway
	for _, gatewayInfo := range gvGetInfoRes.GatewaysList {
		gateway := data.Gateway{}
		gateway.Ip = gatewayInfo.GatewayIpAddress
		gateway.Port = gatewayInfo.GatewayPort
		gateway.PublicKey = gatewayInfo.GatewayPublicKeySignature
		gatewaysInfo = append(gatewaysInfo, gateway)
	}
	return gatewaysInfo
}

func extractVerifiersInfo(gvGetInfoRes gateway_verifier.GatewayVerifierInitInfoOperationResponse, bootStrapVerifierSymmetricKey string) []data.Verifier {
	var verifiersInfo []data.Verifier
	for _, verifierInfo := range gvGetInfoRes.VerifiersList {
		verifier := data.Verifier{}
		verifier.SymmetricKey = bootStrapVerifierSymmetricKey
		verifier.Ip = verifierInfo.VerifierIpAddress
		verifier.Port = verifierInfo.VerifierPort
		verifier.PublicKey = verifierInfo.VerifierPublicKeySignature
		verifiersInfo = append(verifiersInfo, verifier)
	}
	return verifiersInfo
}
