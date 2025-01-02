package data_access

import (
	"gateway/data"
)

type GatewayDA struct {
}

func (gl *GatewayDA) GetGateways() ([]data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return nil, err
	}
	return data.GetGateways(db)
}

func (gl *GatewayDA) AddGateway(gateway data.Gateway) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.AddGateway(db, gateway)
}

func (gl *GatewayDA) UpdateGateway(gateway data.Gateway) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.UpdateGateway(db, gateway)
}

func (gl *GatewayDA) RemoveGateway(id int) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.RemoveGateway(db, id)
}

func (gl *GatewayDA) GetGateway(id int) (data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Gateway{}, err
	}
	return data.GetGateway(db, id)
}

func (gl *GatewayDA) GetGatewayByIpAndPort(ip string, port string) (data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Gateway{}, err
	}
	return data.GetGatewayByIpAndPort(db, ip, port)
}

func (gl *GatewayDA) GetGatewayByIP(ip string) (data.Gateway, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Gateway{}, err
	}
	return data.GetGatewayByIP(db, ip)
}

func (gl *GatewayDA) IfGatewayExist(gateway data.Gateway) (bool, error) {
	db, err := getDbConnection()
	if err != nil {
		return false, err
	}
	return data.IfGatewayExist(db, gateway)
}

func (gl *GatewayDA) AddUpdateGateways(gateways []data.Gateway) error {
	db, err := getDbConnection()
	if err != nil {
		return err
	}
	for _, gateway := range gateways {
		if exist, _ := data.IfGatewayExist(db, gateway); exist {
			_, err = data.UpdateGateway(db, gateway)
			if err != nil {
				return err
			}
		} else {
			_, err = data.AddGateway(db, gateway)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
