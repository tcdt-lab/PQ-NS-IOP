package gateway_verifier

import (
	"gateway/config"
	"gateway/data_access"
	"gateway/message_handler"
	"gateway/message_handler/gateway_verifier/message_applier"
	"gateway/message_handler/gateway_verifier/message_creator"
	"gateway/network"
)

func GetInformationLogicApply() error {
	cfg, err := config.ReadYaml()
	if err != nil {
		return err
	}
	vDa := data_access.VerifierDA{}
	bootstrapVerifier, err := vDa.GetVerifierByIpAndPort(cfg.BootstrapNode.Ip, cfg.BootstrapNode.Port)
	msgBytes, err := message_creator.CreateGatewayVerifierGetInfoOperationMessage(cfg, bootstrapVerifier.Ip, bootstrapVerifier.Port)
	if err != nil {
		return err
	}
	res, err := network.SendAndAwaitReplyToVerifier(bootstrapVerifier, msgBytes)
	if err != nil {
		return err
	}
	msgdata, err := message_handler.ParseGatewayVerifierResponse(res, bootstrapVerifier.Ip, bootstrapVerifier.Port)
	if err != nil {
		return err
	}
	err = message_applier.ApplyGatewayVerifierGetInfoResponse(msgdata)
	if err != nil {
		return err
	}
	return nil
}
