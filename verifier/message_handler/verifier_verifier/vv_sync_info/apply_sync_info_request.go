package vv_sync_info

import (
	"database/sql"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/verifier_verifier"
	"verifier/config"
	"verifier/data"
	"verifier/data_access"
)

func ApplySyncInfoRequest(msgInfo pkg.MessageInfo, db *sql.DB, c config.Config) error {

	verifierDa := data_access.GenerateVerifierDA(db)
	gatewayDa := data_access.GenerateGatewayDA(db)

	operationRequest := msgInfo.Params.(verifier_verifier.VVSyncInfoOperationRequest)
	dataVerifier := formatToDataVerifier(operationRequest.CurrentVerifierInfo)
	dataGateways := formatToDataGateways(operationRequest.GatewaysList)
	dataVerifiers := formatToDataVerifiers(operationRequest.VerifiersList)
	_, err := verifierDa.AddUpdateVerifier(dataVerifier)
	if err != nil {
		return err
	}
	_, err = verifierDa.AddUpdateVerifiers(dataVerifiers)
	if err != nil {
		return err
	}
	_, err = gatewayDa.AddUpdateGateways(dataGateways)
	if err != nil {
		return err
	}

	return nil

}

func formatToDataGateways(vvGateways []verifier_verifier.VVSyncInfoStructureGateway) []data.Gateway {
	var gateways []data.Gateway
	for _, vvGateway := range vvGateways {
		var gateway = data.Gateway{}
		gateway.Ip = vvGateway.GatewayIpAddress
		gateway.Port = vvGateway.GatewayPort
		gateway.PublicKeySig = vvGateway.GatewayPublicKeySignature
		gateway.SigScheme = vvGateway.SigScheme
		gateway.PublicKeyKem = vvGateway.GatewayPublicKeyKem
		gateway.KemScheme = vvGateway.KemScheme
		gateway.SymmetricKey = vvGateway.SymmetricKey
		gateways = append(gateways, gateway)
	}

	return gateways
}

func formatToDataVerifier(vvVerifier verifier_verifier.VVSyncInfoStructureVerifier) data.Verifier {
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

func formatToDataVerifiers(vvVerifiers []verifier_verifier.VVSyncInfoStructureVerifier) []data.Verifier {
	var verifiers []data.Verifier
	for _, vvVerifier := range vvVerifiers {
		var verifier = formatToDataVerifier(vvVerifier)
		verifiers = append(verifiers, verifier)
	}
	return verifiers
}
