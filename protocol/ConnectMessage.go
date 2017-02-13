package protocol

import "time"
import utils "github.com/adhuri/Compel-Monitoring/utils"

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

func ValidateResponse(connectMessage ConnectRequest, ConnectAck ConnectReply) bool {
	// validate the response from the server
	// verify if the messageID of request and response is same and even the host
	return connectMessage.AgentIP == ConnectAck.AgentIP && connectMessage.MessageId == ConnectAck.MessageId
}

func NewConnectRequest() *ConnectRequest {
	// Get External IP of host
	var hostIP [4]byte
	err := utils.GetIPAddressOfHost(hostIP[0:])
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
