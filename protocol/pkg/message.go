package pkg

import (
	"time"
)

const (
	GENERAL_ERROR                                          = -1
	GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID = 1

	GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID = 2
	GATEWAY_VERIFIER_OPERATION_ERROR_ID                     = 3

	GATEWAY_VERIFIER_GET_INFO_OPERATION_REQEST   = 4
	GATEWAY_VERIFIER_GET_INFO_OPERATION_RESPONSE = 5
)

type Message struct {
	IsEncrypted bool   "json:isEncrypted"
	Data        string "json:data"
	MsgTicket   string "json:ticket"
}
type MessageData struct {
	MsgInfo   MessageInfo "json:data"
	Signature string      "json:signature"
	Hmac      string      "json:hmac"
}

type MessageInfo struct {
	OperationTypeId int         "json:operationTypeId"
	Params          interface{} "json:params"
	SourceId        string      "json:sourceId"
	DestinationId   string      "json:destinationId"
	Nonce           string      "json:nonce"
}

type Ticket struct {
	SharedKey     string    "json:sharedKey"
	Deadline      time.Time "json:deadline"
	SourceIp      string    "json:sourceIp"
	DestinationIp string    "json:destinationIp"
}

type ErrorParams struct {
	ErrorCode    int    "json:errorCode"
	ErrorMessage string "json:errorMessage"
}
