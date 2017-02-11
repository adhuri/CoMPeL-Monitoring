package protocol

import (
	"errors"
	"net"
	"time"
)

type ConnectRequest struct {
	MessageId int64
	AgentIP   [4]byte
	AgentPort uint16
}

type ConnectReply struct {
	MessageId     int64
	AgentIP       [4]byte
	IsSuccessfull uint8
}

func getIPAddressOfHost(hostIP []byte) error {
	// Get all Interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	// Iterate over interface to find the right interface
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			// Interfae down
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			// Loopback interface
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			// If the interface is up but ip has not been set
			return err
		}

		// iterate over interface addresses
		for _, addr := range addrs {
			var ip net.IP
			// check the type
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// If Loopback IP i.e. addresses like 127.*.*.*
			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Convert address to 4-byte form
			ip = ip.To4()
			if ip == nil {
				continue
			}

			// Copy the 4 bytes of IP to the slice passed as argument
			for i, val := range ip {
				hostIP[i] = val
			}
			return nil
		}
	}

	// If no interface is connected to the network
	return errors.New("Not Connected To Network")
}

func NewConnectRequest() *ConnectRequest {
	// Get External IP of host
	var hostIP [4]byte
	err := getIPAddressOfHost(hostIP[0:])
	// If external IP of the host is not found then panic
	if err != nil {
		panic("Error Fetching Valid IP Address")
	}

	// Generate a Connect Request
	return &ConnectRequest{
		MessageId: time.Now().Unix(),
		AgentIP:   hostIP,
		AgentPort: 6969,
	}
}
