package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
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

func connectToServer(done chan bool) {
	defer close(done)
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
				defer conn.Close()
			} else {
				// Client Registration successful
				fmt.Println("\n Connected to Server")
				connectedToServer = true
			}
		}
	}
}

func sendMonitorPackets(stringToSend string) {
	udpAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:8081")
	if err != nil {
		fmt.Println("Error in Resolving Address " + err.Error())
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error in Resolving Address")
		return
	}

	_, err = conn.Write([]byte(stringToSend))
	if err != nil {
		fmt.Println("Error in Resolving Address")
		return
	}

	var buf [512]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		fmt.Println("Error")
	}

	fmt.Println(string(buf[0:n]))

}

func numberOfContaiers() int {
	return 4
}

func worker(index int, containerStats chan string) {
	time.Sleep(time.Second * 3)

	stats := " # Container " + strconv.Itoa(index)
	fmt.Println(stats)
	containerStats <- stats
}

func main() {
	done := make(chan bool)
	go connectToServer(done)
	<-done

	numOfWorkers := numberOfContaiers()

	containerStats := make(chan string, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		go worker(i, containerStats)
	}

	var buffer bytes.Buffer

	for i := 0; i < numOfWorkers; i++ {
		buffer.WriteString(<-containerStats)
	}

	stringToSend := buffer.String()

	fmt.Println(stringToSend)
	sendMonitorPackets(stringToSend)

}
