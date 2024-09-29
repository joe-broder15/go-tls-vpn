package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func Start(certFile string, keyFile string, port string) {
	fmt.Println("Starting server...")

	// load the x509 Key Pair
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// create the tls configuration
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	// start the listener and set up a defer to close it
	listener, err := tls.Listen("tcp", ":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// listen for incoming connections
	for {
		// accept connections and hand them off to a new goroutine
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {

			// send some data to the connected client
			conn.Write([]byte("hello from the server over TLS\n"))

		}(conn)
	}

}
