package gateway_verifier

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"verifier/data"
	"verifier/message_parser/util"

	"verifier/config"
)

func GatewayVerifierKeyDistributionHandler(msgData pkg.MessageData, c config.Config) ([]byte, error) {
	var gvKeyDistributionReq gateway_verifier.GatewayVerifierKeyDistributionRequest
	gvKeyDistributionReq = msgData.MsgInfo.Params.(gateway_verifier.GatewayVerifierKeyDistributionRequest)
	var gateway data.Gateway
	msgUtil := util.MessageUtilGenerator(c.Security.CryptographyScheme)
	db, err := util.GetDBConnection(c)
	if err != nil {
		zap.L().Error("Error while getting db connection", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}
	gateway.Ip = gvKeyDistributionReq.GatewayIP
	gateway.Port = gvKeyDistributionReq.GatewayPort
	gateway.SigScheme = gvKeyDistributionReq.GatewaySignatureScheme
	gateway.PublicKeySig = gvKeyDistributionReq.GatewayPublicKeySignature
	gateway.PublicKeyKem = gvKeyDistributionReq.GatewayPublicKeyKem
	gateway.KemScheme = gvKeyDistributionReq.GatewayKemScheme

	res, err2 := checkSignature(msgData, msgUtil, gateway)
	if err2 != nil {
		db, err := util.GetDBConnection(c)
		if err != nil {
			zap.L().Error("Error while getting db connection", zap.Error(err))
			return generateErrorResponse(db, c, err, gateway.SigScheme)
		}
	}
	if !res {
		err = errors.New("Signature verification failed")
		zap.L().Error("Signature verification failed", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}

	cipherText, sharedSymmetricKey, err := msgUtil.AsymmetricHandler.KemGenerateSecretKey("", gateway.PublicKeyKem, "", gateway.KemScheme)
	if err != nil {

		zap.L().Error("Error while generating shared symmetric key", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}
	gateway.SymmetricKey = msgUtil.AesHandler.ConvertKeyBytesToStr64(sharedSymmetricKey)
	_, err = data.AddGateway(db, gateway)
	if err != nil {
		fmt.Println(err)
		zap.L().Error("Error while adding gateway", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}
	return generateResponse(gateway.SymmetricKey, cipherText, msgData.MsgInfo.Nonce, c)

}

func checkSignature(msgData pkg.MessageData, msgUtil pkg.MessageUtil, gateway data.Gateway) (bool, error) {
	if msgData.Signature != "" {
		res, err := msgUtil.VerifyMessageDataSignature(msgData, gateway.PublicKeySig, gateway.SigScheme)
		if err != nil {
			return false, err
		}
		if !res {
			return res, errors.New("Signature verification failed")
		}
		return res, nil

	} else {
		return false, errors.New("No Signature found")
	}
}

func generateErrorResponse(db *sql.DB, c config.Config, err error, schemeName string) ([]byte, error) {
	var gvKeyDistributionReq gateway_verifier.GatewayVerifierKeyDistributionResponse
	var messageInfo pkg.MessageInfo
	var messgaeData pkg.MessageData
	var message pkg.Message
	gvKeyDistributionReq.CipherText = ""
	gvKeyDistributionReq.OperationError = err.Error()
	messageInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_OPERATION_ERROR_ID
	messageInfo.Params = gvKeyDistributionReq
	messgaeData.MsgInfo = messageInfo
	msgUtil := util.MessageUtilGenerator(c.Security.CryptographyScheme)
	privateKeyStr, err := getUserPrivateKey(db, c)
	if err != nil {
		return nil, err
	}
	msgUtil.SignMessageInfo(&messgaeData, privateKeyStr, schemeName)
	msgDataBytes, err := msgUtil.ConvertMessageDataToByte(messgaeData)
	if err != nil {
		zap.L().Error("Error while converting message data to byte", zap.Error(err))
		return nil, err
	}
	msgDataStr := base64.StdEncoding.EncodeToString(msgDataBytes)
	message.Data = msgDataStr
	message.IsEncrypted = false
	msgByte, err := msgUtil.ConvertMessageToByte(message)
	if err != nil {
		zap.L().Error("Error while converting message to byte", zap.Error(err))
		return nil, err
	}
	return msgByte, nil

}

func getUserPrivateKey(db *sql.DB, c config.Config) (string, error) {
	verifeirUser, err := data.GetVerifierUserByPassword(db, c.User.Password)
	if err != nil {
		zap.L().Error("Error while getting verifier user", zap.Error(err))
		return "", err
	}
	return verifeirUser.SecretKey, nil

}

func generateResponse(symmetricKey string, cipherText []byte, nonce int, c config.Config) ([]byte, error) {
	var gvKeyDistributionReq gateway_verifier.GatewayVerifierKeyDistributionResponse
	msgUtil := util.MessageUtilGenerator(c.Security.CryptographyScheme)
	var messageInfo pkg.MessageInfo
	var messgaeData pkg.MessageData
	var message pkg.Message
	gvKeyDistributionReq.CipherText = base64.StdEncoding.EncodeToString(cipherText)
	gvKeyDistributionReq.OperationError = ""
	messageInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID
	messageInfo.Params = gvKeyDistributionReq
	messageInfo.Nonce = nonce
	messgaeData.MsgInfo = messageInfo
	msgUtil.GenerateHmac(&messgaeData, symmetricKey)
	msgDataBytes, err := msgUtil.ConvertMessageDataToByte(messgaeData)
	if err != nil {
		zap.L().Error("Error while converting message data to byte", zap.Error(err))
		return nil, err
	}
	msgDataStr := base64.StdEncoding.EncodeToString(msgDataBytes)

	message.Data = msgDataStr
	message.IsEncrypted = false
	msgByte, err := msgUtil.ConvertMessageToByte(message)
	if err != nil {
		zap.L().Error("Error while converting message to byte", zap.Error(err))
		return nil, err
	}
	return msgByte, nil
}
