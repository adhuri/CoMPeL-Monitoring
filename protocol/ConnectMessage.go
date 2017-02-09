package protocol

import "net"

type ConnectRequest struct {
	messageId uint8
	agentIP   net.IP
	agentPort uint16
}

type ConnectReply struct {
	messageId     uint8
	isSuccessfull bool
}
