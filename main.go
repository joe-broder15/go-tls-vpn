package main

import (
	"flag"
	"fmt"
	"go-tls-vpn/client"
	"go-tls-vpn/server"
)

func main() {
	// set the command line flags
	serverPtr := flag.Bool("s", false, "run as a server")
	flag.Parse()

	// determine whether to run the client or server
	if *serverPtr {
		fmt.Println(server.ServerHello())
	} else {
		client.Hello()
	}
}
