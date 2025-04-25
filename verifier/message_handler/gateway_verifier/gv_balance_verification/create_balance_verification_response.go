package gv_balance_verification

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"
	"go.uber.org/zap"
	"os"
	"sync"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_verifier"
	"verifier/config"
	"verifier/data_access"
	"verifier/message_handler/util"
)

func CreateBalanceVerificationResponse(proof string, publicInputs string, requestId int64, senderPublicKey string, db *sql.DB, c config.Config, mutex *sync.Mutex) ([]byte, error) {
	// This function is not implemented yet.
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	vuDa := data_access.GenerateVerifierUserDA(db)
	gDa := data_access.GenerateGatewayDA(db)
	cacheHandler := data_access.GenerateCacheHandlerDA()

	var msg = pkg.Message{}
	var msgInfo = pkg.MessageInfo{}
	adminId, err := cacheHandler.GetUserAdminId()
	if err != nil {
		return nil, err
	}
	verifierUser, err := vuDa.GetVerifierUser(adminId)
	if err != nil {
		return nil, err
	}
	senderGateway, err := gDa.GetGatewayByPublicKeySig(senderPublicKey)
	if err != nil {
		return nil, err
	}
	//mutex.Lock()
	proofResult, err := CheckProof(proof, publicInputs, c.ZKP.VerificationKeyPath)
	//mutex.Unlock()

	if err != nil {
		return nil, err
	}
	resp := gateway_verifier.VerificationResponse{}
	resp.RequestId = requestId
	resp.VerificationResult = proofResult
	msgInfo.Params = resp
	msgInfo.OperationTypeId = pkg.GATEWAY_VERIFIER_BALANCE_VERIFICATION_RESPONSE_ID
	msgInfoBytes, err := protocolUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err
	}

	protocolUtil.SignMessageInfo(&msg, msgInfoBytes, verifierUser.SecretKeySig, c.Security.DSAScheme)
	msgInfoEncStr, _, err := protocolUtil.EncryptMessageInfo(msgInfoBytes, senderGateway.SymmetricKey)
	if err != nil {
		return nil, err
	}

	msg.MsgInfo = msgInfoEncStr
	msg.PublicKeySig = verifierUser.PublicKeySig
	msg.IsEncrypted = true

	msgBytes, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}

func CheckProof(proof string, publicInputs string, verificationKeyPath string) (bool, error) {
	zap.L().Info("CheckProof", zap.String("proof", proof), zap.String("publicInputs", publicInputs))
	vk, err := os.ReadFile(verificationKeyPath)
	var proofData types.ProofData
	var pubSignals []string

	err = json.Unmarshal([]byte(proof), &proofData)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(publicInputs), &pubSignals)
	if err != nil {
		fmt.Println(err)
	}

	zkProof := &types.ZKProof{Proof: &proofData, PubSignals: pubSignals}
	err = verifier.VerifyGroth16(*zkProof, vk)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	fmt.Println("Verifer passed")
	return true, nil
}
