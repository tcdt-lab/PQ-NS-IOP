package network

import (
	"bytes"
	"io"
	"net"
	data "verifier/data"
)

func SendAndAwaitReplyToGateway(gateway data.Gateway, encryptedMsg []byte) ([]byte, error) {

	//socket client to send the message to the gateway
	conn, err := net.Dial("tcp", gateway.Ip+":"+gateway.Port)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(encryptedMsg)
	defer conn.Close()

	var response bytes.Buffer

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err.Error() == "EOF" {
			break
		} else if err != nil {
			return nil, err
		}
		response.Write(buffer[:n])
		if n < 1024 {
			break
		}

	}
	return response.Bytes(), nil
}

func SendToGateway(gateway data.Gateway, encryptedMsg []byte) error {

	//socket client to send the encryptedMsg to the gateway
	conn, err := net.Dial("tcp", gateway.Ip+":"+gateway.Port)
	if err != nil {
		return nil
	}
	defer conn.Close()

	_, err = conn.Write(encryptedMsg)
	return nil
}

func SendToVerifier(verifier data.Verifier, encryptedMsg []byte) error {

	//socket client to send the encryptedMsg to the verifier_verifier
	conn, err := net.Dial("tcp", verifier.Ip+":"+verifier.Port)
	if err != nil {
		return nil
	}
	defer conn.Close()

	_, err = conn.Write(encryptedMsg)
	return nil
}

func SendAndAwaitReplyToVerifier(verifier data.Verifier, msg []byte) ([]byte, error) {

	tcpAdr, _ := net.ResolveTCPAddr("tcp", verifier.Ip+":"+verifier.Port)
	conn, err := net.DialTCP("tcp", nil, tcpAdr)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(msg)
	defer conn.Close()

	conn.CloseWrite()
	var response bytes.Buffer

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		response.Write(buffer[:n])
		if n < 1024 {
			break
		}

	}

	return response.Bytes(), nil
}
