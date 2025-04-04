package gv_key_distribution

import (
	"database/sql"
	b64 "encoding/base64"
	"go.uber.org/zap"
	"strconv"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func CreateGatewayVerifierKeyDistributionResponse(cipherTextStr string, db *sql.DB) ([]byte, error) {
	cfg, err := config.ReadYaml()
	if err != nil {
		return nil, err
	}

	vuda := data_access.GenerateVerifierUserDA(db)

	cacheHandler := data_access.GenerateCacheHandlerDA()
	adminId, _ := cacheHandler.GetUserAdminId()
	vuda.GetVerifierUser(adminId)
	protoUtil := util.ProtocolUtilGenerator(cfg.Security.CryptographyScheme)
	msg := pkg.Message{}

	msgInfo := pkg.MessageInfo{}
	adminVerifierUSer, err := vuda.GetAdminVerifierUser()
	if err != nil {
		return nil, err
	}
	gvKeyDistributionResponse := gateway_verifier.GatewayVerifierKeyDistributionResponse{}
	gvKeyDistributionResponse.CipherText = cipherTextStr
	gvKeyDistributionResponse.PublicKeyKem = adminVerifierUSer.PublicKeyKem
	msgInfo.Params = gvKeyDistributionResponse
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID
	msgInfo.Nonce = "123"
	msgInfoBytes, _ := protoUtil.ConvertMessageInfoToByte(msgInfo)

	hmacStr, _, err := protoUtil.GenerateHmacMsgInfo(msgInfoBytes, adminVerifierUSer.SymmetricKey)
	if err != nil {
		return nil, err
	}
	msg.Hmac = hmacStr
	msg.MsgInfo = b64.StdEncoding.EncodeToString(msgInfoBytes)
	msg.IsEncrypted = false
	msg.MsgTicket = ""
	msgBytes, err := protoUtil.ConvertMessageToByte(msg)
	zap.L().Info("Generated Response is", zap.String("requnumber", strconv.FormatInt(msgInfo.RequestId, 10)))
	return msgBytes, nil
}
