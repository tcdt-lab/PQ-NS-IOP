package key_distribution

import (
	"database/sql"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
)

func ApplyGatewayVerifierKeyDistributionResponse(msgInfo pkg.MessageInfo, db *sql.DB) error {
	cfg, err := config.ReadYaml()
	if err != nil {
		return err
	}
	gul := data_access.GenerateGatewayUserDA(db)
	vl := data_access.GenerateVerifierDA(db)
	cacheHandler := data_access.NewCacheHandlerDA()
	adminId, _ := cacheHandler.GetUserAdminId()
	gtUser, err := gul.GetGatewayUser(adminId)
	if err != nil {
		zap.L().Error("Error while getting gateway user", zap.Error(err))
		return err
	}
	pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	gvKeyDistributionRes := msgInfo.Params.(gateway_verifier.GatewayVerifierKeyDistributionResponse)
	_, sharedKey, err := pkgUtil.AsymmetricHandler.KemGenerateSecretKey(gtUser.SecretKeyKem, gvKeyDistributionRes.PublicKeyKem, gvKeyDistributionRes.CipherText, cfg.Security.KEMScheme)

	bootstrapVerifier, err := vl.GetVerifierByIpAndPort(cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)

	bootstrapVerifier.SymmetricKey = pkgUtil.AesHandler.ConvertKeyBytesToStr64(sharedKey)
	zap.L().Info("Symmetric key is generated", zap.String("Symmetric key", bootstrapVerifier.SymmetricKey))
	//bootstrapVerifier.PublicKey = gvKeyDistributionRes.PublicKeyKem
	_, err = vl.UpdateVerifier(bootstrapVerifier)

	if err != nil {
		zap.L().Error("Error while updating verifier", zap.Error(err))
		return err
	}
	return nil
}
