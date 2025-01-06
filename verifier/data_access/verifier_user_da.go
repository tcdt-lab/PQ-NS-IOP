package data_access

import "verifier/data"

type VerifierUserDA struct {
}

func (vuda *VerifierUserDA) GetVerifierUsers() ([]data.VerifierUser, error) {
	db, err := getDbConnection()
	if err != nil {
		return nil, err
	}
	return data.GetAllVerifierUsers(db)
}

func (vuda *VerifierUserDA) AddVerifierUser(vu data.VerifierUser) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.AddVerifierUser(&vu, db)
}

func (vuda *VerifierUserDA) UpdateVerifierUser(vu data.VerifierUser) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.UpdateVerifierUser(&vu, db)
}

func (vuda *VerifierUserDA) UpdateVerifierUserByPassword(vu data.VerifierUser) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.UpdateVerifeirUserByPassword(&vu, db)
}

func (vuda *VerifierUserDA) UpdateVerifierUserByPublicKeySig(vu data.VerifierUser) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.UpdateVerifierUSerByPublicKeySig(&vu, db)
}

func (vuda *VerifierUserDA) DeleteVerifierUser(id int) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	return data.DeleteVerifierUser(db, id)
}

func (vuda *VerifierUserDA) GetVerifierUser(id int) (data.VerifierUser, error) {
	db, err := getDbConnection()
	if err != nil {
		return data.VerifierUser{}, err
	}
	return data.GetVerifierUserById(db, id)
}

func (vuda *VerifierUserDA) GetVerifierUserByPublicKeySig(publicKeySig string) (data.VerifierUser, error) {

	db, err := getDbConnection()
	if err != nil {
		return data.VerifierUser{}, err
	}
	return data.GetVerifierUserByPublicKeySig(db, publicKeySig)
}

func (vuda *VerifierUserDA) GetVerifierUserByPassword(password string) (data.VerifierUser, error) {

	db, err := getDbConnection()
	if err != nil {
		return data.VerifierUser{}, err
	}
	return data.GetVerifierUserByPassword(db, password)
}
func (vuda *VerifierUserDA) AddUpdateVerifierUSer(vu data.VerifierUser) (int64, error) {
	db, err := getDbConnection()
	if err != nil {
		return 0, err
	}
	if exist, _ := data.IsVerifierUserExist(db, 1); exist {
		return data.UpdateVerifierUser(&vu, db)
	} else {
		return data.AddVerifierUser(&vu, db)
	}
}
func (vuda *VerifierUserDA) SetUpAdminVerifierUser(publicKeyKem string, secKeyKem string, pubKeySig string, secKeySig string) (int64, error) {
	verifierUser := data.VerifierUser{}
	verifierUser.Id = 1
	verifierUser.PublicKeyKem = publicKeyKem
	verifierUser.SecretKeyKem = secKeyKem
	verifierUser.PublicKeySig = pubKeySig
	verifierUser.SecretKeySig = secKeySig
	return vuda.AddUpdateVerifierUSer(verifierUser)
}

func (vuda *VerifierUserDA) GetAdminVerifierUser() (data.VerifierUser, error) {

	db, err := getDbConnection()
	if err != nil {
		return data.VerifierUser{}, err
	}
	return data.GetVerifierUserById(db, 1)
}
