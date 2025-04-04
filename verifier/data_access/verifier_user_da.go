package data_access

import (
	"database/sql"
	"go.uber.org/zap"
	"verifier/data"
)

type VerifierUserDA struct {
	db *sql.DB
}

func GenerateVerifierUserDA(databse *sql.DB) VerifierUserDA {
	var vDa = VerifierUserDA{db: databse}
	return vDa
}
func (vuda *VerifierUserDA) CloseDbConnection() {
	vuda.db.Close()
}
func (vuda *VerifierUserDA) GetVerifierUsers() ([]data.VerifierUser, error) {

	return data.GetAllVerifierUsers(vuda.db)
}

func (vuda *VerifierUserDA) AddVerifierUser(vu data.VerifierUser) (int64, error) {

	return data.AddVerifierUser(&vu, vuda.db)
}

func (vuda *VerifierUserDA) UpdateVerifierUser(vu data.VerifierUser) (int64, error) {

	return data.UpdateVerifierUser(&vu, vuda.db)
}

func (vuda *VerifierUserDA) UpdateVerifierUserByPassword(vu data.VerifierUser) (int64, error) {

	return data.UpdateVerifeirUserByPassword(&vu, vuda.db)
}

func (vuda *VerifierUserDA) UpdateVerifierUserByPublicKeySig(vu data.VerifierUser) (int64, error) {

	return data.UpdateVerifierUSerByPublicKeySig(&vu, vuda.db)
}

func (vuda *VerifierUserDA) DeleteVerifierUser(id int) (int64, error) {

	return data.DeleteVerifierUser(vuda.db, id)
}

func (vuda *VerifierUserDA) GetVerifierUser(id int64) (data.VerifierUser, error) {

	return data.GetVerifierUserById(vuda.db, int64(id))
}

func (vuda *VerifierUserDA) GetVerifierUserByPublicKeySig(publicKeySig string) (data.VerifierUser, error) {

	return data.GetVerifierUserByPublicKeySig(vuda.db, publicKeySig)
}

func (vuda *VerifierUserDA) GetVerifierUserByPassword(password string) (data.VerifierUser, error) {

	return data.GetVerifierUserByPassword(vuda.db, password)
}
func (vuda *VerifierUserDA) AddUpdateVerifierUSer(vu data.VerifierUser) (int64, error) {

	if exist, _ := data.IsVerifierUserExist(vuda.db, 1); exist {
		return data.UpdateVerifierUser(&vu, vuda.db)
	} else {
		return data.AddVerifierUser(&vu, vuda.db)
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

	cacheHandler := GenerateCacheHandlerDA()
	adminId, err := cacheHandler.GetUserAdminId()
	if err != nil {
		zap.L().Error("Error while getting admin id", zap.Error(err))
		return data.VerifierUser{}, err
	}
	return data.GetVerifierUserById(vuda.db, adminId)
}
