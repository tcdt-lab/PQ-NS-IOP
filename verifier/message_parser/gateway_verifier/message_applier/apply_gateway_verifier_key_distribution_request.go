package message_applier

import (
	b64 "encoding/base64"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"verifier/config"
	"verifier/data"
	"verifier/data_access"
	"verifier/message_parser/util"
)

func ApplyGatewayVerifierKeyDistributionRequest(msgData pkg.MessageData) (string, error) {
	cfg, err := config.ReadYaml()
	if err != nil {
		return "", err
	}

	protoUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)

	gvKeyDistributionParams := msgData.MsgInfo.Params.(gateway_verifier.GatewayVerifierKeyDistributionRequest)
	cipherTextBytes, sharedKey, _ := protoUtil.AsymmetricHandler.KemGenerateSecretKey("", gvKeyDistributionParams.GatewayPublicKeyKem, "", cfg.Security.KEMScheme)
	cipherTextStr := b64.StdEncoding.EncodeToString(cipherTextBytes)

	gateway := data.Gateway{}
	gateway.PublicKeySig = gvKeyDistributionParams.GatewayPublicKeySignature
	gateway.PublicKeyKem = gvKeyDistributionParams.GatewayPublicKeyKem
	gateway.KemScheme = gvKeyDistributionParams.GatewayKemScheme
	gateway.SigScheme = gvKeyDistributionParams.GatewaySignatureScheme
	gateway.Ip = gvKeyDistributionParams.GatewayIP
	gateway.Port = gvKeyDistributionParams.GatewayPort
	gateway.SymmetricKey = protoUtil.AesHandler.ConvertKeyBytesToStr64(sharedKey)
	gateway.Ticket = ""

	gDA := data_access.GatewayDA{}
	_, err = gDA.AddGateway(gateway)
	if err != nil {
		return "", err
	}

	return cipherTextStr, nil
}
