package key_distribution

import (
	"database/sql"
	b64 "encoding/base64"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"verifier/config"
	"verifier/data"
	"verifier/data/transactions/tx_gateway_verifier"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func ApplyGatewayVerifierKeyDistributionRequest(msgData pkg.MessageInfo, db *sql.DB) (string, error) {
	cfg, err := config.ReadYaml()
	if err != nil {
		return "", err
	}

	protoUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	gDA := data_access.GenerateGatewayDA(db)
	vuDA := data_access.GenerateVerifierUserDA(db)
	verifierUSer, err := vuDA.GetAdminVerifierUser()
	gvKeyDistributionParams := msgData.Params.(gateway_verifier.GatewayVerifierKeyDistributionRequest)
	cipherTextBytes, sharedKey, _ := protoUtil.AsymmetricHandler.KemGenerateSecretKey(verifierUSer.SecretKeyKem, gvKeyDistributionParams.GatewayPublicKeyKem, "", cfg.Security.KEMScheme)
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

	if err != nil {
		zap.L().Error("Error while getting verifier user", zap.Error(err))
		return "", err
	}
	//verifierUSer.SymmetricKey = gateway.SymmetricKey
	pubKeySigExist, _ := gDA.IfGatewayExistByPublicKeySig(gateway.PublicKeySig)
	if !pubKeySigExist {
		err = tx_gateway_verifier.SharedKeyAndGatewayRegistration(verifierUSer, gateway, db)
		if err != nil {
			zap.L().Error("Error while registering gateway and Shared Key", zap.Error(err))
			return "", err
		}
	} else {
		gt, err := gDA.GetGatewayByPublicKeySig(gateway.PublicKeySig)
		gateway.Id = gt.Id
		err = tx_gateway_verifier.SharedKeyAndGatewayUpdate(verifierUSer, gateway, db)
		if err != nil {
			zap.L().Error("Error while registering gateway and Shared Key", zap.Error(err))
			return "", err
		}
	}
	return cipherTextStr, nil
}
