package main

import (
	"flag"
	"go-tls-vpn/client"
	"go-tls-vpn/server"
)

func main() {
	// set the command line flags
	serverPtr := flag.Bool("s", false, "run as a server")
	pubkeyPtr := flag.String("pubkey", "server.crt", "run as a server")
	privkeyPtr := flag.String("privkey", "server.key", "run as a server")
	portPtr := flag.String("p", "8888", "server port")
	ipPtr := flag.String("ip", "127.0.0.1", "ip address of server")
	vNet := flag.String("vnet", "10.0.0.0/24", "cidr of the virtual network")

	flag.Parse()

	// determine whether to run the client or server
	if *serverPtr {
		server.Start(*pubkeyPtr, *privkeyPtr, *portPtr, *vNet)
	} else {
		client.Start(*ipPtr, *portPtr)
	}
}
