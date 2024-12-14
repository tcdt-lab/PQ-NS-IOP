package data

import "database/sql"

type Verifier struct {
	Id           int
	Ip           string
	Port         string
	PublicKey    string
	SymmetricKey string
}

func GetVerifiers(db *sql.DB) ([]Verifier, error) {
	rows, err := db.Query("SELECT * FROM verifiers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	verifiers := []Verifier{}

	for rows.Next() {
		var verifier Verifier
		if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKey, &verifier.SymmetricKey); err != nil {
			return nil, err
		}
		verifiers = append(verifiers, verifier)
	}
	return verifiers, nil
}

func AddVerifier(db *sql.DB, verifier Verifier) (int64, error) {
	result, err := db.Exec("INSERT INTO verifiers (Ip, Port, public_key, Symmetric_Key) VALUES (?, ?, ?, ?)", verifier.Ip, verifier.Port, verifier.PublicKey, verifier.SymmetricKey)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateVerifier(db *sql.DB, verifier Verifier) (int64, error) {
	result, err := db.Exec("UPDATE verifiers SET Ip = ?, Port = ?, public_key = ?, symmetric_key = ? WHERE Id = ?", verifier.Ip, verifier.Port, verifier.PublicKey, verifier.SymmetricKey, verifier.Id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func RemoveVerifier(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM verifiers WHERE Id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetVerifier(db *sql.DB, id int) (Verifier, error) {
	var verifier Verifier
	rows := db.QueryRow("SELECT * FROM verifiers WHERE Id = ?", id)
	if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKey, &verifier.SymmetricKey); err != nil {
		return Verifier{}, err
	}
	return verifier, nil
}

func GetVerifierByIP(db *sql.DB, ip string) (Verifier, error) {
	var verifier Verifier
	rows := db.QueryRow("SELECT * FROM verifiers WHERE Ip = ?", ip)
	if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKey, &verifier.SymmetricKey); err != nil {
		return Verifier{}, err
	}
	return verifier, nil
}

func GetVerifierByPublicKey(db *sql.DB, publicKey string) (Verifier, error) {
	var verifier Verifier
	rows := db.QueryRow("SELECT * FROM verifiers WHERE public_key = ?", publicKey)
	if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKey, &verifier.SymmetricKey); err != nil {
		return Verifier{}, err
	}
	return verifier, nil
}

func GetVerifierByIpandPort(db *sql.DB, ip string, port string) (Verifier, error) {
	var verifier Verifier
	rows := db.QueryRow("SELECT * FROM verifiers WHERE Ip = ? AND Port = ?", ip, port)
	if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKey, &verifier.SymmetricKey); err != nil {
		return Verifier{}, err
	}
	return verifier, nil
}
