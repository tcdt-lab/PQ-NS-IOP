package logic

import "gateway/data"

type GatewayUserLogic struct {
}

func (gul *GatewayUserLogic) GetGatewayUsers() ([]data.GatewayUser, error) {
	db, err := getDbConnection()
	if err != nil {
		return nil, err
	}
	return data.GetAllGatewayUsers(db)
}

func (gul *GatewayUserLogic) AddGatewayUser(gatewayUser data.GatewayUser) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.AddNewGatewayUser(db, gatewayUser)
}

func (gul *GatewayUserLogic) UpdateGatewayUser(gatewayUser data.GatewayUser) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.UpdateGatewayUser(db, gatewayUser)
}

func (gul *GatewayUserLogic) RemoveGatewayUser(id int) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.DeleteGatewayUser(db, id)
}

func (gul *GatewayUserLogic) GetGatewayUser(id int) (data.GatewayUser, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.GatewayUser{}, err
	}
	return data.GetGatewayUser(db, id)
}

func (gul *GatewayUserLogic) GetGatewayUserByPublicKeyDsa(pubklicKeyDSA string) (data.GatewayUser, error) {

	db, err := getDbConnection()
	if err != nil {
		return data.GatewayUser{}, err
	}
	return data.GetGatewayUserByPublicKeyDsa(db, pubklicKeyDSA)
}

func (gul *GatewayUserLogic) GetGatewayUserByPublicKeyKem(pubklicKeyKEM string) (data.GatewayUser, error) {

	db, err := getDbConnection()
	if err != nil {
		return data.GatewayUser{}, err
	}
	return data.GetGatewayUserByPublicKeyKem(db, pubklicKeyKEM)
}
