package util

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"gateway/config"
	_ "github.com/go-sql-driver/mysql"
	"test.org/cryptography/pkg/asymmetric"
	"test.org/cryptography/pkg/symmetric"
	"test.org/protocol/pkg"
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
	util.PBKDF2Handler = symmetric.PBKDF2{}
	util.RegisterInterfacesInGob()
	return util
}

func GenerateNonce() (string, error) {
	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(nonce), nil

}

func GenerateRequestNumber() (int64, error) {
	var requestNumber int64
	err := binary.Read(rand.Reader, binary.LittleEndian, &requestNumber)
	if err != nil {
		return 0, err
	}
	return requestNumber, nil
}
