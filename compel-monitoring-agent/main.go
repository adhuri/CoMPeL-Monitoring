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
	// verify if the messageID of request and response is same and even the host
	return connectMessage.AgentIP == ConnectAck.AgentIP && connectMessage.MessageId == ConnectAck.MessageId
}

func generateInitMessage() *monitorProtocol.ConnectRequest {
	// creates a ConnectRequest
	return monitorProtocol.NewConnectRequest()
}

func sendInitMessage(conn net.Conn) error {
	// send init message to server
	connectMessage := *generateInitMessage()
	err := binary.Write(conn, binary.LittleEndian, connectMessage)
	if err != nil {
		// If error occurs in sending a connect message to server then return
		return err
	}

	// read ack from the server
	serverReply := monitorProtocol.ConnectReply{}
	err = binary.Read(conn, binary.LittleEndian, &serverReply)
	if err != nil {
		// If error occurs while reading ACK from server then return
		fmt.Println("ERROR : Bad Reply From Server" + err.Error())
		return err

	} else {
		// Print the ACK received from the server
		fmt.Printf("INFO: Reply Received %+v \n", serverReply)
	}

	// Validate server respose
	var isSucees bool = validateResponse(connectMessage, serverReply)
	if !isSucees {
		// If Ack validation fails then return error
		return errors.New("Invalid Response")
	}

	// If everything goes well return nil error
	return nil
}

func main() {
	// Try connecting to the monitoring server
	// If connection fails try reconnecting after 3 seconds again
	connectedToServer := false
	fmt.Print("Connecting to Server ...")

	// Register client with server
	for !connectedToServer {
		// Setup a TCP connection for communication
		conn, err := net.Dial("tcp", "127.0.0.1:8081")
		if err != nil {
			// Before trying to reconnect to the server wait for 3 seconds
			fmt.Print(".")
			time.Sleep(time.Second * 3)
		} else {
			// If connection successful send a connect message
			err = sendInitMessage(conn)
			if err != nil {
				// Connect Protocol failed midway; Retry
				fmt.Println("\n Try Reconnecting to server")
				conn.Close()
			} else {
				// Client Registration successful
				fmt.Println("\n Connected to Server")
				connectedToServer = true
			}
		}
	}
}
