package logic

import (
	"verifier/message_parser/gateway_verifier/message_applier"
	"verifier/message_parser/gateway_verifier/message_creator"
	"verifier/message_parser/gateway_verifier/message_parser"
)

type ResponseLogic struct{}

func (rl *ResponseLogic) HandleKeyDistributionResponse(message []byte, senderIp string, senderPort string) ([]byte, error) {

	msgData, err := message_parser.ParseRequest(message, senderIp, senderPort)
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

func (rl *ResponseLogic) HandleGetInfoResponse(message []byte, senderIp string, senderPort string) ([]byte, error) {
	_, err := message_parser.ParseRequest(message, senderIp, senderPort)
	if err != nil {
		return nil, err
	}

	res, err := message_creator.CreateGateVerifierGetInfoResponse()
	if err != nil {
		return nil, err
	}
	return res, nil
}
