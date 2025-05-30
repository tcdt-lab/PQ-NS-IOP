package network

import (
	"bytes"
	"gateway/data"
	"io"
	"net"
)

func SendAndAwaitReplyToGateway(ip string, port string, msg []byte) ([]byte, error) {

	//socket client to send the message to the gateway
	//localaddR, _ := net.ResolveTCPAddr("tcp", cfg.Client.Ip+":"+cfg.Client.Port)
	tcpAdr, _ := net.ResolveTCPAddr("tcp", ip+":"+port)
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
	//zap.L().Info("Response from verifier", zap.ByteString("response", response.Bytes()))
	//fmt.Println("Response from verifier", response.Bytes())

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

func SendToVerifier(verifier data.Verifier, msg []byte) error {

	//socket client to send the msg to the verifier_verifier
	conn, err := net.Dial("tcp", verifier.Ip+":"+verifier.Port)
	if err != nil {
		return nil
	}
	defer conn.Close()

	_, err = conn.Write(msg)
	return nil
}

func SendAndAwaitReply(ip string, port string, msg []byte) ([]byte, error) {

	//socket client to send the msg to the verifier_verifier
	//conn, err := net.Dial("tcp", verifier.Ip+":"+verifier.Port)
	//cfg, err := config.ReadYaml()
	//if err != nil {
	//	return nil, err
	//}

	//localaddR, _ := net.ResolveTCPAddr("tcp", cfg.Client.Ip+":"+cfg.Client.Port)
	tcpAdr, _ := net.ResolveTCPAddr("tcp", ip+":"+port)
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
	//zap.L().Info("Response from verifier", zap.ByteString("response", response.Bytes()))
	//fmt.Println("Response from verifier", response.Bytes())

	return response.Bytes(), nil
}

func parseMessage(buffer []byte) ([]byte, error) {
	return buffer, nil
}
