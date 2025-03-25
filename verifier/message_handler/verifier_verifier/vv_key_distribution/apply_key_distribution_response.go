package vv_key_distribution

import (
	"database/sql"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/verifier_verifier"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func ApplyKeyDistributionResponse(msgInfo pkg.MessageInfo, db *sql.DB, c *config.Config) error {

	vud := data_access.GenerateVerifierUserDA(db)
	VDa := data_access.GenerateVerifierDA(db)
	cachHandler := data_access.NewCacheHandlerDA()
	adminId, _ := cachHandler.GetUserAdminId()
	verifierUser, _ := vud.GetVerifierUser(adminId)
	util := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)

	vvKeyDistributionRes := msgInfo.Params.(verifier_verifier.VVerifierKeyDistributionResponse)

	_, sharedKey, _ := util.AsymmetricHandler.KemGenerateSecretKey(verifierUser.SecretKeyKem, vvKeyDistributionRes.PublicKeyKem, vvKeyDistributionRes.CipherText, c.Security.KEMScheme)
	bootstrapVerifier, err := VDa.GetVerifierByIpAndPort(c.BootstrapNode.Ip, c.BootstrapNode.Port)
	if err != nil {
		return err
	}
	bootstrapVerifier.SymmetricKey = util.AesHandler.ConvertKeyBytesToStr64(sharedKey)
	_, err = VDa.UpdateVerifier(bootstrapVerifier)
	if err != nil {
		return err
	}
	return nil
}
