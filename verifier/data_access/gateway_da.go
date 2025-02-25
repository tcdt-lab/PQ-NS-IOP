package data_access

import (
	"database/sql"
	"verifier/data"
)

type GatewayDA struct {
	db *sql.DB
}

func GenerateGatewayDA(database *sql.DB) GatewayDA {
	var gda = GatewayDA{db: database}
	return gda
}

func (gDa *GatewayDA) GetGateways() ([]data.Gateway, error) {

	return data.GetGateways(gDa.db)
}

func (gDa *GatewayDA) AddGateway(gateway data.Gateway) (int64, error) {

	return data.AddGateway(gDa.db, gateway)
}

func (gDa *GatewayDA) UpdateGateway(gateway data.Gateway) (int64, error) {

	return data.UpdateGateway(gDa.db, gateway)
}

func (gDa *GatewayDA) RemoveGateway(id int) (int64, error) {

	return data.RemoveGateway(gDa.db, id)
}

func (gDa *GatewayDA) GetGateway(id int) (data.Gateway, error) {

	return data.GetGateway(gDa.db, id)
}

func (gDa *GatewayDA) GetGatewayByIpAndPort(ip string, port string) (data.Gateway, error) {

	return data.GetGatewayByIpAndPort(gDa.db, ip, port)
}

func (gDa *GatewayDA) GetGatewayByIP(ip string) (data.Gateway, error) {

	return data.GetGatewayByIp(gDa.db, ip)
}

func (gDa *GatewayDA) GetGatewayByPublicKeyKem(publicKeyKem string) (data.Gateway, error) {

	return data.GetGatewayByPublicKeyKem(gDa.db, publicKeyKem)
}

func (gDa *GatewayDA) GetGatewayByPublicKeySig(publicKeySig string) (data.Gateway, error) {

	return data.GetGatewayByPublicKeySig(gDa.db, publicKeySig)
}

func (gDa *GatewayDA) AddUpdateGateway(gateway data.Gateway) (int64, error) {

	if gatewayExist, _ := data.IfGatewayExists(gDa.db, gateway.Ip, gateway.Port); gatewayExist {
		return data.UpdateGateway(gDa.db, gateway)
	} else {
		return data.AddGateway(gDa.db, gateway)
	}
}

func (gda *GatewayDA) IfGatewayExist(ip string, port string) (bool, error) {

	return data.IfGatewayExists(gda.db, ip, port)
}

func (gda *GatewayDA) IfGatewayExistByPublicKeySig(publicKeySig string) (bool, error) {

	return data.IfGatewayExistPublicKeySig(gda.db, publicKeySig)
}
