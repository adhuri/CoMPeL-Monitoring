package main

import "net"
import "fmt"
import "bufio"
import "time"
import github.com/adhuri/Compel-Monitoring/protocol "protocol"

func validateResponse(ackReply string, connectMessage string) bool {
	// validate the response from the server
	return true
}

func generateInitMessage() string {
	return "1\n"
}
func sendInitMessage(conn net.Conn) {
	//send init message to server
	connectMessage := generateInitMessage()
	fmt.Fprintf(conn, connectMessage)
	// read ack from the server
	serverReply, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + serverReply)
	validateResponse(serverReply, connectMessage)
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
