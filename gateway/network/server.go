package network

import (
	"bytes"
	"encoding/hex"
	"test.org/protocol/pkg"

	"go.uber.org/zap"
	"io"
	"net"
)

type Server struct {
	messageParser pkg.Message
}

func (s *Server) StartSocketServer(port string) error {

	ln, err := net.Listen("tcp", port)
	defer ln.Close()
	if err != nil {
		zap.L().Error("Error listening:", zap.Error(err))
	}
	zap.L().Info("Listening on " + port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			zap.L().Error("Error accepting a client: ", zap.Error(err))
			continue
		}

		// Handle the connection in a new goroutine
		go s.handleConnection(conn)
	}
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	var buffer bytes.Buffer
	var bufferSize = 0
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		buffer.Write(buf[:n])
		bufferSize += n
		if err == io.EOF && err != nil {
			zap.L().Error("No more info to read:", zap.Error(err))
			break
		}
		if err != nil {
			zap.L().Error("Error reading:", zap.Error(err))
			return
		}

	}
	zap.L().Info("Received data: ", zap.String("data", hex.EncodeToString(buffer.Bytes())))
	//go func() {
	//	err := s.messageParser.ParseMessage(buffer.Bytes(), bufferSize)
	//	if err != nil {
	//		zap.L().Error("Error parsing message: ", zap.Error(err))
	//	}
	//}()

}
