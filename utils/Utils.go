package utils

import (
	"errors"
	"net"
	"time"
	"log"
)

func GetIPAddressOfHost(hostIP []byte) error {
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

func IpToString(hostIP []byte) string {
	return string(hostIP[0]) + string(hostIP[1]) + string(hostIP[2]) + string(hostIP[3])
}
// Time any function in the repository - 
// Usage - defer utils.TimeTrack(time.Now(), "Filename.go-FunctionName")
func TimeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
        log.Printf("%s took %s", name, elapsed)
}
