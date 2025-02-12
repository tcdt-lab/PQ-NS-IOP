package network

import (
	"bytes"
	"encoding/hex"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
	"verifier/config"
	"verifier/message_parser"
)

func StartServer(config *config.Config) {
	// Start the server
	zap.L().Info("Starting server")
	listener, err := net.Listen(config.Server.Protocol, "127.0.0.1:"+config.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn, config)
	}
}

func handleConnection(conn net.Conn, config *config.Config) {
	// Handle the connection
	defer conn.Close()
	var buffer bytes.Buffer
	var bufferSize = 0
	var i = 0
	for {
		i += 1

		buf := make([]byte, config.Server.BufferSize)
		n, err := conn.Read(buf)

		buffer.Write(buf[:n])
		bufferSize += n

		if err == io.EOF {
			zap.L().Error("No more info to read:", zap.Error(err))
			break
		}
		if err != nil {
			zap.L().Error("Error reading:", zap.Error(err))
			return
		}

	}
	zap.L().Info("Received data: ", zap.String("data", hex.EncodeToString(buffer.Bytes())))

	var messageParser = message_parser.MessageHandler{}
	response, err := messageParser.HandleRequests(buffer.Bytes(), conn.RemoteAddr().String(), conn.RemoteAddr().Network(), *config)
	if err != nil {
		zap.L().Error("Error parsing message: ", zap.Error(err))
	}
	_, err = conn.Write(response)
	if err != nil {
		zap.L().Error("Error writing response: ", zap.Error(err))
	}
	buffer.Reset()
}
