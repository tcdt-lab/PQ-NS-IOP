package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"verifier/config"
)

func generateVerifierUSerModels() (VerifierUser, VerifierUser) {
	verifierUser1 := VerifierUser{
		Salt:      "v_salt1",
		Password:  "v_password1",
		SecretKey: "v_secret",
		PublicKey: "v_publicKey1",
	}
	verifierUser2 := VerifierUser{
		Salt:      "v_salt2",
		Password:  "v_password2",
		SecretKey: "v_secret_2",
		PublicKey: "v_publicKey2",
	}
	return verifierUser1, verifierUser2
}

func Test_CrudVerifierUserModel(t *testing.T) {
	c, err := config.ReadYaml()
	assert.NoError(t, err, "Error in ReadYaml")
	verifierUser1, verifierUser2 := generateVerifierUSerModels()
	db, err := getDBConnection(*c)
	assert.NoError(t, err, "Error opening database")
	var verifierUserList []VerifierUser
	t.Run("AddVerifierUser", func(t *testing.T) {
		_, err := AddVerifierUser(&verifierUser1, db)
		assert.NoError(t, err, "Error adding verifier user")
		_, err = AddVerifierUser(&verifierUser2, db)
		assert.NoError(t, err, "Error adding verifier user")
	})
	t.Run("GetVerifierUsers", func(t *testing.T) {
		verifierUserList, err = GetVerifierUsers(db)
		assert.NoError(t, err, "Error getting verifier users")
		assert.Greater(t, len(verifierUserList), 0, "Expected at least one verifier user, got 0")
		t.Log(verifierUserList)
	})
	t.Run("GetVerifierUser", func(t *testing.T) {
		verifierUser, err := GetVerifierUser(db, verifierUserList[0].Id)
		assert.NoError(t, err, "Error getting verifier user")
		assert.Equal(t, verifierUserList[0], verifierUser, "Expected same verifier user")
	})
	t.Run("UpdateVerifierUser", func(t *testing.T) {
		verifierUser1.Salt = "v_salt1_updated"
		verifierUser1.Password = "v_password1_updated"
		verifierUser1.SecretKey = "v_secret_updated"
		verifierUser1.PublicKey = "v_publicKey1_updated"
		verifierUser1.Id = verifierUserList[0].Id

		_, err = UpdateVerifierUser(&verifierUser1, db)
		verifierUserList, _ = GetVerifierUsers(db)
		assert.NoError(t, err, "Error updating verifier user")
		assert.Equal(t, verifierUser1.Password, verifierUserList[0].Password, "Expected same verifier user")
		assert.Equal(t, verifierUser1.Salt, verifierUserList[0].Salt, "Expected same verifier user")
		assert.Equal(t, verifierUser1.SecretKey, verifierUserList[0].SecretKey, "Expected same verifier user")
	})
	t.Run("GetVerifierUserByPassword", func(t *testing.T) {
		verifierUser, err := GetVerifierUserByPassword(db, verifierUser1.Password)
		assert.NoError(t, err, "Error getting verifier user")
		assert.Equal(t, verifierUser1.PublicKey, verifierUser.PublicKey, "Expected same verifier user")
	})

	t.Run("GetVerifierUserByPublicKey", func(t *testing.T) {
		verifierUser, err := GetVerifierUserByPublicKey(db, verifierUser1.PublicKey)
		assert.NoError(t, err, "Error getting verifier user")
		assert.Equal(t, verifierUser1.SecretKey, verifierUser.SecretKey, "Expected same verifier user")
	})
	t.Run("GetVerifierUserBySecretKey", func(t *testing.T) {
		verifierUser, err := GetVerifierUserBySecretKey(db, verifierUser1.SecretKey)
		assert.NoError(t, err, "Error getting verifier user")
		assert.Equal(t, verifierUser1.Password, verifierUser.Password, "Expected same verifier user")
	})
	t.Run("RemoveVerifierUser", func(t *testing.T) {
		_, err := RemoveVerifierUser(db, verifierUserList[0].Id)
		assert.NoError(t, err, "Error removing verifier user")
		verifierUserList, _ = GetVerifierUsers(db)
		assert.NotEqual(t, verifierUserList[0].Id, verifierUser1.Id, "Expected different verifier user")

	})
}
