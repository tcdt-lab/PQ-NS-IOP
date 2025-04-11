package gateway_verifier

type GatewayVerifierTicketRequest struct {
	RequestId             int64  `json:"requestId"`
	SourceServerIP        string `json:"sourceServerIP"`
	SourceServerPort      string `json:"sourceServerPort"`
	DestinationServerIP   string `json:"destinatingServerIP"`
	DestinationServerPort string `json:"destinatingServerPort"`
}

type GatewayVerifierTicketResponse struct {
	RequestId    int64  `json:"requestId"`
	TicketKey    string `json:"ticketKey"`
	TicketString string `json:"ticketString"`
}
