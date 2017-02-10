package protocol

import (
	"errors"
	"net"
	"time"
)

type ConnectRequest struct {
	messageId int64
	agentIP   net.IP
	agentPort uint16
}

type ConnectReply struct {
	messageId     int64
	agentIP       net.IP
	isSuccessfull bool
}

func getIPAddressOfHost() (*net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return &ip, nil
		}
	}
	return nil, errors.New("Not Connected To Network")
}

func NewConnectRequest() (*ConnectRequest, error) {
	ip, err := getIPAddressOfHost()
	if err != nil {
		return nil, err
	}
	return &ConnectRequest{
		messageId: time.Now().Unix(),
		agentIP:   *ip,
		agentPort: 6969,
	}, nil
}
