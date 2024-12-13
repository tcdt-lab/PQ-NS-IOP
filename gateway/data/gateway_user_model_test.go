package data

import (
	b64 "encoding/base64"
	"gateway/config"
	"github.com/stretchr/testify/assert"
	"os"
	"test.org/cryptography/pkg/asymmetric"
	"test.org/cryptography/pkg/symmetric"
	"test.org/protocol/pkg"
	"testing"
)

func generateFakeGatewayUser(c *config.Config) GatewayUser {

	protocolUtil := generateProtocolUtil(*c)

	var user GatewayUser
	user.Salt = b64.StdEncoding.EncodeToString([]byte("Testsalt"))
	user.Password = b64.StdEncoding.EncodeToString([]byte(os.Getenv("PQ_NS_IOP_GU_PASS")))
	secKeyDsa, pubKeyDsa, err := protocolUtil.AsymmetricHandler.DSKeyGen(c.Security.DSAScheme)
	if err != nil {
		panic(err)
	}
	secKeyKem, pubKeyKem, err := protocolUtil.AsymmetricHandler.KEMKeyGen(c.Security.KEMScheme)
	if err != nil {
		panic(err)
	}
	user.PublicKeyDsa = pubKeyDsa
	user.SecretKeyDsa = secKeyDsa
	user.PublicKeyKem = pubKeyKem
	user.SecretKeyKem = secKeyKem
	user.Dsa_scheme = c.Security.DSAScheme
	user.Kem_scheme = c.Security.KEMScheme
	return user
}

func generateProtocolUtil(c config.Config) pkg.ProtocolUtil {
	var util pkg.ProtocolUtil
	util.AesHandler = symmetric.AesGcm{}
	util.AsymmetricHandler = asymmetric.NewAsymmetricHandler(c.Security.CryptographyScheme)
	util.HmacHandler = symmetric.HMAC{}
	return util
}

func Test_GatewayUserCrud(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		panic(err)
	}
	user := generateFakeGatewayUser(c)
	var users []GatewayUser
	db := getDBConnection(*c)
	defer db.Close()
	t.Run("AddGatewayUser", func(t *testing.T) {
		_, err := AddNewGatewayUser(db, user)
		assert.NoError(t, err, "Error adding gateway user")
	})
	t.Run("GetAllGatewayUsers", func(t *testing.T) {
		users, err = GetAllGatewayUsers(db)
		assert.NoError(t, err, "Error getting all gateway users")
		assert.NotEmpty(t, users, "No gateway users found")
		t.Log("Users: ", users)
	})
	t.Run("GetGatewayUser", func(t *testing.T) {
		userNew, err := GetGatewayUser(db, users[0].Id)
		assert.NoError(t, err, "Error getting gateway user")
		t.Log("User: ", userNew)
	})
	t.Run("GetGatewayUserByPassword", func(t *testing.T) {
		userNew, err := GetGatewayUserByPassword(db, user.Password)
		assert.NoError(t, err, "Error getting gateway user by password")
		t.Log("User: ", userNew)
	})
	t.Run("GetGatewayUserByPublicKeyDsa", func(t *testing.T) {
		userNew, err := GetGatewayUserByPublicKeyDsa(db, user.PublicKeyDsa)
		assert.NoError(t, err, "Error getting gateway user by public key dsa")
		t.Log("User: ", userNew)
	})
	t.Run("GetGatewayUserByPublicKeyKem", func(t *testing.T) {
		_, err := GetGatewayUserByPublicKeyKem(db, user.PublicKeyKem)
		assert.NoError(t, err, "Error getting gateway user by public key kem")
	})
	t.Run("UpdateGatewayUser", func(t *testing.T) {
		users[0].Salt = b64.StdEncoding.EncodeToString([]byte("TestsaltUpdated"))
		_, err := UpdateGatewayUser(db, users[0])
		assert.NoError(t, err, "Error updating gateway user")
		UpdatedUSer, err := GetGatewayUser(db, users[0].Id)
		assert.NoError(t, err, "Error getting updated gateway user")
		assert.Equal(t, users[0].Salt, UpdatedUSer.Salt, "Salt not updated")
	})
	t.Run("DeleteGatewayUser", func(t *testing.T) {
		_, err := DeleteGatewayUser(db, user.Id)
		assert.NoError(t, err, "Error deleting gateway user")
	})
}
