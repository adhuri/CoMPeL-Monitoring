package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func validateResponse(connectMessage monitorProtocol.ConnectRequest, ConnectAck monitorProtocol.ConnectReply) bool {
	// validate the response from the server
	return connectMessage.AgentIP == ConnectAck.AgentIP && connectMessage.MessageId == ConnectAck.MessageId
}

func generateInitMessage() *monitorProtocol.ConnectRequest {
	return monitorProtocol.NewConnectRequest()
}

func sendInitMessage(conn net.Conn) error {
	//send init message to server
	connectMessage := *generateInitMessage()
	err := binary.Write(conn, binary.LittleEndian, connectMessage)
	if err != nil {
		return err
	}

	// read ack from the server
	serverReply := monitorProtocol.ConnectReply{}
	err = binary.Read(conn, binary.LittleEndian, &serverReply)
	if err != nil {
		fmt.Println("ERROR : Bad Reply From Server" + err.Error())
		return err
	} else {
		fmt.Println("INFO: Reply Received")
	}

	//validate server respose
	var isSucees bool = validateResponse(connectMessage, serverReply)
	if !isSucees {
		return errors.New("Invalid Response")
	}
	return nil
}

func main() {
	// Try connecting to the monitoring server
	// If connection fails try reconnecting after 3 seconds again
	connectedToServer := false
	fmt.Print("Connecting to Server ...")
	for !connectedToServer {
		conn, err := net.Dial("tcp", "127.0.0.1:8081")
		if err != nil {
			fmt.Print(".")
			time.Sleep(time.Second * 3)
		} else {
			// If connection successful send a connect message
			err = sendInitMessage(conn)
			if err != nil {
				fmt.Println("\n Try Reconnecting to server")
				conn.Close()
			} else {
				fmt.Println("\n Connected to Server")
				connectedToServer = true
			}
		}
	}
}
