package protocol

import (
	"net"
	"time"

	utils "github.com/adhuri/Compel-Monitoring/utils"
)

type ConnectRequest struct {
	MessageId int64
	AgentIP   net.IP
	AgentPort uint16
}

type ConnectReply struct {
	MessageId     int64
	AgentIP       net.IP
	IsSuccessfull uint8
}

func ValidateResponse(connectMessage ConnectRequest, ConnectAck ConnectReply) bool {
	// validate the response from the server
	// verify if the messageID of request and response is same and even the host
	return utils.CheckIPAddressesEqual(connectMessage.AgentIP, ConnectAck.AgentIP) && connectMessage.MessageId == ConnectAck.MessageId && ConnectAck.IsSuccessfull == 1
}

func NewConnectRequest() *ConnectRequest {
	// Get External IP of host
	//var hostIP [4]byte
	hostIP, err := utils.GetIPAddressOfHost()
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
