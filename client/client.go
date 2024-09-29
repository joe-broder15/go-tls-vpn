package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
)

func Start(ip string, port string) {
	fmt.Println("client is starting...")

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", ip+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//read data
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Println("Client read:", string(buf[:n]))

}
