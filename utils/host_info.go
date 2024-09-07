package utils

import (
	"log"
	"net"
)

// MacAddressRecord represents mac address record with interface name and address
type MacAddressRecord struct {
	Name    string
	Address string
}

// MacAddr returns list of MAC addresses of the host
func MacAddr() []MacAddressRecord {
	var alist []MacAddressRecord
	// Get a list of all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error getting network interfaces: ", err)
		return alist
	}

	// Loop over all the interfaces
	for _, iface := range interfaces {
		// Skip interfaces that are down or don't have a hardware (MAC) address
		if iface.Flags&net.FlagUp == 0 || len(iface.HardwareAddr) == 0 {
			continue
		}
		alist = append(alist, MacAddressRecord{Name: iface.Name, Address: iface.HardwareAddr.String()})
	}
	return alist
}

// IpAddresses returns list of IP addresses of the host
func IpAddr() []string {
	var ilist []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("Error getting network interfaces: ", err)
		return ilist
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ilist = append(ilist, ipNet.IP.String())
			}
		}
	}
	return ilist
}
