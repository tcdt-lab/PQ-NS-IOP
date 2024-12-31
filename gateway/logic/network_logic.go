package logic

import (
	"gateway/config"
	"gateway/data"
	"gateway/message_handler/gateway_verifier/message_creator"
	"gateway/network"
	"go.uber.org/zap"
)

type NetworkLogic struct {
}

func (networkLogic *NetworkLogic) SendKeyDistributionOperationGatewayVerifier() (bool, error) {
	cfg, err := config.ReadYaml()
	if err != nil {
		zap.L().Error("Error reading config.yaml file", zap.Error(err))
		return false, err
	}
	bootStreapVerifier := data.Verifier{}
	bootStreapVerifier.Ip = cfg.BootstrapNode.Ip
	bootStreapVerifier.Port = cfg.BootstrapNode.Port
	msg := message_creator.CreateGatewayVerifierKeyDistributionMessage(cfg)
	network.SendAndAwaitReplyToVerifier(bootStreapVerifier, msg)
	//zap.L().Info("Response from verifier", zap.ByteString("response", res))
	return true, nil
}
