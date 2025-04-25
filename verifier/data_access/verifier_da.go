package data_access

import (
	"database/sql"
	"verifier/data"
)

type VerifierDA struct {
	db *sql.DB
}

func GenerateVerifierDA(database *sql.DB) VerifierDA {
	var vda = VerifierDA{db: database}
	return vda
}
func (vl *VerifierDA) GetVerifiers() ([]data.Verifier, error) {

	return data.GetVerifiers(vl.db)
}

func (vl *VerifierDA) AddVerifier(verifier data.Verifier) (int64, error) {

	return data.AddNewVerifier(&verifier, vl.db)
}

func (vl *VerifierDA) UpdateVerifier(verifier data.Verifier) (int64, error) {

	return data.UpdateVerifier(&verifier, vl.db)
}

func (vl *VerifierDA) RemoveVerifier(id int64) (int64, error) {

	return data.RemoveVerifier(vl.db, id)
}

func (vl *VerifierDA) GetVerifier(id int64) (data.Verifier, error) {

	return data.GetVerifier(vl.db, id)
}

func (vl *VerifierDA) GetVerifierByIpAndPort(ip string, port string) (data.Verifier, error) {

	return data.GetVerifierByIpAndPort(vl.db, ip, port)
}

func (vl *VerifierDA) GetVerifierByPublicKeySig(publicKey string) (data.Verifier, error) {

	return data.GetVerifierByPublicKeySig(vl.db, publicKey)
}

func (vl *VerifierDA) GetVerifierByIP(ip string) (data.Verifier, error) {

	return data.GetVerifierByIp(vl.db, ip)
}

func (vl *VerifierDA) AddUpdateVerifier(verifier data.Verifier) (int64, error) {

	if verifierExist, _ := data.IfVerifierExists(vl.db, verifier); verifierExist {
		return data.UpdateVerifierByIpandPort(&verifier, vl.db)
	}
	return data.AddNewVerifier(&verifier, vl.db)
}

func (vl *VerifierDA) AddUpdateVerifiers(verifiers []data.Verifier) (int64, error) {
	var err error
	var rowsAffected int64
	for _, verifier := range verifiers {
		if rowsAffected, err = vl.AddUpdateVerifier(verifier); err != nil {
			return 0, err
		}
	}
	return rowsAffected, nil
}

func (vl *VerifierDA) IfVerifierExistByPublicKeySig(pubKey string) (bool, error) {
	return data.IfVerifierExistByPublicKeySig(vl.db, pubKey)
}
