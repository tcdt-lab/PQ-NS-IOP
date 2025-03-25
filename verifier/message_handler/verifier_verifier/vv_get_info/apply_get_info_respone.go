package vv_get_info

import (
	"database/sql"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/verifier_verifier"
	"verifier/config"
	"verifier/data"
	"verifier/data_access"
)

func ApplyGetInfoResponse(msgInfo pkg.MessageInfo, db *sql.DB, c *config.Config) error {
	veriferDataAccess := data_access.GenerateVerifierDA(db)
	gatewayDataAccess := data_access.GenerateGatewayDA(db)
	operationResponse := msgInfo.Params.(verifier_verifier.VVInitInfoOperationResponse)
	senderVerifier := formatToDataVerifier(operationResponse.CurrentVerifierInfo)
	verifiers := formatToDataVerifiers(operationResponse.VerifiersList)
	gateways := formatToDataGateways(operationResponse.GatewaysList)

	_, err := gatewayDataAccess.AddUpdateGateways(gateways)
	if err != nil {
		return err
	}
	_, err = veriferDataAccess.AddUpdateVerifier(senderVerifier)
	if err != nil {
		return err
	}
	_, err = veriferDataAccess.AddUpdateVerifiers(verifiers)
	if err != nil {
		return err
	}
	return nil
}

func formatToDataGateways(vvGateways []verifier_verifier.VVInitInfoStructureGateway) []data.Gateway {
	var gateways []data.Gateway
	for _, vvGateway := range vvGateways {
		var gateway = data.Gateway{}
		gateway.Ip = vvGateway.GatewayIpAddress
		gateway.Port = vvGateway.GatewayPort
		gateway.PublicKeySig = vvGateway.GatewayPublicKeySignature
		gateway.SigScheme = vvGateway.SigScheme
		gateway.PublicKeyKem = vvGateway.GatewayPublicKeyKem
		gateway.KemScheme = vvGateway.KemScheme
		gateways = append(gateways, gateway)
	}

	return gateways
}

func formatToDataVerifier(vvVerifier verifier_verifier.VVInitInfoStructureVerifier) data.Verifier {
	var verifier = data.Verifier{}
	verifier.Port = vvVerifier.VerifierPort
	verifier.Ip = vvVerifier.VerifierIpAddress
	verifier.PublicKeySig = vvVerifier.VerifierPublicKeySignature
	verifier.SigScheme = vvVerifier.SigScheme
	verifier.PublicKeyKem = vvVerifier.VerifierPublicKeyKem
	verifier.TrustScore = vvVerifier.TrustScore
	verifier.IsInCommittee = vvVerifier.IsInCommittee
	return verifier
}

func formatToDataVerifiers(vvVerifiers []verifier_verifier.VVInitInfoStructureVerifier) []data.Verifier {
	var verifiers []data.Verifier
	for _, vvVerifier := range vvVerifiers {
		var verifier = formatToDataVerifier(vvVerifier)
		verifiers = append(verifiers, verifier)
	}
	return verifiers
}
