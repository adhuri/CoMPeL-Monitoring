package protocol

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net"
	"time"

	logrus "github.com/Sirupsen/logrus"
)

func sendInitMessage(conn net.Conn, log *logrus.Logger) error {
	// send init message to server
	connectMessage := *NewConnectRequest()
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(connectMessage)
	if err != nil {
		// If error occurs in sending a connect message to server then return
		log.Errorln("Failure While Sending Data To Server " + err.Error())
		return err
	}
	log.Infoln("Connect Message Successfully Sent")

	// read ack from the server
	serverReply := ConnectReply{}
	decoder := gob.NewDecoder(conn)
	err = decoder.Decode(&serverReply)
	if err != nil {
		// If error occurs while reading ACK from server then return
		log.Errorln("Bad Reply From Server " + err.Error())
		return err

	} else {
		// Print the ACK received from the server
		log.Infoln("Connect ACK Received")
	}

	// Validate server respose
	var isSucees bool = ValidateResponse(connectMessage, serverReply)
	if !isSucees {
		// If Ack validation fails then return error
		log.Errorln("Invalid Connect ACK")
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

	// Register client with server
	for !connectedToServer {
		// Setup a TCP connection for communication
		addr := serverIp + ":" + tcpPort

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			// Before trying to reconnect to the server wait for 3 seconds
			log.Warn("Server Not Alive")
			time.Sleep(time.Second * 3)
		} else {
			// If connection successful send a connect message
			err = sendInitMessage(conn, log)
			if err != nil {
				// Connect Protocol failed midway; Retry
				log.Warn("Connect Protocol failed. Try Reconnecting to server")
				defer conn.Close()
			} else {
				// Client Registration successful
				log.Infoln("Connected to Server")
				connectedToServer = true
			}
		}
	}
}

func SendContainerStatistics(stringToSend []ContainerStats, serverIp string, udpPort string, log *logrus.Logger) {
	udpAddr, err := net.ResolveUDPAddr("udp4", serverIp+":"+udpPort)
	if err != nil {
		log.Errorln("Error in Resolving Address " + err.Error())
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Errorln("Error in Dialing UDP" + err.Error())
		return
	}

	// construct a StatsMessage
	startPointer := 0
	endPointer := 0
	totalSizeOfData := 0
	for i := 0; i < len(stringToSend); i++ {
		totalSizeOfData += stringToSend[i].Size()
		if totalSizeOfData <= 576 {
			// Increment End Pointer if we can add more
			endPointer++

			// Debug Logging
			log.WithFields(logrus.Fields{
				"Total_Data_To_Send":  len(stringToSend),
				"Current_Packet_Size": totalSizeOfData,
			}).Debugln("End Pointer Value : " + string(endPointer))

		} else {
			statsMessage := *NewStatsMessage(stringToSend[startPointer:endPointer])
			var buf bytes.Buffer
			if err := gob.NewEncoder(&buf).Encode(statsMessage); err != nil {
				log.Errorln("Error in Encoding the StatMessage using GOB Encoder")
				return
			}
			_, err := conn.Write(buf.Bytes())
			if err != nil {
				log.Errorln("Error in Sending Stat Message")
				return
			}
			log.Debugln("Stats Message = ", statsMessage)
			log.Infoln("Stats Message Sent.\n")

			// Debug Logging
			log.WithFields(logrus.Fields{
				"Start_Pointer":       startPointer,
				"End_Pointer":         endPointer,
				"Total_Data_To_Send":  len(stringToSend),
				"Current_Packet_Size": totalSizeOfData,
			}).Debugln("Information Regarding Stat Message Sent")

			// Update Pointers
			startPointer = endPointer
			totalSizeOfData = 0
			totalSizeOfData += stringToSend[i].Size()
			endPointer++
		}
	}

	if startPointer < len(stringToSend) {

		statsMessage := *NewStatsMessage(stringToSend[startPointer:])
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(statsMessage); err != nil {
			log.Errorln("Error in Encoding the StatMessage using GOB Encoder")
			return
		}
		_, err := conn.Write(buf.Bytes())
		if err != nil {
			log.Errorln("Error in Sending Stat Message")
			return
		}
		log.Debugln("Stats Message = ", statsMessage)
		log.Infoln("Stats Message Sent \n")

	}
}
