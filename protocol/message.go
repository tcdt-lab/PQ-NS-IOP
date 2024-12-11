package protocol

import (
	"time"
)

const (
	INIT_OPEARTION_ID                = 1
	INIT_RESPONSE_OPERATION_ID       = 2
	INIT_RESPONSE_OPERATION_ERROR_ID = 3
)

type Message struct {
	Data      string "json:data"
	MsgTicket string "json:ticket"
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
	Nonce           int         "json:nonce"
}

type Ticket struct {
	SharedKey     string    "json:sharedKey"
	Deadline      time.Time "json:deadline"
	SourceIp      string    "json:sourceIp"
	DestinationIp string    "json:destinationIp"
}
