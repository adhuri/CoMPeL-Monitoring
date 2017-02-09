package main

import "net"
import "fmt"
import "bufio"

func main() {

	ln, _ := net.Listen("tcp", ":8081")
	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Received:", string(message))
		newmessage := "2"
		conn.Write([]byte(newmessage + "\n"))
	}
}
