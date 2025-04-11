package balance_verification

import (
	"database/sql"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
)

func CreateBalanceVerificationRequest(verifierPublicKey string, requestId int64, publicInput string, proof string, c config.Config, db *sql.DB) ([]byte, error) {
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	guDa := data_access.GenerateGatewayUserDA(db)
	vDa := data_access.GenerateVerifierDA(db)
	cacheHandler := data_access.NewCacheHandlerDA()
	adminId, err := cacheHandler.GetUserAdminId()
	if err != nil {
		return nil, err
	}
	gatewayUser, err := guDa.GetGatewayUser(adminId)
	if err != nil {
		return nil, err
	}
	var msg = pkg.Message{}
	var msgInfo = pkg.MessageInfo{}
	var balanceVerificationRequest = gateway_verifier.VerificationRequest{}

	balanceVerificationRequest.RequestId = requestId
	balanceVerificationRequest.PublicInputs = publicInput
	balanceVerificationRequest.Proof = proof

	msgInfo.Params = balanceVerificationRequest
	msgInfo.RequestId = requestId
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_BALANCE_VERIFICATION_REQUEST_ID

	msgInfoByte, err := protocolUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err
	}

	destinationVerifier, err := vDa.GetVerifierByPublicKey(verifierPublicKey)
	protocolUtil.SignMessageInfo(&msg, msgInfoByte, gatewayUser.SecretKeyDsa, c.Security.DSAScheme)
	encMsgInfoStr, _, err := protocolUtil.EncryptMessageInfo(msgInfoByte, destinationVerifier.SymmetricKey)

	msg.PublicKeySig = gatewayUser.PublicKeyDsa
	msg.MsgInfo = encMsgInfoStr
	msg.IsEncrypted = true
	msgBytes, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}
