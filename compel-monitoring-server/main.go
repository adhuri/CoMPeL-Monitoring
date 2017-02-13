package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func handleConnectMessage(conn net.Conn) {
	// When everything is done close the connection
	defer conn.Close()

	// Read the ConnectRequest
	connectMessage := monitorProtocol.ConnectRequest{}
	err := binary.Read(conn, binary.LittleEndian, &connectMessage)
	if err != nil {
		// If failure in parsing, close the connection and return
		fmt.Println("ERROR : Bad Message From Client" + err.Error())
		return
	} else {
		// If success, print the message received
		fmt.Println("INFO: Connect Request Received")
		fmt.Printf("%+v\n", connectMessage)
	}

	// Create a ConnectAck Message
	connectAck := monitorProtocol.ConnectReply{
		MessageId:     connectMessage.MessageId,
		AgentIP:       connectMessage.AgentIP,
		IsSuccessfull: 1,
	}

	// Send Connect Ack back to the client
	err = binary.Write(conn, binary.LittleEndian, connectAck)
	if err != nil {
		// If failure in parsing, close the connection and return
		return
	}

}

func tcpListener(wg *sync.WaitGroup) {
	defer wg.Done()
	// Server listens on all interfaces for TCP connestion
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic("Server Failed to Start")
	}

	// Wait for clients to connect
	for {
		// Accept a connection and spin-off a goroutine
		conn, err := listener.Accept()
		if err != nil {
			// If error continue to wait for other clients to connect
			continue
		}
		go handleConnectMessage(conn)
	}
}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)
	go tcpListener(&wg)
	wg.Wait()

}
