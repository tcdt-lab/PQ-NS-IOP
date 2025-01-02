package gateway_verifier

import (
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler/gateway_verifier/message_applier"
	"gateway/message_handler/gateway_verifier/message_creator"
	"gateway/message_handler/gateway_verifier/message_parser"
	"gateway/network"
)

func KeyDistributionLogicApply() error {
	cfg, err := config.ReadYaml()
	if err != nil {
		return err
	}
	vDa := data_access.VerifierDA{}
	bootstrapVerifier, err := vDa.GetVerifierByIpAndPort(cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)
	msgBytes := message_creator.CreateGatewayVerifierKeyDistributionMessage(cfg)
	responseBytes, err := network.SendAndAwaitReplyToVerifier(bootstrapVerifier, msgBytes)
	if err != nil {
		return err
	}
	msgData, err := message_parser.ParseGatewayVerifierKeyDistributionResponse(responseBytes, bootstrapVerifier.Ip, bootstrapVerifier.Port)
	if err != nil {
		return err
	}
	err = message_applier.ApplyGatewayVerifierKeyDistributionResponse(msgData)
	if err != nil {
		return err
	}
	return nil
}
