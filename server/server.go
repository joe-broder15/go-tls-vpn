package server

import (
	"crypto/tls"
	"fmt"
)

func Start(certFile string, keyFile string) {
	fmt.Println("Starting server...")

	_, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// config := &tls.Config{Certificates: []tls.Certificate{cert}}
}
