package data

import (
	b64 "encoding/base64"
	"github.com/stretchr/testify/assert"
	"test.org/cryptography/pkg/asymmetric"
	"test.org/cryptography/pkg/symmetric"
	"test.org/protocol/pkg"
	"testing"
	"verifier/config"
)

func generateProtocolUtil(c *config.Config) pkg.ProtocolUtil {
	var util pkg.ProtocolUtil
	util.AesHandler = symmetric.AesGcm{}
	util.AsymmetricHandler = asymmetric.NewAsymmetricHandler(c.Security.CryptographyScheme)
	util.HmacHandler = symmetric.HMAC{}
	return util

}
func generateVerifierUSerModels(protocolUtil pkg.ProtocolUtil, c *config.Config) (VerifierUser, VerifierUser) {
	var verifierUser1 = VerifierUser{}
	var verifierUser2 = VerifierUser{}
	verifierUser1.Salt = b64.StdEncoding.EncodeToString([]byte("salt1"))
	verifierUser1.Password = b64.StdEncoding.EncodeToString([]byte("vu_password1"))
	secKeySig1, pubKeySig1, err := protocolUtil.AsymmetricHandler.DSKeyGen(c.Security.MlDSAScheme)
	if err != nil {
		panic(err)
	}
	secKeyKem1, pubKeyKem1, err := protocolUtil.AsymmetricHandler.KEMKeyGen(c.Security.MlKEMScheme)
	if err != nil {
		panic(err)
	}
	verifierUser1.PublicKeyKem = pubKeyKem1
	verifierUser1.PublicKeySig = pubKeySig1
	verifierUser1.SecretKeyKem = secKeyKem1
	verifierUser1.SecretKeySig = secKeySig1

	verifierUser2.Salt = b64.StdEncoding.EncodeToString([]byte("salt2"))
	verifierUser2.Password = "kk1234"
	secKeySig2, pubKeySig2, err := protocolUtil.AsymmetricHandler.DSKeyGen(c.Security.MlDSAScheme)
	if err != nil {
		panic(err)
	}
	secKeyKem2, pubKeyKem2, err := protocolUtil.AsymmetricHandler.KEMKeyGen(c.Security.MlKEMScheme)
	if err != nil {
		panic(err)
	}
	verifierUser2.PublicKeyKem = pubKeyKem2
	verifierUser2.PublicKeySig = pubKeySig2
	verifierUser2.SecretKeyKem = secKeyKem2
	verifierUser2.SecretKeySig = secKeySig2

	return verifierUser1, verifierUser2
}

func Test_CrudVerifierUserModel(t *testing.T) {
	c, err := config.ReadYaml()
	assert.NoError(t, err, "Error in ReadYaml")
	verifierUser1, verifierUser2 := generateVerifierUSerModels(generateProtocolUtil(c), c)
	db, err := getDBConnection(*c)
	assert.NoError(t, err, "Error opening database")
	var verifierUserList []VerifierUser
	t.Run("AddVerifierUser", func(t *testing.T) {
		_, err := AddVerifierUser(&verifierUser1, db)
		assert.NoError(t, err, "Error adding verifier_verifier user")
		_, err = AddVerifierUser(&verifierUser2, db)
		assert.NoError(t, err, "Error adding verifier_verifier user")
	})
	t.Run("GetVerifierUsers", func(t *testing.T) {
		verifierUserList, err = GetAllVerifierUsers(db)
		assert.NoError(t, err, "Error getting verifier_verifier users")
		assert.Greater(t, len(verifierUserList), 0, "Expected at least one verifier_verifier user, got 0")
		t.Log(verifierUserList)
	})
	t.Run("GetVerifierUser", func(t *testing.T) {
		verifierUser, err := GetVerifierUserById(db, verifierUserList[0].Id)
		assert.NoError(t, err, "Error getting verifier_verifier user")
		assert.Equal(t, verifierUserList[0], verifierUser, "Expected same verifier_verifier user")
	})
	t.Run("UpdateVerifierUser", func(t *testing.T) {
		verifierUser1.Salt = "v_salt1_updated"
		verifierUser1.Password = "v_password1_updated"
		verifierUser1.SecretKeyKem = "v_secret_kem_updated"
		verifierUser1.PublicKeyKem = "v_publicKey1_kem_updated"
		verifierUser1.PublicKeySig = "v_publicKey1_sig_updated"
		verifierUser1.SecretKeySig = "v_secret_sig_updated"
		verifierUser1.Id = verifierUserList[0].Id

		_, err = UpdateVerifierUser(&verifierUser1, db)
		verifierUserList, _ = GetAllVerifierUsers(db)
		assert.NoError(t, err, "Error updating verifier_verifier user")
		assert.Equal(t, verifierUser1.Password, verifierUserList[0].Password, "Expected same verifier_verifier user")
		assert.Equal(t, verifierUser1.Salt, verifierUserList[0].Salt, "Expected same verifier_verifier user")
		assert.Equal(t, verifierUser1.SecretKeySig, verifierUserList[0].SecretKeySig, "Expected same verifier_verifier user")
	})
	t.Run("GetVerifierUserByPassword", func(t *testing.T) {
		verifierUser, err := GetVerifierUserByPassword(db, verifierUser1.Password)
		assert.NoError(t, err, "Error getting verifier_verifier user")
		assert.Equal(t, verifierUser1.PublicKeySig, verifierUser.PublicKeySig, "Expected same verifier_verifier user")
	})

	t.Run("GetVerifierUserByPublicKey", func(t *testing.T) {
		verifierUser, err := GetVerifierUserByPublicKeySig(db, verifierUser1.PublicKeySig)
		assert.NoError(t, err, "Error getting verifier_verifier user")
		assert.Equal(t, verifierUser1.SecretKeySig, verifierUser.SecretKeySig, "Expected same verifier_verifier user")
	})

	t.Run("RemoveVerifierUser", func(t *testing.T) {
		_, err := DeleteVerifierUser(db, verifierUserList[0].Id)
		assert.NoError(t, err, "Error removing verifier_verifier user")
		verifierUserList, _ = GetAllVerifierUsers(db)
		assert.NotEqual(t, verifierUserList[0].Id, verifierUser1.Id, "Expected different verifier_verifier user")

	})
}
