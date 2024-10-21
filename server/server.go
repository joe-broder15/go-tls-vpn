package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

// function to handle a connection with a client
func handleConnection(c net.Conn, addressPool * AddressPool, ifce * water.Interface, routingTable * RoutingTable) {

	readBuf := make([]byte, 1500)
	
	// get client ip address from the address pool
	clientAddress, err := addressPool.GetUnusedAddress()
	if err != nil {
		log.Println(err);
		return
	}
	defer addressPool.ReleaseAddress(clientAddress)

	// add a routing entry for the new connection
	err = routingTable.AddEntry(clientAddress.String(), c)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer routingTable.RemoveEntry(clientAddress.String())

	
	// send ip address to the client
	c.Write([]byte(clientAddress.String()))

	// listen for traffic from the client
	for {
		// read from the connection
		n, err := c.Read(readBuf)
		if err != nil {
			log.Fatal(err)
		}
		// write to the interface
		_, err = ifce.Write(readBuf[:n])
		if err != nil {
			log.Fatal(err)
		}
		// clean the buffer
		for i := 0 ; i < n; i++ {
			readBuf[i] = 0x00
		}
	}
}

func handleTun(ifce * water.Interface, routingTable * RoutingTable) {
	readBuf := make([]byte, 1500)
	// listen for traffic from the client
	for {
		// read from the tun device
		n, err := ifce.Read(readBuf)
		if err != nil {
			log.Println(err)
			continue
		}

		// get 
		fmt.Println(readBuf)

		// // write to the interface
		// _, err = ifce.Write(readBuf[:n])
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// clean the buffer
		for i := 0 ; i < n; i++ {
			readBuf[i] = 0x00
		}
	}
}

// start the server
func Start(certFile string, keyFile string, port string, vNet string) {
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

	// create the ip address pool
	addressPool, err := NewAddressPool(vNet)
	if err != nil {
		log.Fatal(err)
	}

	// get server ip address
	serverAddress, err := addressPool.GetUnusedAddress()
	if err != nil {
		log.Fatal(err)
	}

	// create tun interface
	water_config := water.Config{
		DeviceType: water.TUN,
	}

	ifce, err := water.New(water_config)
	if err != nil {
		log.Fatal(err)
	}

	// assign ip to tun interface
	ipAssignCmd := exec.Command("sudo", "ip", "addr", "add", serverAddress.String()+"/24", "dev", "tun0")
	if err := ipAssignCmd.Run(); err != nil {
		log.Fatal(err)
	}

	// bring up tun device
	interfaceUpCmd := exec.Command("sudo", "ip", "link", "set", "tun0", "up")
	if err := interfaceUpCmd.Run(); err != nil {
		log.Fatal(err)
	}

	// create the routing table
	routingTable := NewRoutingTable()

	go handleTun(ifce, routingTable)

	fmt.Println("server ip will be", serverAddress)
	fmt.Printf("Interface Name: %s\n", ifce.Name())

	// listen for incoming connections
	for {
		// accept connections and hand them off to a new goroutine
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("got connection from client")
		go handleConnection(conn, addressPool, ifce, routingTable)
	}

}
