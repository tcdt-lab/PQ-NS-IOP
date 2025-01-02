package data_access

import "gateway/data"

type GatewayUserDA struct {
}

func (gul *GatewayUserDA) GetGatewayUsers() ([]data.GatewayUser, error) {
	db, err := getDbConnection()
	if err != nil {
		return nil, err
	}
	return data.GetAllGatewayUsers(db)
}

func (gul *GatewayUserDA) AddGatewayUser(gatewayUser data.GatewayUser) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.AddNewGatewayUser(db, gatewayUser)
}

func (gul *GatewayUserDA) UpdateGatewayUser(gatewayUser data.GatewayUser) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.UpdateGatewayUser(db, gatewayUser)
}

func (gul *GatewayUserDA) RemoveGatewayUser(id int) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.DeleteGatewayUser(db, id)
}

func (gul *GatewayUserDA) GetGatewayUser(id int) (data.GatewayUser, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.GatewayUser{}, err
	}
	return data.GetGatewayUser(db, id)
}

func (gul *GatewayUserDA) GetGatewayUserByPublicKeyDsa(pubklicKeyDSA string) (data.GatewayUser, error) {

	db, err := getDbConnection()
	if err != nil {
		return data.GatewayUser{}, err
	}
	return data.GetGatewayUserByPublicKeyDsa(db, pubklicKeyDSA)
}

func (gul *GatewayUserDA) GetGatewayUserByPublicKeyKem(pubklicKeyKEM string) (data.GatewayUser, error) {

	db, err := getDbConnection()
	if err != nil {
		return data.GatewayUser{}, err
	}
	return data.GetGatewayUserByPublicKeyKem(db, pubklicKeyKEM)
}
