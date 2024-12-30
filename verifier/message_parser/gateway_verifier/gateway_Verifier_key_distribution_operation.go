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

// As this step is for key setup, it is the only place we have to verify the signiture
// In other requests, we can verify the signiture in the parser
func GatewayVerifierKeyDistributionHandler(msgData pkg.MessageData, c config.Config) ([]byte, error) {
	var gvKeyDistributionReq gateway_verifier.GatewayVerifierKeyDistributionRequest
	gvKeyDistributionReq = msgData.MsgInfo.Params.(gateway_verifier.GatewayVerifierKeyDistributionRequest)
	var gateway data.Gateway
	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	db, err := util.GetDBConnection(c)
	if err != nil {
		zap.L().Error("ErrorParams while getting db connection", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}
	gateway.Ip = gvKeyDistributionReq.GatewayIP
	gateway.Port = gvKeyDistributionReq.GatewayPort
	gateway.SigScheme = gvKeyDistributionReq.GatewaySignatureScheme
	gateway.PublicKeySig = gvKeyDistributionReq.GatewayPublicKeySignature
	gateway.PublicKeyKem = gvKeyDistributionReq.GatewayPublicKeyKem
	gateway.KemScheme = gvKeyDistributionReq.GatewayKemScheme

	signatureCheck, err := checkSignature(msgData, msgUtil, gateway)
	if err != nil {
		zap.L().Error("ErrorParams while checking signature", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}
	if !signatureCheck {
		err = errors.New("Signature verification failed")
		zap.L().Error("ErrorParams while checking signature", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}

	currentUserInfo, err := util.GetUserInformation(db, c)
	if err != nil {
		zap.L().Error("ErrorParams while getting user secret key", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}
	cipherText, sharedSymmetricKey, err := msgUtil.AsymmetricHandler.KemGenerateSecretKey(currentUserInfo.SecretKeyKem, gateway.PublicKeyKem, "", gateway.KemScheme)
	if err != nil {

		zap.L().Error("ErrorParams while generating shared symmetric key", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}
	gateway.SymmetricKey = msgUtil.AesHandler.ConvertKeyBytesToStr64(sharedSymmetricKey)
	_, err = data.AddGateway(db, gateway)
	if err != nil {
		fmt.Println(err)
		zap.L().Error("ErrorParams while adding gateway", zap.Error(err))
		return generateErrorResponse(db, c, err, gateway.SigScheme)
	}
	return generateResponse(gateway.SymmetricKey, currentUserInfo.PublicKeyKem, cipherText, msgData.MsgInfo.Nonce, c)

}

func checkSignature(msgData pkg.MessageData, protocolUtil pkg.ProtocolUtil, gateway data.Gateway) (bool, error) {
	if msgData.Signature != "" {
		res, err := protocolUtil.VerifyMessageDataSignature(msgData, gateway.PublicKeySig, gateway.SigScheme)
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
	msgUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	verifeirCurrentUser, err := util.GetUserInformation(db, c)
	if err != nil {
		return nil, err
	}
	msgUtil.SignMessageInfo(&messgaeData, verifeirCurrentUser.SecretKeySig, schemeName)
	msgDataBytes, err := msgUtil.ConvertMessageDataToByte(messgaeData)
	if err != nil {
		zap.L().Error("ErrorParams while converting message data to byte", zap.Error(err))
		return nil, err
	}
	msgDataStr := base64.StdEncoding.EncodeToString(msgDataBytes)
	message.Data = msgDataStr
	message.IsEncrypted = false
	msgByte, err := msgUtil.ConvertMessageToByte(message)
	if err != nil {
		zap.L().Error("ErrorParams while converting message to byte", zap.Error(err))
		return nil, err
	}
	return msgByte, nil

}

func generateResponse(symmetricKey string, verifierUSerPubKeyKem string, cipherText []byte, nonce string, c config.Config) ([]byte, error) {
	var gvKeyDistributionRes gateway_verifier.GatewayVerifierKeyDistributionResponse
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	var messageInfo pkg.MessageInfo
	var messgaeData pkg.MessageData
	var message pkg.Message
	gvKeyDistributionRes.CipherText = base64.StdEncoding.EncodeToString(cipherText)
	gvKeyDistributionRes.OperationError = ""
	gvKeyDistributionRes.PublicKeyKem = verifierUSerPubKeyKem
	messageInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID
	messageInfo.Params = gvKeyDistributionRes
	messageInfo.Nonce = nonce
	messgaeData.MsgInfo = messageInfo
	protocolUtil.GenerateHmacMsgInfo(&messgaeData, symmetricKey)
	msgDataBytes, err := protocolUtil.ConvertMessageDataToByte(messgaeData)
	if err != nil {
		zap.L().Error("ErrorParams while converting message data to byte", zap.Error(err))
		return nil, err
	}
	msgDataStr := base64.StdEncoding.EncodeToString(msgDataBytes)

	message.Data = msgDataStr
	message.IsEncrypted = false
	msgByte, err := protocolUtil.ConvertMessageToByte(message)
	if err != nil {
		zap.L().Error("ErrorParams while converting message to byte", zap.Error(err))
		return nil, err
	}
	return msgByte, nil
}
