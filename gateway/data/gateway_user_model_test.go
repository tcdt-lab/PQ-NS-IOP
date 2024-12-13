package data

import (
	"gateway/config"
	"testing"
)

func TestAddUser(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	var user User
	user.Salt = "test"
	user.Password = "test"
	user.PublicKey = "test"
	user.SecretKey = "test"
	_, err = AddUser(db, user)
	if err != nil {
		t.Errorf("Error adding logic: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	user, err := GetUser(db, 1)
	if err != nil {
		t.Errorf("Error getting logic: %v", err)
	}
	t.Log(user)
}

func TestUpdateUser(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	var user User
	user.Id = 1
	user.Salt = "test_update"
	user.Password = "test_update"
	user.PublicKey = "test_update"
	user.SecretKey = "test_update"
	_, err = UpdateUser(db, user)
	if err != nil {
		t.Errorf("Error updating logic: %v", err)
	}
}

func TestRemoveUser(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	_, err = RemoveUser(db, 1)
	if err != nil {
		t.Errorf("Error removing logic: %v", err)
	}
}

func TestGetUsers(t *testing.T) {

	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	users, err := GetUsers(db)
	if err != nil {
		t.Errorf("Error getting users: %v", err)
	}
	t.Log(users)
}

func TestGetUserByID(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	user, err := GetUserByID(db, 1)
	if err != nil {
		t.Errorf("Error getting logic by ID: %v", err)
	}
	t.Log(user)
}

func TestGetUserByPublicKey(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	user, err := GetUserByPublicKey(db, "test")
	if err != nil {
		t.Errorf("Error getting logic by public key: %v", err)
	}
	t.Log(user)
}

func TestGetUserByPassword(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	user, err := GetUserByPassword(db, "test")
	if err != nil {
		t.Errorf("Error getting logic by Password: %v", err)
	}
	t.Log(user)
}
