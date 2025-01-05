package data_access

import "verifier/data"

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
	return data.AddNewVerifier(&verifier, db)
}

func (vl *VerifierDA) UpdateVerifier(verifier data.Verifier) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.UpdateVerifier(&verifier, db)
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
	return data.GetVerifierByIpAndPort(db, ip, port)
}

func (vl *VerifierDA) GetVerifierByPublicKeySig(publicKey string) (data.Verifier, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Verifier{}, err
	}
	return data.GetVerifierByPublicKeySig(db, publicKey)
}

func (vl *VerifierDA) GetVerifierByIP(ip string) (data.Verifier, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.Verifier{}, err
	}
	return data.GetVerifierByIp(db, ip)
}

func (vl *VerifierDA) AddUpdateVerifier(verifier data.Verifier) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	if verifierExist, _ := data.IfVerifierExists(db, verifier); verifierExist {
		return data.UpdateVerifier(&verifier, db)
	}
	return data.AddNewVerifier(&verifier, db)
}
