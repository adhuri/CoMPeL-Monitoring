package protocol

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"time"

	logrus "github.com/Sirupsen/logrus"
)

func sendInitMessage(conn net.Conn) error {
	// send init message to server
	connectMessage := *NewConnectRequest()
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(connectMessage)
	//err := binary.Write(conn, binary.LittleEndian, connectMessage)
	if err != nil {
		// If error occurs in sending a connect message to server then return
		fmt.Printf("ERROR : Failure While Sending Data To Server " + err.Error())
		return err
	}

	// read ack from the server
	serverReply := ConnectReply{}
	decoder := gob.NewDecoder(conn)
	err = decoder.Decode(&serverReply)
	// err = binary.Read(conn, binary.LittleEndian, &serverReply)
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

func ConnectToServer(serverIp, tcpPort string, log *logrus.Logger) {
	// Try connecting to the monitoring server
	// If connection fails try reconnecting after 3 seconds again
	connectedToServer := false
	log.Info("Connecting to Server ...")
	//fmt.Print("Connecting to Server ...")

	// Register client with server
	for !connectedToServer {
		// Setup a TCP connection for communication
		addr := serverIp + ":" + tcpPort

		conn, err := net.Dial("tcp", addr)
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

func SendContainerStatistics(stringToSend []ContainerStats, serverIp string, udpPort string) {
	udpAddr, err := net.ResolveUDPAddr("udp4", serverIp+":"+udpPort)
	if err != nil {
		fmt.Println("Error in Resolving Address " + err.Error())
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error in Resolving Address")
		return
	}

	// construct a StatsMessage
	//var buffer bytes.Buffer
	//data := ""
	startPointer := 0
	endPointer := 0
	totalSizeOfData := 0
	for i := 0; i < len(stringToSend); i++ {
		//buffer.WriteString(stringToSend[i])
		totalSizeOfData += stringToSend[i].Size()
		if totalSizeOfData <= 576 {
			endPointer++
			fmt.Println(endPointer)
		} else {
			statsMessage := *NewStatsMessage(stringToSend[startPointer:endPointer])
			var buf bytes.Buffer
			if err := gob.NewEncoder(&buf).Encode(statsMessage); err != nil {
				// handle error
				fmt.Println("Error in Encoding the StatMessage using GOB Encoder")
				return
			}
			_, err := conn.Write(buf.Bytes())
			//
			// // Send StatsMessage
			// _, err = conn.Write([]byte(stringToSend))
			if err != nil {
				fmt.Println("Error in Resolving Address")
				return
			}
			startPointer = endPointer
			totalSizeOfData = 0
			totalSizeOfData += stringToSend[i].Size()
			//buffer.Reset()
			//buffer.WriteString(stringToSend[i])
			endPointer++
		}
	}

	if startPointer < len(stringToSend) {

		statsMessage := *NewStatsMessage(stringToSend[startPointer:])
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(statsMessage); err != nil {
			// handle error
			fmt.Println("Error in Encoding the StatMessage using GOB Encoder")
			return
		}
		_, err := conn.Write(buf.Bytes())
		//
		// // Send StatsMessage
		// _, err = conn.Write([]byte(stringToSend))
		if err != nil {
			fmt.Println("Error in Resolving Address")
			return
		}

	}

	// var buf [512]byte
	// n, err := conn.Read(buf[0:])
	// if err != nil {
	// 	fmt.Println("Error")
	// }
	//
	// fmt.Println(string(buf[0:n]))

}
