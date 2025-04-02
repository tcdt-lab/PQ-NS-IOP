package util

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
	"test.org/cryptography/pkg/asymmetric"
	"test.org/cryptography/pkg/symmetric"
	"test.org/protocol/pkg"
	"verifier/config"
	"verifier/data"
)

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

func GenerateRequestNumber() (int64, error) {
	var requestNumber int64
	err := binary.Read(rand.Reader, binary.LittleEndian, &requestNumber)
	if err != nil {
		return 0, err
	}
	return requestNumber, nil
}
