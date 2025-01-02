package data_access

//
//import (
//	"crypto/sha256"
//	"test.org/cryptography/pkg/asymmetric/pq"
//
//	"database/sql"
//	b64 "encoding/base64"
//	"encoding/hex"
//	"fmt"
//	"gateway/config"
//	"gateway/data"
//	"go.uber.org/zap"
//	"os"
//	"test.org/cryptography/pkg/symmetric"
//)
//
//func Login(db *sql.DB) (bool, []byte) {
//	c, err := config.ReadYaml()
//	if err != nil {
//		zap.L().Error("Error while reading config.yaml file", zap.Error(err))
//		return false, nil
//	}
//	if c.Security.LoginType == "PASSWORD" {
//		fmt.Println("If you have already a public key and password, enter 1. If you want to create a new account, enter 2")
//		var choice int
//		_, err := fmt.Scanln(&choice)
//		if err != nil {
//			zap.L().Error("Error while reading user input", zap.Error(err))
//			return false, nil
//		}
//		if choice == 1 {
//			fmt.Println("Enter the first five characters of your public key")
//			var firstFivePublicKey string
//			fmt.Scanln(&firstFivePublicKey)
//			fmt.Println("Enter your rawPassword")
//			var rawPassword string
//			fmt.Scanln(&rawPassword)
//			credentialResult, user := checkCredentials(db, firstFivePublicKey, rawPassword)
//			if credentialResult {
//				fmt.Println("GatewayUser signed in successfully")
//				return true, extractUserAESKey(rawPassword, user.Salt)
//
//			}
//			return false, nil
//		} else if choice == 2 {
//			fmt.Println("Enter your rawPassword")
//			var rawPassword string
//			fmt.Scanln(&rawPassword)
//			user, err := CreateUser(rawPassword, db)
//			if err != nil {
//				zap.L().Error("Error while creating data_access", zap.Error(err))
//				return false, nil
//			}
//			fmt.Println("GatewayUser signed in successfully")
//			return true, extractUserAESKey(rawPassword, user.Salt)
//		}
//
//	} else if c.Security.LoginType == "ENV" {
//		pass := os.Getenv("GATEWAY_PATH")
//		credentialResult, user := checkCredentialsWithPass(db, pass)
//		if credentialResult {
//			return true, extractUserAESKey(pass, user.Salt)
//		}
//		return false, nil
//	}
//	fmt.Println("Wrong YAML file configuration!")
//	return false, nil
//}
//
//func CreateUser(password string, db *sql.DB) (data.GatewayUser, error) {
//	keyDrivation := symmetric.PBKDF2{}
//	salt, _ := keyDrivation.GeneratingSalt(8)
//	var mldsa = pq.MLDSA{}
//	c, err := config.ReadYaml()
//	if err != nil {
//		zap.L().Error("Error while reading config.yaml file", zap.Error(err))
//		return data.GatewayUser{}, err
//	}
//	publicKey, secretKey, err := mldsa.KeyGen(c.Security.DSAScheme)
//	if err != nil {
//		zap.L().Error("Error while generating ML-DSA keys", zap.Error(err))
//		return data.GatewayUser{}, err
//	}
//	pkStr64 := mldsa.ConvertPubKeyToBase64String(publicKey)
//	encodedSalt := b64.StdEncoding.EncodeToString(salt)
//	secKeyBytes, _ := secretKey.MarshalBinary()
//	encodedSecKey, err := secKeyEncryption(secKeyBytes, password, encodedSalt)
//	if err != nil {
//		return data.GatewayUser{}, err
//	}
//	user := data.GatewayUser{Password: userPasswordEncoder(password), PublicKey: pkStr64, SecretKey: encodedSecKey, Salt: encodedSalt}
//	_, errAdd := data.AddUser(db, user)
//	if errAdd != nil {
//		return data.GatewayUser{}, errAdd
//	}
//	fmt.Println("GatewayUser created successfully")
//	fmt.Println("first five characters of your encoded public key: ", user.PublicKey[:5])
//	return user, nil
//
//}
//
//// Password parameter is the raw password
//func checkCredentialsWithPass(db *sql.DB, rawPassword string) (bool, data.GatewayUser) {
//
//	user, err := data.GetUserByPassword(db, userPasswordEncoder(rawPassword))
//	if err != nil {
//		zap.L().Error("Error while getting user by password", zap.Error(err))
//		return false, data.GatewayUser{}
//	}
//	if user.Password == userPasswordEncoder(rawPassword) {
//		return true, user
//	}
//	return false, data.GatewayUser{}
//}
//
//// Password parameter is the raw password
//func checkCredentials(db *sql.DB, firstFivePublicKey string, rawPassword string) (bool, data.GatewayUser) {
//
//	user, err := data.GetUserByPassword(db, userPasswordEncoder(rawPassword))
//	if err != nil {
//		zap.L().Error("Error while getting user by password", zap.Error(err))
//		return false, data.GatewayUser{}
//	}
//	if user.Password == userPasswordEncoder(rawPassword) && user.PublicKey[:5] == firstFivePublicKey {
//		return true, user
//	}
//	return false, data.GatewayUser{}
//}
//
//// Input salt is b64 encoded
//func extractUserAESKey(rawPassword string, salt string) []byte {
//	decodedSalt, err := b64.StdEncoding.DecodeString(salt)
//	if err != nil {
//		zap.L().Error("Error while decoding salt", zap.Error(err))
//		return nil
//	}
//	keyDerivation := symmetric.PBKDF2{}
//	key := keyDerivation.KeyDerivation([]byte(rawPassword), decodedSalt, 4096)
//	return key
//}
//
//func userPasswordEncoder(password string) string {
//	h := sha256.New()
//	h.Write([]byte(password))
//	hashedPass := h.Sum(nil)
//	return hex.EncodeToString(hashedPass)
//}
//
//func secKeyEncryption(secretKey []byte, password string, encodedSalt string) (string, error) {
//	aesKey := extractUserAESKey(password, encodedSalt)
//	aesGcm := symmetric.AesGcm{}
//	encryptedKey, err := aesGcm.Encrypt(secretKey, aesKey)
//	if err != nil {
//		zap.L().Error("Error while encrypting secret key", zap.Error(err))
//		return "", err
//	}
//	return b64.StdEncoding.EncodeToString(encryptedKey), nil
//}
