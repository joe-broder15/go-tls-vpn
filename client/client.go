package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
)

func Start(ip string, port string) {

	fmt.Println("client is starting...")
	//create tls configuration
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	//connect to server
	conn, err := tls.Dial("tcp", ip+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("client is connected!")
	defer conn.Close()

	// create tun interface
	/*
		tun, err := water.New(water.Config{
			DeviceType: water.TUN,
		})
		if err != nil {
			log.Fatal(err)
		}
	*/
	//read data from server
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Println("Client read:", string(buf[:n]))
}
