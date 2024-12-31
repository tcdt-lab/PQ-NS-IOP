package data

import "database/sql"

type Gateway struct {
	Id           int
	Ip           string
	Port         string
	PublicKey    string
	Ticket       string
	SymmetricKey string
}

func GetGateways(db *sql.DB) ([]Gateway, error) {

	rows, err := db.Query("SELECT * FROM gateways")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gateways := []Gateway{}

	for rows.Next() {
		var gateway Gateway
		if err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKey, &gateway.Ticket, &gateway.SymmetricKey); err != nil {
			return nil, err
		}
		gateways = append(gateways, gateway)
	}
	return gateways, nil
}

func AddGateway(db *sql.DB, gateway Gateway) (int64, error) {
	result, err := db.Exec("INSERT INTO gateways (Ip, Port, Public_Key, ticket, Symmetric_Key) VALUES (?, ?, ?, ?, ?)", gateway.Ip, gateway.Port, gateway.PublicKey, gateway.Ticket, gateway.SymmetricKey)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateGateway(db *sql.DB, gateway Gateway) (int64, error) {
	result, err := db.Exec("UPDATE gateways SET Ip = ?, Port = ?, Public_Key = ?, ticket = ?, Symmetric_Key = ? WHERE Id = ?", gateway.Ip, gateway.Port, gateway.PublicKey, gateway.Ticket, gateway.SymmetricKey, gateway.Id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func RemoveGateway(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM gateways WHERE Id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetGateway(db *sql.DB, id int) (Gateway, error) {
	var gateway Gateway
	rows := db.QueryRow("SELECT * FROM gateways WHERE Id = ?", id)
	err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKey, &gateway.Ticket, &gateway.SymmetricKey)
	if err != nil {
		return Gateway{}, err
	}
	return gateway, nil
}

func GetGatewayByIP(db *sql.DB, ip string) (Gateway, error) {

	var gateway Gateway
	rows := db.QueryRow("SELECT * FROM gateways WHERE Ip = ?", ip)
	err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKey, &gateway.Ticket, &gateway.SymmetricKey)
	if err != nil {
		return Gateway{}, err
	}
	return gateway, nil
}

func GetGatewayByIpAndPort(db *sql.DB, ip string, port string) (Gateway, error) {

	var gateway Gateway
	rows := db.QueryRow("SELECT * FROM gateways WHERE Ip = ? AND Port = ?", ip, port)
	err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKey, &gateway.Ticket, &gateway.SymmetricKey)
	if err != nil {
		return Gateway{}, err
	}
	return gateway, nil
}

func IfGatewayExist(db *sql.DB, gateway Gateway) (bool, error) {
	rows, err := db.Query("SELECT * FROM gateways WHERE Ip = ? AND Port = ?", gateway.Ip, gateway.Port)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}
