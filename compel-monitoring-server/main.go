package main

import "net"
import "fmt"
import "bufio"
import "os"

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println("Connection failed")
	} else {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Text to send: ")
			text, _ := reader.ReadString('\n')
			// send to socket
			fmt.Fprintf(conn, text+"\n")
			// listen for reply
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Print("Message from server: " + message)
		}
	}
}
