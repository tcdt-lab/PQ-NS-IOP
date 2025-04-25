package network

import (
	"bytes"
	"database/sql"
	"gateway/config"
	"log"
	"sync"

	"gateway/message_handler"

	"go.uber.org/zap"
	"io"
	"net"
)

func StartServer(config *config.Config, db *sql.DB) {
	zap.L().Info("Starting server")
	var mu sync.Mutex
	listener, err := net.Listen(config.Server.Protocol, "127.0.0.1:"+config.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn, config, db, &mu)
	}
}

func handleConnection(conn net.Conn, config *config.Config, db *sql.DB, mutex *sync.Mutex) {
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

	var messageParser = message_handler.GenerateNewMessageHandler(db)
	response, err := messageParser.HandleRequests(buffer.Bytes(), conn.RemoteAddr().String(), conn.RemoteAddr().Network(), *config, mutex)
	if err != nil {
		zap.L().Error("Error parsing message: ", zap.Error(err))
	}
	_, err = conn.Write(response)
	if err != nil {
		zap.L().Error("Error writing response: ", zap.Error(err))
	}
	buffer.Reset()
}
