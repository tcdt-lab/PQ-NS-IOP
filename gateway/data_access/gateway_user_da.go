package data_access

import (
	"database/sql"
	"gateway/data"
)

type GatewayUserDA struct {
	db *sql.DB
}

func GenerateGatewayUserDA(db *sql.DB) *GatewayUserDA {
	return &GatewayUserDA{
		db: db,
	}
}

func (gul *GatewayUserDA) GetGatewayUsers() ([]data.GatewayUser, error) {
	return data.GetAllGatewayUsers(gul.db)
}

func (gul *GatewayUserDA) AddGatewayUser(gatewayUser data.GatewayUser) (int64, error) {
	return data.AddNewGatewayUser(gul.db, gatewayUser)
}

func (gul *GatewayUserDA) UpdateGatewayUser(gatewayUser data.GatewayUser) (int64, error) {
	return data.UpdateGatewayUser(gul.db, gatewayUser)
}

func (gul *GatewayUserDA) RemoveGatewayUser(id int) (int64, error) {
	return data.DeleteGatewayUser(gul.db, id)
}

func (gul *GatewayUserDA) GetGatewayUser(id int64) (data.GatewayUser, error) {
	return data.GetGatewayUser(gul.db, id)
}

func (gul *GatewayUserDA) GetGatewayUserByPublicKeyDsa(pubklicKeyDSA string) (data.GatewayUser, error) {
	return data.GetGatewayUserByPublicKeyDsa(gul.db, pubklicKeyDSA)
}

func (gul *GatewayUserDA) GetGatewayUserByPublicKeyKem(pubklicKeyKEM string) (data.GatewayUser, error) {
	return data.GetGatewayUserByPublicKeyKem(gul.db, pubklicKeyKEM)
}
