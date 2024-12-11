package data

import "database/sql"

type VerifierUser struct {
	Id        int
	Salt      string
	Password  string
	SecretKey string
	PublicKey string
}

func AddVerifierUser(v *VerifierUser, db *sql.DB) (int64, error) {

	result, err := db.Exec("INSERT INTO verifier_user (Salt, Password, Secret_Key, Public_Key) VALUES (?, ?, ?, ?)", v.Salt, v.Password, v.SecretKey, v.PublicKey)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetVerifierUsers(db *sql.DB) ([]VerifierUser, error) {
	rows, err := db.Query("SELECT * FROM verifier_user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	verifierUsers := []VerifierUser{}
	for rows.Next() {
		var verifierUser VerifierUser
		if err := rows.Scan(&verifierUser.Id, &verifierUser.Salt, &verifierUser.Password, &verifierUser.SecretKey, &verifierUser.PublicKey); err != nil {
			return nil, err
		}
		verifierUsers = append(verifierUsers, verifierUser)
	}
	return verifierUsers, nil
}

func GetVerifierUser(db *sql.DB, id int) (VerifierUser, error) {
	var verifierUser VerifierUser
	rows := db.QueryRow("SELECT * FROM verifier_user WHERE Id = ?", id)
	if err := rows.Scan(&verifierUser.Id, &verifierUser.Salt, &verifierUser.Password, &verifierUser.SecretKey, &verifierUser.PublicKey); err != nil {
		return verifierUser, err
	}
	return verifierUser, nil
}

func UpdateVerifierUser(v *VerifierUser, db *sql.DB) (int64, error) {
	result, err := db.Exec("UPDATE verifier_user SET Salt = ?, Password = ?, Secret_Key = ?, Public_Key = ? WHERE Id = ?", v.Salt, v.Password, v.SecretKey, v.PublicKey, v.Id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func RemoveVerifierUser(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM verifier_user WHERE Id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetVerifierUserByPublicKey(db *sql.DB, publicKey string) (VerifierUser, error) {
	var verifierUser VerifierUser
	rows := db.QueryRow("SELECT * FROM verifier_user WHERE Public_Key = ?", publicKey)
	if err := rows.Scan(&verifierUser.Id, &verifierUser.Salt, &verifierUser.Password, &verifierUser.SecretKey, &verifierUser.PublicKey); err != nil {
		return verifierUser, err
	}
	return verifierUser, nil
}

func GetVerifierUserBySecretKey(db *sql.DB, secretKey string) (VerifierUser, error) {
	var verifierUser VerifierUser
	rows := db.QueryRow("SELECT * FROM verifier_user WHERE Secret_Key = ?", secretKey)
	if err := rows.Scan(&verifierUser.Id, &verifierUser.Salt, &verifierUser.Password, &verifierUser.SecretKey, &verifierUser.PublicKey); err != nil {
		return verifierUser, err
	}
	return verifierUser, nil
}

func GetVerifierUserByPassword(db *sql.DB, password string) (VerifierUser, error) {
	var verifierUser VerifierUser
	rows := db.QueryRow("SELECT * FROM verifier_user WHERE Password = ?", password)
	if err := rows.Scan(&verifierUser.Id, &verifierUser.Salt, &verifierUser.Password, &verifierUser.SecretKey, &verifierUser.PublicKey); err != nil {
		return verifierUser, err
	}
	return verifierUser, nil
}
