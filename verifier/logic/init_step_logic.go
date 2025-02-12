package logic

import (
	b64 "encoding/base64"
	"os"
	"verifier/config"
	"verifier/data"
	"verifier/data_access"
	"verifier/message_parser/util"
)

func InitStepLogic() (int64, int64, error) {
	//Here we create a new Instance of verfiifer user and build keys for it
	c, err := config.ReadYaml()
	if err != nil {
		panic(err)
	}
	var verfierUser = data.VerifierUser{}
	var verifierUserDataAccess = data_access.VerifierUserDA{}
	var bootstrapVerfier = data.Verifier{}
	var verifierDataAccess = data_access.VerifierDA{}
	protoUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	secKeyKem, pubKeyKem, err := protoUtil.AsymmetricHandler.KEMKeyGen(c.Security.KEMScheme)
	if err != nil {
		return 0, 0, err
	}
	secKeyDSA, pubKeyDSA, err := protoUtil.AsymmetricHandler.DSKeyGen(c.Security.DSAScheme)
	if err != nil {
		return 0, 0, err
	}
	verfierUser.SecretKeyKem = secKeyKem
	verfierUser.SecretKeySig = secKeyDSA
	verfierUser.PublicKeyKem = pubKeyKem
	verfierUser.PublicKeySig = pubKeyDSA
	verfierUser.Password = os.Getenv("PQ_NS_IOP_VU_PASS")
	saltByte, err := protoUtil.PBKDF2Handler.GeneratingSalt(64)
	if err != nil {
		return 0, 0, err
	}
	verfierUser.Salt = b64.StdEncoding.EncodeToString(saltByte)
	adminId, err := verifierUserDataAccess.AddVerifierUser(verfierUser)
	if err != nil {
		return 0, 0, err
	}

	bootstrapVerfier.Ip = c.BootstrapNode.Ip
	bootstrapVerfier.Port = c.BootstrapNode.Port
	bootstrapVerfier.PublicKeySig = c.BootstrapNode.PubKeySig

	bootstrapId, err := verifierDataAccess.AddVerifier(bootstrapVerfier)
	if err != nil {
		return 0, 0, err
	}
	return adminId, bootstrapId, nil
}
