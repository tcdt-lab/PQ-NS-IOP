package logic

import (
	"encoding/base64"
	"gateway/config"
	"gateway/data"
	"gateway/data_access"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/cryptography/pkg/asymmetric"
)

func KeyGenLogic(cfg *config.Config) (int64, int64, error) {

	admin := data.GatewayUser{}
	pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	salt, err := pkgUtil.PBKDF2Handler.GeneratingSalt(32)
	if err != nil {
		zap.L().Error("Error while generating salt", zap.Error(err))
		return 0, 0, err
	}

	vDa := data_access.VerifierDA{}
	verifier := data.Verifier{}
	verifier.PublicKey = cfg.BootstrapNode.PubKeySig
	verifier.Ip = cfg.BootstrapNode.Ip
	verifier.Port = cfg.BootstrapNode.Port

	bootstrapExists, err := vDa.IfVerifierExistsByIpandPort(verifier)
	if bootstrapExists {
		zap.L().Info("Bootstrap verifier already exists")
		return 0, 0, nil
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

	return bootstrapId, gatewayId, nil
}
