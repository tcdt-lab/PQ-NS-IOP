package data

import "database/sql"

type GatewayUser struct {
	Id           int
	Salt         string
	Password     string
	PublicKeyDsa string
	SecretKeyDsa string
	PublicKeyKem string
	SecretKeyKem string
	Dsa_scheme   string
	Kem_scheme   string
}

func AddNewGatewayUser(db *sql.DB, user GatewayUser) (int64, error) {
	result, err := db.Exec("INSERT INTO gateway_user (salt, password, public_key_dsa, secret_key_dsa, public_key_kem, secret_key_kem, dsa_scheme, kem_scheme) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", user.Salt, user.Password, user.PublicKeyDsa, user.SecretKeyDsa, user.PublicKeyKem, user.SecretKeyKem, user.Dsa_scheme, user.Kem_scheme)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateGatewayUser(db *sql.DB, user GatewayUser) (int64, error) {
	result, err := db.Exec("UPDATE gateway_user SET salt = ?, password = ?, public_key_dsa = ?, secret_key_dsa = ?, public_key_kem = ?, secret_key_kem = ?, dsa_scheme = ?, kem_scheme = ? WHERE id = ?", user.Salt, user.Password, user.PublicKeyDsa, user.SecretKeyDsa, user.PublicKeyKem, user.SecretKeyKem, user.Dsa_scheme, user.Kem_scheme, user.Id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteGatewayUser(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM gateway_user WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetGatewayUser(db *sql.DB, id int) (GatewayUser, error) {
	var user GatewayUser
	rows := db.QueryRow("SELECT * FROM gateway_user WHERE id = ?", id)
	err := rows.Scan(&user.Id, &user.Salt, &user.Password, &user.PublicKeyDsa, &user.SecretKeyDsa, &user.PublicKeyKem, &user.SecretKeyKem, &user.Dsa_scheme, &user.Kem_scheme)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetAllGatewayUsers(db *sql.DB) ([]GatewayUser, error) {
	rows, err := db.Query("SELECT * FROM gateway_user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []GatewayUser{}

	for rows.Next() {
		var user GatewayUser
		if err := rows.Scan(&user.Id, &user.Salt, &user.Password, &user.PublicKeyDsa, &user.SecretKeyDsa, &user.PublicKeyKem, &user.SecretKeyKem, &user.Dsa_scheme, &user.Kem_scheme); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetGatewayUserByPassword(db *sql.DB, password string) (GatewayUser, error) {
	var user GatewayUser
	err := db.QueryRow("SELECT * FROM gateway_user WHERE password = ?", password).Scan(&user.Id, &user.Salt, &user.Password, &user.PublicKeyDsa, &user.SecretKeyDsa, &user.PublicKeyKem, &user.SecretKeyKem, &user.Dsa_scheme, &user.Kem_scheme)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetGatewayUserByPublicKeyDsa(db *sql.DB, publicKeyDsa string) (GatewayUser, error) {
	var user GatewayUser
	err := db.QueryRow("SELECT * FROM gateway_user WHERE public_key_dsa = ?", publicKeyDsa).Scan(&user.Id, &user.Salt, &user.Password, &user.PublicKeyDsa, &user.SecretKeyDsa, &user.PublicKeyKem, &user.SecretKeyKem, &user.Dsa_scheme, &user.Kem_scheme)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetGatewayUserByPublicKeyKem(db *sql.DB, publicKeyKem string) (GatewayUser, error) {
	var user GatewayUser
	err := db.QueryRow("SELECT * FROM gateway_user WHERE public_key_kem = ?", publicKeyKem).Scan(&user.Id, &user.Salt, &user.Password, &user.PublicKeyDsa, &user.SecretKeyDsa, &user.PublicKeyKem, &user.SecretKeyKem, &user.Dsa_scheme, &user.Kem_scheme)
	if err != nil {
		return user, err
	}
	return user, nil
}
