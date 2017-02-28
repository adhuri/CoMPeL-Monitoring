package protocol

import (
	"time"

	"github.com/adhuri/Compel-Monitoring/utils"
)

type ServerHeartBeat struct {
	MessageId  int64
	ServerIP   [4]byte
	ServerPort uint16
}

func NewServerHeartBeat() *ServerHeartBeat {
	// Get External IP of host
	var hostIP [4]byte
	err := utils.GetIPAddressOfHost(hostIP[0:])

	// If external IP of the host is not found then panic
	if err != nil {
		panic("Error Fetching Valid IP Address")
	}

	// Generate a Connect Request
	return &ServerHeartBeat{
		MessageId:  time.Now().Unix(),
		ServerIP:   hostIP,
		ServerPort: 6969,
	}
}
