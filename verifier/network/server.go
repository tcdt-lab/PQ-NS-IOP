package network

import (
	"log"
	"net"
	"verifier/config"
)

func startServer() {
	// Start the server
	listener, err := net.Listen(config.SERVER_PROTOCOL, config.SERVER_PORT)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Handle the connection
	defer conn.Close()

	buf := make([]byte, config.SERVER_BUFFER_SIZE)
	input, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	print(string(buf[0:input]))

}
