package data_access

import (
	"database/sql"
	"gateway/data"
)

type VerifierDA struct {
	db *sql.DB
}

func GenerateVerifierDA(db *sql.DB) *VerifierDA {
	return &VerifierDA{
		db: db,
	}
}

func (vl *VerifierDA) GetVerifiers() ([]data.Verifier, error) {

	return data.GetVerifiers(vl.db)
}

func (vl *VerifierDA) AddVerifier(verifier data.Verifier) (int64, error) {

	return data.AddVerifier(vl.db, verifier)
}

func (vl *VerifierDA) UpdateVerifier(verifier data.Verifier) (int64, error) {

	return data.UpdateVerifier(vl.db, verifier)
}

func (vl *VerifierDA) RemoveVerifier(id int) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.RemoveVerifier(db, id)
}

func (vl *VerifierDA) GetVerifier(id int) (data.Verifier, error) {

	return data.GetVerifier(vl.db, id)
}

func (vl *VerifierDA) GetVerifierByIpAndPort(ip string, port string) (data.Verifier, error) {

	return data.GetVerifierByIpandPort(vl.db, ip, port)
}

func (vl *VerifierDA) GetVerifierByPublicKey(publicKey string) (data.Verifier, error) {

	return data.GetVerifierByPublicKey(vl.db, publicKey)
}
func (vl *VerifierDA) GetVerifierByIP(ip string) (data.Verifier, error) {

	return data.GetVerifierByIP(vl.db, ip)
}

func (vl *VerifierDA) IfVerifierExistsByIpandPort(verifier data.Verifier) (bool, error) {

	return data.IfVerifierExistsWithIPandPort(vl.db, verifier)
}

func (vl *VerifierDA) AddUpdateVerifiers(verifier []data.Verifier) error {

	for _, v := range verifier {
		if exist, _ := data.IfVerifierExistsWithIPandPort(vl.db, v); exist {
			if _, err := data.UpdateVerifierWithIpandPort(vl.db, v); err != nil {
				return err
			}
		} else {
			if _, err := data.AddVerifier(vl.db, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (vl *VerifierDA) AddUpdateVerifier(verifier data.Verifier) error {
	if exist, _ := data.IfVerifierExistsWithIPandPort(vl.db, verifier); exist {
		if _, err := data.UpdateVerifierWithIpandPort(vl.db, verifier); err != nil {
			return err
		}
	} else {
		if _, err := data.AddVerifier(vl.db, verifier); err != nil {
			return err
		}
	}
	return nil
}

func (vl *VerifierDA) IfVerifierExistsBuPublicKeySig(publicKeySig string) (bool, error) {
	return data.IfVerifierExistByPubKeySign(vl.db, publicKeySig)
}

func (vl *VerifierDA) GetVerifierByPublicKeySig(publicKey string) (data.Verifier, error) {
	return data.GetVerifierByPublicKey(vl.db, publicKey)
}
