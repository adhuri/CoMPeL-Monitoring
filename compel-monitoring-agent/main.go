package main

import "net"
import "fmt"

import "time"
import "encoding/binary"
import monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"

func validateResponse(connectMessage monitorProtocol.ConnectRequest, ack monitorProtocol.ConnectReply) bool {
	// validate the response from the server
	return true
}

func generateInitMessage() *monitorProtocol.ConnectRequest {
	return monitorProtocol.NewConnectRequest()
}

func sendInitMessage(conn net.Conn) {
	//send init message to server
	connectMessage := *generateInitMessage()
	binary.Write(conn, binary.LittleEndian, connectMessage)
	//fmt.Fprintf(conn, connectMessage)
	// read ack from the server
	serverReply := monitorProtocol.ConnectReply{}
	err := binary.Read(conn, binary.LittleEndian, &serverReply)
	if err != nil {
		fmt.Println("ERROR : Bad Reply From Server" + err.Error())
	} else {
		fmt.Println("INFO: Reply Received")
	}
	//serverReply, _ := bufio.NewReader(conn).ReadString('\n')
	//fmt.Print("Message from server: " + serverReply)
	validateResponse(connectMessage, serverReply)
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
			sendInitMessage(conn)
			fmt.Println("\n Connected to Server")
			connectedToServer = true
		}
	}
}
