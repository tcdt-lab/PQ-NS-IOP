package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"verifier/config"
)

func getDBConnection(c config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCrudGateways(t *testing.T) {

	c, err := config.ReadYaml()
	assert.NoError(t, err, "Error in ReadYaml")
	db, err := getDBConnection(*c)
	assert.NoError(t, err, "Error opening database")
	defer db.Close()
	t.Run("AddGateway", func(t *testing.T) {
		var gateway Gateway
		gateway.Ip = "test_ip"
		gateway.Port = "8080"
		gateway.Ticket = "test_ticket"
		gateway.SymmetricKey = "test_SK"
		gateway.KemScheme = "test_KS"
		gateway.SigScheme = "test_SS"
		gateway.PublicKeyKem = "test_PKK"
		gateway.PublicKeySig = "test_PKS"
		_, err = AddGateway(db, gateway)

		assert.NoError(t, err, "Error adding gateway")

		var gateway2 Gateway
		gateway2.Ip = "test2_ip"
		gateway2.Port = "8080"
		gateway2.Ticket = "test2_ticket"
		gateway2.SymmetricKey = "test2_SK"
		gateway2.KemScheme = "test2_KS"
		gateway2.SigScheme = "test2_SS"
		gateway2.PublicKeyKem = "test2_PKK"
		gateway2.PublicKeySig = "test2_PKS"
		_, err = AddGateway(db, gateway2)

		assert.NoError(t, err, "Error adding gateway")

	})
	t.Run("GetGateways", func(t *testing.T) {
		gateways, err := GetGateways(db)
		assert.NoError(t, err, "Error getting gateways")
		assert.Greater(t, len(gateways), 0, "Expected at least one gateway, got 0")
		t.Log(gateways)
	})
	t.Run("UpdateGateway", func(t *testing.T) {
		gateways, err := GetGateways(db)
		assert.NoError(t, err, "Error getting gateways")
		assert.Greater(t, len(gateways), 0, "Expected at least one gateway, got 0")
		gateways[0].Ip = "updated_ip"
		gateways[0].Port = "1010U"
		gateways[0].Ticket = "updated_ticket"
		gateways[0].SymmetricKey = "updated_SK"
		gateways[0].KemScheme = "updated_KS"
		gateways[0].SigScheme = "updated_SS"
		gateways[0].PublicKeyKem = "updated_PKK"
		gateways[0].PublicKeySig = "updated_PKS"
		_, err = UpdateGateway(db, gateways[0])
		assert.NoError(t, err, "Error updating gateway")

		updatedGt, err := GetGateway(db, gateways[0].Id)
		assert.NoError(t, err, "Error getting updated gateway")
		assert.Equal(t, "updated_ip", updatedGt.Ip)
		assert.Equal(t, "1010U", updatedGt.Port)
		assert.Equal(t, "updated_ticket", updatedGt.Ticket)
		assert.Equal(t, "updated_SK", updatedGt.SymmetricKey)
		assert.Equal(t, "updated_KS", updatedGt.KemScheme)
		assert.Equal(t, "updated_SS", updatedGt.SigScheme)
		assert.Equal(t, "updated_PKK", updatedGt.PublicKeyKem)
		assert.Equal(t, "updated_PKS", updatedGt.PublicKeySig)
	})

	t.Run("RemoveGateway", func(t *testing.T) {
		gateways, err := GetGateways(db)
		assert.NoError(t, err, "Error getting gateways")
		assert.Greater(t, len(gateways), 0, "Expected at least one gateway, got 0")
		_, err = RemoveGateway(db, gateways[0].Id)
		assert.NoError(t, err, "Error removing gateway")
		newGtList, err := GetGateways(db)
		assert.NoError(t, err, "Error getting gateways")
		assert.Equal(t, len(gateways)-1, len(newGtList))
		t.Log(newGtList)
	})

	t.Run("GetGatewayByIP", func(t *testing.T) {
		gateways, err := GetGateways(db)
		assert.NoError(t, err, "Error getting gateways")
		assert.Greater(t, len(gateways), 0, "Expected at least one gateway, got 0")
		gateway, err := GetGatewayByIp(db, gateways[0].Ip)
		assert.NoError(t, err, "Error getting gateway by IP")
		assert.Equal(t, gateways[0].Ip, gateway.Ip)
	})

	t.Run("GetGatewayByPublicKeySig", func(t *testing.T) {
		gateways, err := GetGateways(db)
		assert.NoError(t, err, "Error getting gateways")
		assert.Greater(t, len(gateways), 0, "Expected at least one gateway, got 0")
		gateway, err := GetGatewayByPublicKeySig(db, gateways[0].PublicKeySig)
		assert.NoError(t, err, "Error getting gateway by PublicKeySig")
		assert.Equal(t, gateways[0].PublicKeySig, gateway.PublicKeySig)
	})

	t.Run("GetGatewayByPublicKeyKem", func(t *testing.T) {
		gateways, err := GetGateways(db)
		assert.NoError(t, err, "Error getting gateways")
		assert.Greater(t, len(gateways), 0, "Expected at least one gateway, got 0")
		gateway, err := GetGatewayByPublicKeyKem(db, gateways[0].PublicKeyKem)
		assert.NoError(t, err, "Error getting gateway by PublicKeyKem")
		assert.Equal(t, gateways[0].PublicKeyKem, gateway.PublicKeyKem)
	})
}
