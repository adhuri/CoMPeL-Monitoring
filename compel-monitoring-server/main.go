package main

import (
	"encoding/binary"
	"fmt"
	"net"

	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func handleConnectMessage(conn net.Conn) {
	defer conn.Close()
	connectMessage := monitorProtocol.ConnectRequest{}
	err := binary.Read(conn, binary.LittleEndian, &connectMessage)
	if err != nil {
		fmt.Println("ERROR : Bad Message From Client" + err.Error())
		return
	} else {
		fmt.Println("INFO: Connect Request Received")
		fmt.Printf("%+v\n", connectMessage)
	}

	connectAck := monitorProtocol.ConnectReply{
		MessageId:     connectMessage.MessageId,
		AgentIP:       connectMessage.AgentIP,
		IsSuccessfull: 1,
	}
	err = binary.Write(conn, binary.LittleEndian, connectAck)
	if err != nil {
		return
	}

}

func main() {

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic("Server Failed to Start")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnectMessage(conn)
	}

}
