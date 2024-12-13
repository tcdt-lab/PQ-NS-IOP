package util

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"test.org/cryptography/pkg/asymmetric"
	"test.org/cryptography/pkg/symmetric"
	"test.org/protocol/pkg"
	"verifier/config"
)

func GetDBConnection(c config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", c.DB.Username+":"+c.DB.Password+"@/"+c.DB.Name)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MessageUtilGenerator(cryptographySchemeName string) pkg.MessageUtil {
	var util pkg.MessageUtil
	util.AesHandler = symmetric.AesGcm{}
	util.AsymmetricHandler = asymmetric.NewAsymmetricHandler(cryptographySchemeName)
	util.HmacHandler = symmetric.HMAC{}
	util.RegisterInterfacesInGob()
	return util
}
