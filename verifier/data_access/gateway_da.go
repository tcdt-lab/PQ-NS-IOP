package data_access

import "verifier/data"

type GatewayDA struct {
}

func (gDa *GatewayDA) GetGateways() ([]data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return nil, err
	}
	return data.GetGateways(db)
}

func (gDa *GatewayDA) AddGateway(gateway data.Gateway) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.AddGateway(db, gateway)
}

func (gDa *GatewayDA) UpdateGateway(gateway data.Gateway) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.UpdateGateway(db, gateway)
}

func (gDa *GatewayDA) RemoveGateway(id int) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.RemoveGateway(db, id)
}

func (gDa *GatewayDA) GetGateway(id int) (data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Gateway{}, err
	}
	return data.GetGateway(db, id)
}

func (gDa *GatewayDA) GetGatewayByIpAndPort(ip string, port string) (data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Gateway{}, err
	}
	return data.GetGatewayByIpAndPort(db, ip, port)
}

func (gDa *GatewayDA) GetGatewayByIP(ip string) (data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Gateway{}, err
	}
	return data.GetGatewayByIp(db, ip)
}

func (gDa *GatewayDA) GetGatewayByPublicKeyKem(publicKeyKem string) (data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Gateway{}, err
	}
	return data.GetGatewayByPublicKeyKem(db, publicKeyKem)
}

func (gDa *GatewayDA) GetGatewayByPublicKeySig(publicKeySig string) (data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Gateway{}, err
	}
	return data.GetGatewayByPublicKeySig(db, publicKeySig)
}

func (gDa *GatewayDA) AddUpdateGateway(gateway data.Gateway) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	if gatewayExist, _ := data.IfGatewayExists(db, gateway); gatewayExist {
		return data.UpdateGateway(db, gateway)
	} else {
		return data.AddGateway(db, gateway)
	}
}
