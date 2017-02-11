package main

import (
	"encoding/binary"
	"fmt"
	"net"

	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func main() {

	ln, _ := net.Listen("tcp", ":8081")
	conn, _ := ln.Accept()

	connectMessage := monitorProtocol.ConnectRequest{}
	err := binary.Read(conn, binary.LittleEndian, &connectMessage)
	if err != nil {
		fmt.Println("ERROR : Bad Message From Client" + err.Error())
	} else {
		fmt.Println("INFO: Connect Request Received")
		fmt.Printf("%+v\n", connectMessage)
	}

	connectAck := monitorProtocol.ConnectReply{
		MessageId:     connectMessage.MessageId,
		AgentIP:       connectMessage.AgentIP,
		IsSuccessfull: 1,
	}
	binary.Write(conn, binary.LittleEndian, connectAck)

}
