package logic

import (
	"database/sql"
	"encoding/base64"
	"gateway/config"
	"gateway/data"
	"gateway/data_access"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/cryptography/pkg/asymmetric"
)

func InintStepLogic(cfg *config.Config, db *sql.DB) (int64, int64, error) {

	zap.L().Info("Initiating bootstrap and admin")
	admin := data.GatewayUser{}
	cacheHAndler := data_access.NewCacheHandlerDA()
	pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	salt, err := pkgUtil.PBKDF2Handler.GeneratingSalt(32)
	if err != nil {
		zap.L().Error("Error while generating salt", zap.Error(err))
		return 0, 0, err
	}

	vDa := data_access.GenerateVerifierDA(db)
	verifier := data.Verifier{}
	verifier.PublicKey = cfg.BootstrapNode.PubKeySig
	verifier.Ip = cfg.BootstrapNode.Ip
	verifier.Port = cfg.BootstrapNode.Port

	bootstrapExists, err := vDa.IfVerifierExistsByIpandPort(verifier)
	if bootstrapExists {

		adminId, _ := cacheHAndler.GetUserAdminId()
		bootstrapId, _ := cacheHAndler.GetBootstrapVerifierId()
		zap.L().Info("Bootstrap and Admin already exists", zap.Int64("BootstrapId", bootstrapId), zap.Int64("AdminId", adminId))
		return bootstrapId, adminId, nil
	}

	asymHandler := asymmetric.NewAsymmetricHandler(cfg.Security.CryptographyScheme)
	skDSAStr, pkDSAStr, err := asymHandler.DSKeyGen(cfg.Security.DSAScheme)
	skKEMStr, pkKEMStr, err := asymHandler.KEMKeyGen(cfg.Security.KEMScheme)
	if err != nil {
		zap.L().Error("Error while generating DSA key pair", zap.Error(err))
	}

	admin.SecretKeyKem = skKEMStr
	admin.SecretKeyDsa = skDSAStr
	admin.PublicKeyKem = pkKEMStr
	admin.PublicKeyDsa = pkDSAStr
	admin.Salt = base64.StdEncoding.EncodeToString(salt)
	admin.Dsa_scheme = cfg.Security.DSAScheme
	admin.Kem_scheme = cfg.Security.KEMScheme

	gatewayTransaction := data.GatewayTransaction{}
	gatewayId, bootstrapId, err := gatewayTransaction.InitStep(admin, verifier)
	if err != nil {
		zap.L().Error("Error while initializing gateway", zap.Error(err))
		return 0, 0, err
	}
	cacheHAndler.SetBootstrapVerifierId(bootstrapId)
	cacheHAndler.SetUserAdminId(gatewayId)
	zap.L().Info("Bootstrap and Admin generated", zap.String("Step:", "inint_logic_step"), zap.Int64("BootstrapId", bootstrapId), zap.Int64("AdminId", gatewayId))
	return bootstrapId, gatewayId, nil
}
