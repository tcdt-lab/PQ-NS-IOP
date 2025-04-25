package data

import "database/sql"

type Verifier struct {
	Id            int
	Ip            string
	Port          string
	PublicKeySig  string
	PublicKeyKem  string
	SigScheme     string
	SymmetricKey  string
	TrustScore    float64
	IsInCommittee bool
}

func AddNewVerifier(v *Verifier, db *sql.DB) (int64, error) {

	result, err := db.Exec("INSERT INTO verifiers (Ip, Port, Public_Key_Sig,Public_Key_Kem, Sig_Scheme, Symmetric_Key, Trust_Score, Is_In_Committee) VALUES (?, ?, ?, ?, ?, ?, ?,?)", v.Ip, v.Port, v.PublicKeySig, v.PublicKeyKem, v.SigScheme, v.SymmetricKey, v.TrustScore, v.IsInCommittee)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
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
		if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKeySig, &verifier.PublicKeyKem, &verifier.SigScheme, &verifier.SymmetricKey, &verifier.TrustScore, &verifier.IsInCommittee); err != nil {
			return nil, err
		}
		verifiers = append(verifiers, verifier)
	}
	return verifiers, nil
}

func GetVerifier(db *sql.DB, id int64) (Verifier, error) {
	var verifier Verifier
	rows := db.QueryRow("SELECT * FROM verifiers WHERE Id = ?", id)
	if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKeySig, &verifier.PublicKeyKem, &verifier.SigScheme, &verifier.SymmetricKey, &verifier.TrustScore, &verifier.IsInCommittee); err != nil {
		return verifier, err
	}
	return verifier, nil
}

func UpdateVerifier(v *Verifier, db *sql.DB) (int64, error) {
	result, err := db.Exec("UPDATE verifiers SET Ip = ?, Port = ?,Public_Key_Sig=?,Public_Key_Kem =?, Sig_Scheme = ?, Symmetric_Key = ?, Trust_Score = ?, Is_In_Committee = ? WHERE Public_Key_Sig = ?", v.Ip, v.Port, v.PublicKeySig, v.PublicKeyKem, v.SigScheme, v.SymmetricKey, v.TrustScore, v.IsInCommittee, v.PublicKeySig)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func UpdateVerifierByIpandPort(v *Verifier, db *sql.DB) (int64, error) {
	result, err := db.Exec("UPDATE verifiers SET Ip = ?, Port = ?,Public_Key_Sig=?,Public_Key_Kem =?, Sig_Scheme = ?, Symmetric_Key = ?, Trust_Score = ?, Is_In_Committee = ? WHERE Ip = ? AND Port = ?", v.Ip, v.Port, v.PublicKeySig, v.PublicKeyKem, v.SigScheme, v.SymmetricKey, v.TrustScore, v.IsInCommittee, v.Ip, v.Port)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()

}

func RemoveVerifier(db *sql.DB, id int64) (int64, error) {
	result, err := db.Exec("DELETE FROM verifiers WHERE Id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetVerifierByPublicKeySig(db *sql.DB, publicKeySig string) (Verifier, error) {
	var verifier Verifier
	rows := db.QueryRow("SELECT * FROM verifiers WHERE Public_Key_Sig = ?", publicKeySig)
	if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKeySig, &verifier.PublicKeyKem, &verifier.SigScheme, &verifier.SymmetricKey, &verifier.TrustScore, &verifier.IsInCommittee); err != nil {
		return verifier, err
	}
	return verifier, nil
}

func GetVerifierByIp(db *sql.DB, ip string) (Verifier, error) {
	var verifier Verifier
	rows := db.QueryRow("SELECT * FROM verifiers WHERE Ip = ?", ip)
	if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKeySig, &verifier.PublicKeyKem, &verifier.SigScheme, &verifier.SymmetricKey, &verifier.TrustScore, &verifier.IsInCommittee); err != nil {
		return verifier, err
	}
	return verifier, nil
}

func GetVerifierByIpAndPort(db *sql.DB, ip string, port string) (Verifier, error) {
	var verifier Verifier
	rows := db.QueryRow("SELECT * FROM verifiers WHERE Ip = ? AND Port = ?", ip, port)
	if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKeySig, &verifier.PublicKeyKem, &verifier.SigScheme, &verifier.SymmetricKey, &verifier.TrustScore, &verifier.IsInCommittee); err != nil {
		return verifier, err
	}
	return verifier, nil
}

func GetVerifiersInCommittee(db *sql.DB) ([]Verifier, error) {
	rows, err := db.Query("SELECT * FROM verifiers WHERE Is_In_Committee = 1")
	var verifiers []Verifier
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var verifier Verifier
		if err := rows.Scan(&verifier.Id, &verifier.Ip, &verifier.Port, &verifier.PublicKeySig, &verifier.PublicKeyKem, &verifier.SigScheme, &verifier.SymmetricKey, &verifier.TrustScore, &verifier.IsInCommittee); err != nil {
			return nil, err
		}
		verifiers = append(verifiers, verifier)
	}

	return verifiers, nil
}

func IfVerifierExists(db *sql.DB, verifier Verifier) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM verifiers WHERE Ip = ? AND Port = ?", verifier.Ip, verifier.Port).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func IfVerifierExistByPublicKeySig(db *sql.DB, pubKey string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM verifiers WHERE Public_Key_Sig = ?", pubKey).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
