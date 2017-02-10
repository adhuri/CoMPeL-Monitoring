package protocol

import (
	"errors"
	"net"
	"time"
)

type ConnectRequest struct {
	//AgentIP   net.IP
	MessageId int64
	AgentIP   [4]byte
	AgentPort uint16
}

type ConnectReply struct {
	//AgentIP   net.IP
	MessageId     int64
	AgentIP       [4]byte
	IsSuccessfull uint8
}

func getIPAddressOfHost(hostIP []byte) error {
	// get all Interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			// interfae down
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			// loopback interface
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			var ip net.IP
			// check the type
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				// If Loopback IP i.e. addresses like 127.*.*.*
				continue
			}
			// convert address to 4-byte form
			ip = ip.To4()
			if ip == nil {
				// not a valid IPv4 address
				continue
			}
			for i, val := range ip {
				hostIP[i] = val
			}
			return nil
		}
	}
	return errors.New("Not Connected To Network")
}

func NewConnectRequest() *ConnectRequest {
	var hostIP [4]byte
	err := getIPAddressOfHost(hostIP[0:])
	// If external IP of the host is not found then return err
	if err != nil {
		panic("Error Fetching Valid IP Address")
	}
	return &ConnectRequest{
		MessageId: time.Now().Unix(),
		AgentIP:   hostIP,
		AgentPort: 6969,
	}
}
