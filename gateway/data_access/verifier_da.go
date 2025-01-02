package data_access

import "gateway/data"

type VerifierDA struct {
}

func (vl *VerifierDA) GetVerifiers() ([]data.Verifier, error) {
	db, err := getDbConnection()
	if err != nil {
		return nil, err
	}
	return data.GetVerifiers(db)
}

func (vl *VerifierDA) AddVerifier(verifier data.Verifier) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.AddVerifier(db, verifier)
}

func (vl *VerifierDA) UpdateVerifier(verifier data.Verifier) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.UpdateVerifier(db, verifier)
}

func (vl *VerifierDA) RemoveVerifier(id int) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.RemoveVerifier(db, id)
}

func (vl *VerifierDA) GetVerifier(id int) (data.Verifier, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Verifier{}, err
	}
	return data.GetVerifier(db, id)
}

func (vl *VerifierDA) GetVerifierByIpAndPort(ip string, port string) (data.Verifier, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Verifier{}, err
	}
	return data.GetVerifierByIpandPort(db, ip, port)
}

func (vl *VerifierDA) GetVerifierByPublicKey(publicKey string) (data.Verifier, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Verifier{}, err
	}
	return data.GetVerifierByPublicKey(db, publicKey)
}
func (vl *VerifierDA) GetVerifierByIP(ip string) (data.Verifier, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Verifier{}, err
	}
	return data.GetVerifierByIP(db, ip)
}

func (vl *VerifierDA) IfVerifierExists(verifier data.Verifier) (bool, error) {
	db, err := getDbConnection()
	if err != nil {
		return false, err
	}
	return data.IfVerifierExists(db, verifier)
}

func (vl *VerifierDA) AddUpdateVerifiers(verifier []data.Verifier) error {
	db, err := getDbConnection()
	if err != nil {
		return err
	}
	for _, v := range verifier {
		if exist, _ := data.IfVerifierExists(db, v); exist {
			if _, err := data.UpdateVerifier(db, v); err != nil {
				return err
			}
		} else {
			if _, err := data.AddVerifier(db, v); err != nil {
				return err
			}
		}
	}
	return nil
}
