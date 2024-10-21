package server

import (
	"fmt"
	"net"
	"sync"
)

type RoutingTableEntry struct {
	ip string
	conn net.Conn
}

type RoutingTable struct {
	mu sync.Mutex
	entries map[string]*RoutingTableEntry
}

func (table *RoutingTable) GetEntry(ip string) (*RoutingTableEntry, error) {
	table.mu.Lock()
	defer table.mu.Unlock()
	item := table.entries[ip]
	if item == nil {
		return nil, fmt.Errorf("no such entry in routing table for %s", ip)
	} else {
		return item, nil
	}
}

func (table *RoutingTable) AddEntry(ip string, conn net.Conn) (error) {
	table.mu.Lock()
	defer table.mu.Unlock()
	
	// create a routing table entry
	newRoutingEntry := new(RoutingTableEntry)
	newRoutingEntry.conn = conn
	newRoutingEntry.ip = ip

	if table.entries[ip] != nil {
		return fmt.Errorf("routing entry already exists for %s", ip)
	} 

	table.entries[ip] = newRoutingEntry

	return nil

}

func (table *RoutingTable) RemoveEntry(ip string) (error) {
	table.mu.Lock()
	defer table.mu.Unlock()
	
	delete(table.entries, ip)

	return nil
}

func NewRoutingTable() (*RoutingTable) {
	routingTable := new(RoutingTable)
	routingTable.entries = make(map[string]*RoutingTableEntry)
	return routingTable	
}