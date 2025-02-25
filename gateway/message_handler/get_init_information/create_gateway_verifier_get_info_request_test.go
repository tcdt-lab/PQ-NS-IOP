package message_creator

import (
	"encoding/hex"
	"gateway/config"
	"gateway/data"
	"gateway/data_access"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getBootstrapVerifier(t *testing.T, cfg config.Config) data.Verifier {

	vDa := data_access.VerifierDA{}
	bootstrapVerifier, err := vDa.GetVerifierByIpAndPort(cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)
	assert.NoError(t, err)
	return bootstrapVerifier
}

func TestCreateGatewayVerifierGetInfoRequest(t *testing.T) {
	cfg, err := config.ReadYaml()
	assert.NoError(t, err)
	bootstrapVerifier := getBootstrapVerifier(t, *cfg)

	res, err := CreateGatewayVerifierGetInfoOperationMessage(cfg, bootstrapVerifier.Ip, bootstrapVerifier.Port)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	t.Log(hex.EncodeToString(res))
}
