package vv_key_distribution

import (
	"database/sql"
	b64 "encoding/base64"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/verifier_verifier"
	"verifier/config"
	"verifier/data"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func ApplyKeyDistributionRequest(db *sql.DB, c *config.Config, msgInfo pkg.MessageInfo) (string, error) {

	vud := data_access.GenerateVerifierUserDA(db)
	VDa := data_access.GenerateVerifierDA(db)
	cachHandler := data_access.GenerateCacheHandlerDA()
	adminId, _ := cachHandler.GetUserAdminId()
	verifierUser, _ := vud.GetVerifierUser(adminId)
	util := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)

	vvKeyDistributionReq := msgInfo.Params.(verifier_verifier.VVerifierKeyDistributionRequest)
	cipherTextBytes, sharedKey, _ := util.AsymmetricHandler.KemGenerateSecretKey(verifierUser.SecretKeyKem, vvKeyDistributionReq.VerifierPublicKeyKem, "", c.Security.KEMScheme)
	cipherTextStr := b64.StdEncoding.EncodeToString(cipherTextBytes)

	verifier := data.Verifier{}
	verifier.PublicKeySig = vvKeyDistributionReq.VerifierPublicKeySignature
	verifier.Ip = vvKeyDistributionReq.VerifierIP
	verifier.Port = vvKeyDistributionReq.VerifierPort
	verifier.SymmetricKey = util.AesHandler.ConvertKeyBytesToStr64(sharedKey)
	verifier.SigScheme = vvKeyDistributionReq.VerifierSignatureScheme
	verifier.IsInCommittee = false
	verifier.TrustScore = 0

	_, err := VDa.AddUpdateVerifier(verifier)
	if err != nil {
		return "", err
	}

	return cipherTextStr, nil
}
