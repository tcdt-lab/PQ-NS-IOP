package data

import (
	"gateway/config"
	"test.org/cryptography/pkg/asymmetric"

	"testing"
)

func TestAddVerifier(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	var verifier Verifier
	verifier.Ip = "127.0.0.1"
	verifier.Port = "50051"
	verifier.PublicKey = ""
	verifier.SymmetricKey = ""
	_, err = AddVerifier(db, verifier)
	if err != nil {
		t.Errorf("Error adding verifier_verifier: %v", err)
	}
}

func TestGetVerifiers(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	verifiers, err := GetVerifiers(db)
	if err != nil {
		t.Errorf("Error getting verifiers: %v", err)
	}
	if len(verifiers) == 0 {
		t.Errorf("Expected at least one verifier_verifier, got 0")
	}
	t.Log(verifiers)
}

func TestUpdateVerifier(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)
	asymmetricHandler := asymmetric.NewAsymmetricHandler(c.Security.CryptographyScheme)
	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	var verifier Verifier
	_, pkDSa, err := asymmetricHandler.DSKeyGen(c.Security.DSAScheme)
	verifier.Id = 1
	verifier.Ip = "127.0.0.1"
	verifier.Port = "50051"
	verifier.PublicKey = pkDSa
	verifier.SymmetricKey = ""
	_, err = UpdateVerifier(db, verifier)
	if err != nil {
		t.Errorf("Error updating verifier_verifier: %v", err)
	}
}

func TestRemoveVerifier(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	_, err = RemoveVerifier(db, 1)
	if err != nil {
		t.Errorf("Error removing verifier_verifier: %v", err)
	}
}

func TestGetVerifier(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	verifier, err := GetVerifier(db, 1)
	if err != nil {
		t.Errorf("Error getting verifier_verifier: %v", err)
	}
	t.Log(verifier)
}

func TestGetVerifierByIP(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	verifier, err := GetVerifierByIP(db, "test")
	if err != nil {
		t.Errorf("Error getting verifier_verifier by IP and Port: %v", err)
	}
	t.Log(verifier)
}

func TestGetVerifierByPublicKey(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	verifier, err := GetVerifierByPublicKey(db, "test")
	if err != nil {
		t.Errorf("Error getting verifier_verifier by public key: %v", err)
	}
	t.Log(verifier)
}

func TestGetVerifierByIPAndPort(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	verifier, err := GetVerifierByIpandPort(db, "127.0.0.1", "8080")
	if err != nil {
		t.Errorf("Error getting verifier_verifier by IP and Port: %v", err)
	}
	t.Log(verifier)
}
