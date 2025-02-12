package logic

import (
	"test.org/protocol/pkg"
	"verifier/message_parser/gateway_verifier/message_applier"
	"verifier/message_parser/gateway_verifier/message_creator"
)

type ResponseLogic struct{}

func (rl *ResponseLogic) HandleKeyDistributionResponse(msgData pkg.MessageData, senderIp string, senderPort string) ([]byte, error) {

	cipherText, err := message_applier.ApplyGatewayVerifierKeyDistributionRequest(msgData)
	if err != nil {
		return nil, err
	}
	res, err := message_creator.CreateGatewayVerifierKeyDistributionResponse(cipherText)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (rl *ResponseLogic) HandleGetInfoResponse() ([]byte, error) {

	res, err := message_creator.CreateGateVerifierGetInfoResponse()
	if err != nil {
		return nil, err
	}
	return res, nil
}
