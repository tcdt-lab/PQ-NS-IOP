package message_applier

import (
	"gateway/config"
	"gateway/logic"
	"gateway/message_handler/util"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
)

func ApplyGatewayVerifierKeyDistributionResponse(msgData pkg.MessageData) error {
	cfg, err := config.ReadYaml()
	if err != nil {
		return err
	}
	gul := logic.GatewayUserLogic{}
	vl := logic.VerifierLogic{}
	gtUser, err := gul.GetGatewayUser(1)
	if err != nil {
		zap.L().Error("Error while getting gateway user", zap.Error(err))
		return err
	}
	pkgUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	gvKeyDistributionRes := msgData.MsgInfo.Params.(gateway_verifier.GatewayVerifierKeyDistributionResponse)
	_, sharedKey, err := pkgUtil.AsymmetricHandler.KemGenerateSecretKey(gtUser.SecretKeyKem, gvKeyDistributionRes.PublicKeyKem, gvKeyDistributionRes.CipherText, cfg.Security.KEMScheme)

	bootstrapVerifier, err := vl.GetVerifierByIpAndPort(cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)

	bootstrapVerifier.SymmetricKey = pkgUtil.AesHandler.ConvertKeyBytesToStr64(sharedKey)
	bootstrapVerifier.PublicKey = gvKeyDistributionRes.PublicKeyKem
	_, err = vl.UpdateVerifier(bootstrapVerifier)

	if err != nil {
		zap.L().Error("Error while updating verifier", zap.Error(err))
		return err
	}
	return nil
}
