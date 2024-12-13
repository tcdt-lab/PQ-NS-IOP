package data

import "database/sql"

type VerifierUser struct {
	Id           int
	Salt         string
	Password     string
	SecretKeySig string
	PublicKeySig string
	SecretKeyKem string
	PublicKeyKem string
}

func AddVerifierUser(v *VerifierUser, db *sql.DB) (int64, error) {
	rowChange, err := db.Exec("INSERT INTO verifier_user (salt, password, secret_key_sig, public_key_sig, secret_key_kem, public_key_kem) VALUES(?, ?, ?, ?, ?, ?)", v.Salt, v.Password, v.SecretKeySig, v.PublicKeySig, v.SecretKeyKem, v.PublicKeyKem)
	if err != nil {
		return 0, err
	}
	return rowChange.LastInsertId()
}

func UpdateVerifierUser(v *VerifierUser, db *sql.DB) (int64, error) {
	rowChange, err := db.Exec("UPDATE verifier_user SET salt=?, password=?, secret_key_sig=?, public_key_sig=?, secret_key_kem=?, public_key_kem=? WHERE id=?", v.Salt, v.Password, v.SecretKeySig, v.PublicKeySig, v.SecretKeyKem, v.PublicKeyKem, v.Id)
	if err != nil {
		return 0, err
	}
	return rowChange.RowsAffected()
}

func UpdateVerifeirUserByPassword(v *VerifierUser, db *sql.DB) (int64, error) {
	rowChange, err := db.Exec("UPDATE verifier_user SET salt=?, secret_key_sig=?, public_key_sig=?, secret_key_kem=?, public_key_kem=? WHERE password=?", v.Salt, v.SecretKeySig, v.PublicKeySig, v.SecretKeyKem, v.PublicKeyKem, v.Password)
	if err != nil {
		return 0, err
	}
	return rowChange.RowsAffected()
}

func UpdateVerifierUSerByPublicKeySig(v *VerifierUser, db *sql.DB) (int64, error) {
	rowChange, err := db.Exec("UPDATE verifier_user SET salt=?, password=?, secret_key_sig=?, secret_key_kem=?, public_key_kem=? WHERE public_key_sig=?", v.Salt, v.Password, v.SecretKeySig, v.SecretKeyKem, v.PublicKeyKem, v.PublicKeySig)
	if err != nil {
		return 0, err
	}
	return rowChange.RowsAffected()
}

func DeleteVerifierUser(db *sql.DB, id int) (int64, error) {
	rowChange, err := db.Exec("DELETE FROM verifier_user WHERE id=?", id)
	if err != nil {
		return 0, err
	}
	return rowChange.RowsAffected()
}

func GetAllVerifierUsers(db *sql.DB) ([]VerifierUser, error) {
	rows, err := db.Query("SELECT * FROM verifier_user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []VerifierUser
	for rows.Next() {
		var v VerifierUser
		err := rows.Scan(&v.Id, &v.Salt, &v.Password, &v.SecretKeySig, &v.PublicKeySig, &v.SecretKeyKem, &v.PublicKeyKem)
		if err != nil {
			return nil, err
		}
		users = append(users, v)
	}
	return users, nil
}

func GetVerifierUserByPublicKeySig(db *sql.DB, publicKeySig string) (VerifierUser, error) {
	var v VerifierUser
	err := db.QueryRow("SELECT * FROM verifier_user WHERE public_key_sig = ?", publicKeySig).Scan(&v.Id, &v.Salt, &v.Password, &v.SecretKeySig, &v.PublicKeySig, &v.SecretKeyKem, &v.PublicKeyKem)
	return v, err
}

func GetVerifierUserByPassword(db *sql.DB, password string) (VerifierUser, error) {
	var v VerifierUser
	err := db.QueryRow("SELECT * FROM verifier_user WHERE password = ?", password).Scan(&v.Id, &v.Salt, &v.Password, &v.SecretKeySig, &v.PublicKeySig, &v.SecretKeyKem, &v.PublicKeyKem)
	return v, err
}
func GetVerifierUserById(db *sql.DB, id int) (VerifierUser, error) {
	var v VerifierUser
	err := db.QueryRow("SELECT * FROM verifier_user WHERE id = ?", id).Scan(&v.Id, &v.Salt, &v.Password, &v.SecretKeySig, &v.PublicKeySig, &v.SecretKeyKem, &v.PublicKeyKem)
	return v, err
}
