package key_distribution

import (
	"database/sql"
	b64 "encoding/base64"
	"gateway/config"
	"gateway/data_access"
	"go.uber.org/zap"
	"test.org/protocol/pkg"

	"gateway/message_handler/util"
	"test.org/protocol/pkg/gateway_verifier"
)

func CreateGatewayVerifierKeyDistributionMessage(c *config.Config, requestId int64, db *sql.DB) []byte {
	msg := pkg.Message{}
	msgInfo := pkg.MessageInfo{}
	guDa := data_access.GenerateGatewayUserDA(db)
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	var nonce string
	cachHandler := data_access.NewCacheHandlerDA()
	adminId, err := cachHandler.GetUserAdminId()
	currentUSer, err := guDa.GetGatewayUser(adminId)
	if err != nil {

		zap.L().Error("Error while getting current gateway user", zap.Error(err), zap.Any("db", &db))
		return nil
	}
	// Get current gateway user
	gvKeyDistributionReq := gateway_verifier.GatewayVerifierKeyDistributionRequest{}
	gvKeyDistributionReq.GatewayIP = c.Server.Ip
	gvKeyDistributionReq.GatewayPort = c.Server.Port
	gvKeyDistributionReq.GatewaySignatureScheme = c.Security.DSAScheme
	gvKeyDistributionReq.GatewayPublicKeySignature = currentUSer.PublicKeyDsa
	gvKeyDistributionReq.GatewayPublicKeyKem = currentUSer.PublicKeyKem
	gvKeyDistributionReq.GatewayKemScheme = c.Security.KEMScheme

	nonce, err = util.GenerateNonce()
	if err != nil {
		zap.L().Error("Error while generating nonce", zap.Error(err))
		return nil
	}
	msgInfo.Nonce = nonce
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID
	msgInfo.Params = gvKeyDistributionReq
	msgInfo.RequestId = requestId
	msgInfoBytes, err := protocolUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		zap.L().Error("Error while converting message info to byte", zap.Error(err))
		return nil
	}
	msgInfoStr := b64.StdEncoding.EncodeToString(msgInfoBytes)

	msg.MsgInfo = msgInfoStr
	msg.PublicKeySig = currentUSer.PublicKeyDsa
	msg.IsEncrypted = false
	msgByte, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		zap.L().Error("Error while converting message to byte", zap.Error(err))
		return nil
	}

	zap.L().Info("Key Distribution Request Message Generated")
	return msgByte
}
