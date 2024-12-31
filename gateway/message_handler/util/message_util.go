package util

import (
	"crypto/rand"
	"database/sql"
	b64 "encoding/base64"
	"encoding/hex"
	"gateway/config"
	"gateway/data"
	_ "github.com/go-sql-driver/mysql"
	"os"
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
func GetCurrentGatewayUser(c *config.Config) (data.GatewayUser, error) {
	db, err := GetDBConnection(*c)
	if err != nil {
		return data.GatewayUser{}, err
	}
	defer db.Close()
	return data.GetGatewayUserByPassword(db, b64.StdEncoding.EncodeToString([]byte(os.Getenv("PQ_NS_IOP_GU_PASS"))))

}

func GetVerifierByPublicSigKey(c *config.Config, publicKey string) (data.Verifier, error) {
	db, err := GetDBConnection(*c)
	if err != nil {
		return data.Verifier{}, err
	}
	defer db.Close()
	return data.GetVerifierByPublicKey(db, publicKey)
}
