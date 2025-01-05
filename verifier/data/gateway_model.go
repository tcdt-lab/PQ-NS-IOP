package data

import "database/sql"

type Gateway struct {
	Id           int
	Ip           string
	Port         string
	PublicKeyKem string
	PublicKeySig string
	KemScheme    string
	SigScheme    string
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
		if err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKeyKem, &gateway.PublicKeySig, &gateway.KemScheme, &gateway.SigScheme, &gateway.Ticket, &gateway.SymmetricKey); err != nil {
			return nil, err
		}
		gateways = append(gateways, gateway)
	}
	return gateways, nil
}

func AddGateway(db *sql.DB, gateway Gateway) (int64, error) {
	result, err := db.Exec("INSERT INTO gateways (Ip, Port, Public_Key_Kem, Public_Key_Sig, Kem_Scheme, Sig_Scheme, ticket, Symmetric_Key) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", gateway.Ip, gateway.Port, gateway.PublicKeyKem, gateway.PublicKeySig, gateway.KemScheme, gateway.SigScheme, gateway.Ticket, gateway.SymmetricKey)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateGateway(db *sql.DB, gateway Gateway) (int64, error) {
	result, err := db.Exec("UPDATE gateways SET Ip = ?, Port = ?, Public_Key_Kem = ?, Public_Key_Sig = ?, Kem_Scheme = ?, Sig_Scheme = ?, ticket = ?, Symmetric_Key = ? WHERE Id = ?", gateway.Ip, gateway.Port, gateway.PublicKeyKem, gateway.PublicKeySig, gateway.KemScheme, gateway.SigScheme, gateway.Ticket, gateway.SymmetricKey, gateway.Id)
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
	err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKeyKem, &gateway.PublicKeySig, &gateway.KemScheme, &gateway.SigScheme, &gateway.Ticket, &gateway.SymmetricKey)
	if err != nil {
		return Gateway{}, err
	}
	return gateway, nil
}

func GetGatewayByPublicKeySig(db *sql.DB, publicKeySig string) (Gateway, error) {

	var gateway Gateway
	rows := db.QueryRow("SELECT * FROM gateways WHERE Public_Key_Sig = ?", publicKeySig)
	err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKeyKem, &gateway.PublicKeySig, &gateway.KemScheme, &gateway.SigScheme, &gateway.Ticket, &gateway.SymmetricKey)
	if err != nil {
		return Gateway{}, err
	}
	return gateway, nil
}

func GetGatewayByPublicKeyKem(db *sql.DB, publicKeyKem string) (Gateway, error) {

	var gateway Gateway
	rows := db.QueryRow("SELECT * FROM gateways WHERE Public_Key_Kem = ?", publicKeyKem)
	err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKeyKem, &gateway.PublicKeySig, &gateway.KemScheme, &gateway.SigScheme, &gateway.Ticket, &gateway.SymmetricKey)
	if err != nil {
		return Gateway{}, err
	}
	return gateway, nil
}

func GetGatewayByIp(db *sql.DB, ip string) (Gateway, error) {

	var gateway Gateway
	rows := db.QueryRow("SELECT * FROM gateways WHERE Ip = ?", ip)
	err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKeyKem, &gateway.PublicKeySig, &gateway.KemScheme, &gateway.SigScheme, &gateway.Ticket, &gateway.SymmetricKey)
	if err != nil {
		return Gateway{}, err
	}
	return gateway, nil
}

func GetGatewayByIpAndPort(db *sql.DB, ip string, port string) (Gateway, error) {

	var gateway Gateway
	rows := db.QueryRow("SELECT * FROM gateways WHERE Ip = ? AND Port = ?", ip, port)
	err := rows.Scan(&gateway.Id, &gateway.Ip, &gateway.Port, &gateway.PublicKeyKem, &gateway.PublicKeySig, &gateway.KemScheme, &gateway.SigScheme, &gateway.Ticket, &gateway.SymmetricKey)
	if err != nil {
		return Gateway{}, err
	}
	return gateway, nil
}

func IfGatewayExists(db *sql.DB, gateway Gateway) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM gateways WHERE Ip = ? AND Port = ?", gateway.Ip, gateway.Port).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
