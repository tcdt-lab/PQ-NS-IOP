package balance_check

import (
	"database/sql"
	"fmt"
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/util"
	"github.com/iden3/go-rapidsnark/prover"
	"os"
	"test.org/protocol/pkg"
	"test.org/protocol/pkg/gateway_gateway"
)

func CreateBalanceCheckResponse(ticketKey string, requestId int64, c *config.Config, db *sql.DB) ([]byte, error) {
	var msg = pkg.Message{}
	var msgInfo = pkg.MessageInfo{}
	protocolUtil := util.ProtocolUtilGenerator(c.Security.CryptographyScheme)
	guDa := data_access.GenerateGatewayUserDA(db)
	cacheHandler := data_access.NewCacheHandlerDA()

	var balanceCheckResponse = gateway_gateway.BalanceCheckResponse{}

	publicInputs, proof, err := CreateProof(c)
	if err != nil {
		return nil, err
	}
	balanceCheckResponse.RequestId = requestId
	balanceCheckResponse.Proof = proof
	balanceCheckResponse.PublicInputs = publicInputs
	msgInfo.Params = balanceCheckResponse
	msgInfo.RequestId = requestId
	msgInfo.OperationTypeId = pkg.GATEWAY_GATEWAY_BALANCE_CHECK_RESPONSE_ID

	adminId, err := cacheHandler.GetUserAdminId()
	if err != nil {
		return nil, err
	}
	gatewayUser, err := guDa.GetGatewayUser(adminId)
	if err != nil {
		return nil, err
	}
	msgInfoBytes, err := protocolUtil.ConvertMessageInfoToByte(msgInfo)
	if err != nil {
		return nil, err

	}
	protocolUtil.SignMessageInfo(&msg, msgInfoBytes, gatewayUser.SecretKeyDsa, c.Security.DSAScheme)
	msginfoEncrypted, _, err := protocolUtil.EncryptMessageInfo(msgInfoBytes, ticketKey)
	if err != nil {
		return nil, err
	}

	msg.PublicKeySig = gatewayUser.PublicKeyDsa
	msg.MsgInfo = msginfoEncrypted
	msg.IsEncrypted = true
	msgBytes, err := protocolUtil.ConvertMessageToByte(msg)
	if err != nil {
		return nil, err

	}
	return msgBytes, nil
}

func CreateProof(c *config.Config) (pubicInputs string, proof string, err error) {
	zkeyBytes, err := os.ReadFile(c.ZKP.ProverKeyPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read zkey file: %v\n", err)
		os.Exit(1)
	}

	wtnsBytes, err := os.ReadFile(c.ZKP.WitnessPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read witness file: %v\n", err)
		os.Exit(1)
	}

	proof, publicInputs, err := prover.Groth16ProverRaw(zkeyBytes, wtnsBytes)
	prover.Groth16Prover(zkeyBytes, wtnsBytes)
	fmt.Println(proof, publicInputs, err)
	return proof, publicInputs, err
}
