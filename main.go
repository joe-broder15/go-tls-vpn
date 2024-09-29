package main

import (
	"flag"
	"fmt"
	"go-tls-vpn/server"
)

func main() {
	// set the command line flags
	serverPtr := flag.Bool("s", false, "run as a server")
	pubkeyPtr := flag.String("pubkey", "server.crt", "run as a server")
	privkeyPtr := flag.String("privkey", "server.key", "run as a server")
	flag.Parse()

	// determine whether to run the client or server
	if *serverPtr {
		server.Start(*pubkeyPtr, *privkeyPtr)
	} else {
		fmt.Println("Run as client")
	}
}
