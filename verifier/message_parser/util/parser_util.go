package util

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
	"test.org/cryptography/pkg/asymmetric"
	"test.org/cryptography/pkg/symmetric"
	"test.org/protocol/pkg"
	"verifier/config"
	"verifier/data"
)

func GetDBConnection(c config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ProtocolUtilGenerator(cryptographySchemeName string) pkg.ProtocolUtil {
	var util pkg.ProtocolUtil
	util.AesHandler = symmetric.AesGcm{}
	util.AsymmetricHandler = asymmetric.NewAsymmetricHandler(cryptographySchemeName)
	util.HmacHandler = symmetric.HMAC{}
	util.RegisterInterfacesInGob()
	return util
}

func GetUserInformation(db *sql.DB, c config.Config) (data.VerifierUser, error) {
	verifeirUser, err := data.GetVerifierUserByPassword(db, os.Getenv("PQ_NS_IOP_VU_PASS"))
	if err != nil {
		zap.L().Error("ErrorParams while getting verifier_verifier user", zap.Error(err))
		return data.VerifierUser{}, err
	}
	return verifeirUser, nil

}

func GetSenderKeys(senderIp string, c config.Config) (string, string, error) {
	db, err := GetDBConnection(c)
	if err != nil {
		return "", "", err
	}
	var gateway data.Gateway
	var verifier data.Verifier

	gateway, err = data.GetGatewayByIp(db, senderIp)
	if err == nil {
		return gateway.SymmetricKey, gateway.PublicKeySig, nil
	}
	verifier, err = data.GetVerifierByIp(db, senderIp)
	if err != nil {
		return "", "", err
	}
	return verifier.SymmetricKey, verifier.PublicKeySig, nil
}
