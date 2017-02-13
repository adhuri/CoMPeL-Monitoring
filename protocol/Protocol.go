package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

func sendInitMessage(conn net.Conn) error {
	// send init message to server
	connectMessage := *NewConnectRequest()
	err := binary.Write(conn, binary.LittleEndian, connectMessage)
	if err != nil {
		// If error occurs in sending a connect message to server then return
		return err
	}

	// read ack from the server
	serverReply := ConnectReply{}
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
	var isSucees bool = ValidateResponse(connectMessage, serverReply)
	if !isSucees {
		// If Ack validation fails then return error
		return errors.New("Invalid Response")
	}

	// If everything goes well return nil error
	return nil
}

func ConnectToServer(done chan bool) {
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

func SendContainerStatistics(stringToSend string) {
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
