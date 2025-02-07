package util

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"gateway/config"
	"gateway/data_access"
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

func CheckMessageSignature(msg *pkg.MessageData, sourceIp string, sourcePort string, isSourceVerifier bool, util pkg.ProtocolUtil, cfg *config.Config) (bool, error) {
	if isSourceVerifier {
		vDa := data_access.VerifierDA{}
		verfifer, err := vDa.GetVerifierByIpAndPort(sourceIp, sourcePort)
		if err != nil {
			return false, err
		}
		return util.VerifyMessageDataSignature(*msg, verfifer.PublicKey, cfg.Security.DSAScheme)
	} else {
		gDa := data_access.GatewayDA{}
		gateway, err := gDa.GetGatewayByIpAndPort(sourceIp, sourcePort)
		if err != nil {
			return false, err
		}
		return util.VerifyMessageDataSignature(*msg, gateway.PublicKey, cfg.Security.DSAScheme)
	}

}

func CheckMessageHMac(msg *pkg.MessageData, symmetricKey string, util pkg.ProtocolUtil) (bool, error) {
	return util.VerifyHmac(*msg, symmetricKey)
}

func GenerateRequestNumber() (int64, error) {
	var requestNumber int64
	err := binary.Read(rand.Reader, binary.LittleEndian, &requestNumber)
	if err != nil {
		return 0, err
	}
	return requestNumber, nil
}
