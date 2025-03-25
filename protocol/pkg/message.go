package pkg

import (
	"time"
)

const (
	GENERAL_ERROR = -1

	GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID  = 1
	GATEWAY_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID = 2
	GATEWAY_VERIFIER_OPERATION_ERROR_ID                     = 3

	GATEWAY_VERIFIER_GET_INFO_OPERATION_REQEST   = 4
	GATEWAY_VERIFIER_GET_INFO_OPERATION_RESPONSE = 5

	VERIFIER_VERIFIER_KEY_DISTRIBUTION_OPERATION_REQUEST_ID  = 6
	VERIFIER_VERIFIER_KEY_DISTRIBUTION_OPERATION_RESPONSE_ID = 7

	VERIFIER_VERIFIER_GET_INFO_OPERATION_REQEST_ID   = 8
	VERIFIER_VERIFIER_GET_INFO_OPERATION_RESPONSE_ID = 9
	VERIFIER_VERIFIER_OPERATION_ERROR_ID             = 500
)

type Message struct {
	IsEncrypted  bool   "json:isEncrypted"
	MsgInfo      string "json:data"
	Signature    string "json:signature"
	Hmac         string "json:hmac"
	MsgTicket    string "json:ticket"
	PublicKeySig string "json:publicKeySig"
}

type MessageInfo struct {
	RequestId       int64       "json:requestId"
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
