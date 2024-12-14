package data

import (
	"database/sql"
	"fmt"
	"gateway/config"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func getDBConnection(c config.Config) *sql.DB {
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	if err != nil {
		return nil
	}
	return db
}
func TestAddGateway(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	var gateway Gateway
	gateway.Ip = "test"
	gateway.Port = "8080"
	gateway.PublicKey = "test"
	gateway.Ticket = "test"
	gateway.SymmetricKey = "test"
	_, err = AddGateway(db, gateway)
	if err != nil {
		t.Errorf("Error adding gateway: %v", err)
	}
}
func TestGetGateway(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	gateways, err := GetGateways(db)
	if err != nil {
		t.Errorf("Error getting gateways: %v", err)
	}
	if len(gateways) == 0 {
		t.Errorf("Expected at least one gateway, got 0")
	}
	fmt.Println(gateways)
}

func TestUpdateGateway(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	var gateway Gateway
	gateway.Id = 2
	gateway.Ip = "test_update"
	gateway.Port = "9090"
	gateway.PublicKey = "test_update"
	gateway.Ticket = "test_update"
	gateway.SymmetricKey = "test_update"
	_, err = UpdateGateway(db, gateway)
	if err != nil {
		t.Errorf("Error updating gateway: %v", err)
	}
}

func TestRemoveGateway(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	_, err = RemoveGateway(db, 1)
	if err != nil {
		t.Errorf("Error removing gateway: %v", err)
	}
}

func TestGetGatewayByIP(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	gateway, err := GetGatewayByIP(db, "test_update")
	if err != nil {
		t.Errorf("Error getting gateway by IP: %v", err)
	}
	fmt.Println(gateway)
}

func TestGetGatewayByID(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	gateway, err := GetGateway(db, 2)
	if err != nil {
		t.Errorf("Error getting gateway by ID: %v", err)
	}
	fmt.Println(gateway)
}

func TestGetGatewayByIpandPort(t *testing.T) {
	c, err := config.ReadYaml()
	if err != nil {
		t.Errorf("Error reading config.yaml file: %v", err)
	}
	db := getDBConnection(*c)

	if err != nil {
		t.Errorf("Error opening database: %v", err)
	}
	defer db.Close()
	gateway, err := GetGatewayByIpAndPort(db, "test", "8080")
	if err != nil {
		t.Errorf("Error getting gateway by IP and Port: %v", err)
	}
	fmt.Println(gateway)
}
