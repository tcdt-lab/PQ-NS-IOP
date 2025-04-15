package data_access

import (
	"database/sql"
	"gateway/data"
)

type GatewayDA struct {
	db *sql.DB
}

func GenerateGatewayDA(db *sql.DB) *GatewayDA {
	return &GatewayDA{
		db: db,
	}
}
func (gl *GatewayDA) GetGateways() ([]data.Gateway, error) {

	return data.GetGateways(gl.db)
}

func (gl *GatewayDA) AddGateway(gateway data.Gateway) (int64, error) {

	return data.AddGateway(gl.db, gateway)
}

func (gl *GatewayDA) UpdateGateway(gateway data.Gateway) (int64, error) {

	return data.UpdateGateway(gl.db, gateway)
}

func (gl *GatewayDA) RemoveGateway(id int) (int64, error) {

	return data.RemoveGateway(gl.db, id)
}

func (gl *GatewayDA) GetGateway(id int) (data.Gateway, error) {

	return data.GetGateway(gl.db, id)
}

func (gl *GatewayDA) GetGatewayByIpAndPort(ip string, port string) (data.Gateway, error) {

	return data.GetGatewayByIpAndPort(gl.db, ip, port)
}

func (gl *GatewayDA) GetGatewayByIP(ip string) (data.Gateway, error) {

	return data.GetGatewayByIP(gl.db, ip)
}

func (gl *GatewayDA) IfGatewayExist(gateway data.Gateway) (bool, error) {

	return data.IfGatewayExist(gl.db, gateway)
}

func (gl *GatewayDA) AddUpdateGateways(gateways []data.Gateway) error {

	for _, gateway := range gateways {
		if exist, _ := data.IfGatewayExist(gl.db, gateway); exist {
			_, err := data.UpdateGateway(gl.db, gateway)
			if err != nil {
				return err
			}
		} else {
			_, err := data.AddGateway(gl.db, gateway)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gl *GatewayDA) IfGatewayExistByPublicKeySig(publicKeySig string) (bool, error) {
	return data.IfGatewayExistByPubicKeySig(gl.db, publicKeySig)
}

func (gl *GatewayDA) GetGatewayByPublicKey(publicKey string) (data.Gateway, error) {
	return data.GetGatewayByPublicKey(gl.db, publicKey)
}
