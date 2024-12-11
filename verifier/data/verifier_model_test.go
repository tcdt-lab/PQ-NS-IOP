package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"verifier/config"
)

func generateTwoVerifiers() (Verifier, Verifier) {
	verifier1 := Verifier{
		Ip:            "v_IP1",
		Port:          "2222",
		PublicKeySig:  "v_PublicKeySig1",
		SigScheme:     "v_SigScheme1",
		SymmetricKey:  "v_SymmetricKey1",
		TrustScore:    0.5,
		IsInCommittee: false,
	}
	verifier2 := Verifier{
		Ip:            "v_IP2",
		Port:          "1111",
		PublicKeySig:  "v_PublicKeySig2",
		SigScheme:     "v_SigScheme2",
		SymmetricKey:  "v_SymmetricKey2",
		TrustScore:    0.6,
		IsInCommittee: false,
	}
	return verifier1, verifier2
}
func Test_VerfiersCrud(t *testing.T) {
	c, err := config.ReadYaml()
	assert.NoError(t, err, "Error in ReadYaml")
	v1, v2 := generateTwoVerifiers()
	var verifiersList []Verifier
	db, err := getDBConnection(*c)
	assert.NoError(t, err, "Error opening database")
	t.Run("AddVerifier", func(t *testing.T) {
		_, err := AddNewVerifier(&v1, db)
		assert.NoError(t, err, "Error adding verifier")
		_, err = AddNewVerifier(&v2, db)
		assert.NoError(t, err, "Error adding verifier")
	})
	t.Run("GetVerifiers", func(t *testing.T) {
		verifiersList, err = GetVerifiers(db)
		assert.NoError(t, err, "Error getting verifiers")
		assert.Greater(t, len(verifiersList), 0, "Expected at least one verifier, got 0")
		t.Log(verifiersList)
	})
	t.Run("UpdateVerifier", func(t *testing.T) {
		v1.TrustScore = 0.7
		v1.Ip = "v_IP1_updated"
		v1.IsInCommittee = true
		v1.PublicKeySig = "v_PublicKeySig1_updated"
		v1.SigScheme = "v_SigScheme1_updated"
		v1.SymmetricKey = "v_SymmetricKey1_updated"
		v1.Port = "9999"
		v1.Id = verifiersList[0].Id

		_, err = UpdateVerifier(&v1, db)
		assert.NoError(t, err, "Error updating verifier")
		verifiersList, err = GetVerifiers(db)
		assert.NoError(t, err, "Error getting verifiers")
		assert.Equal(t, "v_IP1_updated", verifiersList[0].Ip)
	})
	t.Run("GetVerfiersInCommittee", func(t *testing.T) {
		verifiersList, err = GetVerifiersInCommittee(db)
		assert.NoError(t, err, "Error getting verifiers")
		assert.Greater(t, len(verifiersList), 0, "Expected at least one verifier, got 0")
		t.Log(verifiersList)
	})
	t.Run("GetVerifierByPublicKeySig", func(t *testing.T) {
		verifier, err := GetVerifierByPublicKeySig(db, v1.PublicKeySig)
		assert.NoError(t, err, "Error getting verifier")
		assert.Equal(t, "v_IP1_updated", verifier.Ip)
	})
	t.Run("GetVerifierByIP", func(t *testing.T) {
		verifier, err := GetVerifierByIp(db, v1.Ip)
		assert.NoError(t, err, "Error getting verifier")
		assert.Equal(t, "v_IP1_updated", verifier.Ip)
	})

	t.Run("GetVerifier", func(t *testing.T) {
		verifier, err := GetVerifier(db, verifiersList[0].Id)
		assert.NoError(t, err, "Error getting verifier")
		assert.Equal(t, "v_IP1_updated", verifier.Ip)
	})
	t.Run("RemoveVerifier", func(t *testing.T) {
		verifiersList, _ := GetVerifiers(db)
		_, err = RemoveVerifier(db, verifiersList[0].Id)

		lenBeforeRemove := len(verifiersList)
		assert.NoError(t, err, "Error removing verifier")
		verifiersList, err = GetVerifiers(db)
		assert.NoError(t, err, "Error getting verifiers")
		assert.Equal(t, lenBeforeRemove-1, len(verifiersList))
		t.Log(verifiersList)
	})

}
