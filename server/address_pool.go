package server

import (
	"fmt"
	"net"
	"sync"
)

// AddressPool tracks IP addresses assigned from a subnet
type AddressPool struct {
	ipNet        *net.IPNet
	availableIPs []net.IP
	usedIPs      map[string]struct{}
	mu           sync.Mutex
}

// NewAddressPool creates a new AddressPool from a CIDR notation
func NewAddressPool(cidr string) (*AddressPool, error) {
	// Parse the CIDR string into an IP address and a network.
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation: %v", err)
	}

	// Initialize the AddressPool
	pool := &AddressPool{
		ipNet:   ipnet,
		usedIPs: make(map[string]struct{}),
	}

	// Generate all IP addresses in the subnet
	var ips []net.IP
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ipCopy := copyIP(ip)
		ips = append(ips, ipCopy)
	}

	// Remove network address and broadcast address
	if len(ips) < 2 {
		return nil, fmt.Errorf("CIDR block is too small")
	}
	ips = ips[1 : len(ips)-1] // Remove first and last IPs

	pool.availableIPs = ips

	return pool, nil
}

// GetUnusedAddress returns any free address from the pool
func (pool *AddressPool) GetUnusedAddress() (net.IP, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if len(pool.availableIPs) == 0 {
		return nil, fmt.Errorf("no available IP addresses")
	}

	// Pop an IP from availableIPs
	ip := pool.availableIPs[0]
	pool.availableIPs = pool.availableIPs[1:]

	// Add to usedIPs
	pool.usedIPs[ip.String()] = struct{}{}

	return ip, nil
}

// ReleaseAddress releases an IP address back to the pool
func (pool *AddressPool) ReleaseAddress(addr net.IP) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	ipStr := addr.String()
	if _, exists := pool.usedIPs[ipStr]; exists {
		delete(pool.usedIPs, ipStr)
		pool.availableIPs = append(pool.availableIPs, addr)
	}
}

// incrementIP increments the given IP address by one.
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
}

// copyIP creates a copy of the given IP address.
func copyIP(ip net.IP) net.IP {
	newIP := make(net.IP, len(ip))
	copy(newIP, ip)
	return newIP
}
