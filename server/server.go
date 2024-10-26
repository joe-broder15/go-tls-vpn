package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

// extract ip addresses from ipv4 packet
func extractIPs(packet []byte) (net.IP, net.IP, error) {
	// Check that the packet is at least as long as the minimum IPv4 header size
	if len(packet) < 20 {
		return nil, nil, fmt.Errorf("packet too short")
	}

	// The Internet Header Length (IHL) is in the lower 4 bits of the first byte
	ihl := packet[0] & 0x0F
	ihl *= 4 // Convert IHL to bytes

	// Ensure the packet length matches the header length
	if len(packet) < int(ihl) {
		return nil, nil, fmt.Errorf("packet too short for header length")
	}

	// Extract source IP (bytes 12-15)
	srcIP := net.IPv4(packet[12], packet[13], packet[14], packet[15])

	// Extract destination IP (bytes 16-19)
	dstIP := net.IPv4(packet[16], packet[17], packet[18], packet[19])

	return srcIP, dstIP, nil
}

// function to handle a connection with a client
func handleConnection(c net.Conn, addressPool *AddressPool, clientTrafficChannel chan []byte, routingTable *RoutingTable) {

	// get client ip address from the address pool
	clientAddress, err := addressPool.GetUnusedAddress()
	if err != nil {
		log.Println(err)
		return
	}
	defer addressPool.ReleaseAddress(clientAddress)

	// send ip address to the client
	c.Write([]byte(clientAddress.String()))

	// add a routing entry for the new connection
	err = routingTable.AddEntry(clientAddress.String(), c)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer routingTable.RemoveEntry(clientAddress.String())

	// listen for traffic from the client
	for {
		// allocate buffer for reading
		buffer := make([]byte, 1500)

		// read from the connection
		n, err := c.Read(buffer)
		if err != nil {
			log.Println("error reading from client", clientAddress, err)
			return
		}
		// write to the interface
		clientTrafficChannel <- buffer[:n]
	}
}

// function to handle a connection with a client
func listenForConnections(listener net.Listener, addressPool *AddressPool, clientTrafficChannel chan []byte, routingTable *RoutingTable) {
	// listen for incoming connections
	for {
		// accept connections and hand them off to a new goroutine
		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println("got connection from client")
		go handleConnection(conn, addressPool, clientTrafficChannel, routingTable)
	}
}

func readFromTun(ifce *water.Interface, serverTrafficChannel chan []byte) {

	// listen for traffic from the client
	for {
		// allocate memory for a packet
		buffer := make([]byte, 1500)

		// read from the tun device
		n, err := ifce.Read(buffer)
		if err != nil {
			log.Println(err)
			continue
		}

		serverTrafficChannel <- buffer[:n]
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

	// print info
	log.Println("Server ip:", serverAddress)
	log.Println("Interface:", ifce.Name())

	// create the routing table
	routingTable := NewRoutingTable()

	// create the channels for orchestration
	clientTrafficChannel := make(chan []byte, 1000)
	serverTrafficChannel := make(chan []byte, 1000)

	// start thread that reads from tun device
	go readFromTun(ifce, serverTrafficChannel)

	// start primary server thread that will facilitate client connection
	go listenForConnections(listener, addressPool, clientTrafficChannel, routingTable)

	// select on the channels
	for {
		select {
		case clientPacket := <-clientTrafficChannel:
			ifce.Write(clientPacket)
		case serverPacket := <-serverTrafficChannel:
			if (serverPacket[0] >> 4) != 4 {
				continue
			} else {
				_, dstIp, err := extractIPs(serverPacket)
				if err != nil {
					log.Println(err)
					continue
				}
				routingEntry, err := routingTable.GetEntry(dstIp.String())
				if err != nil {
					log.Println(err)
					continue
				}
				_, err routingEntry.conn.Write(serverPacket)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}
