package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"sync"

	db "github.com/adhuri/Compel-Monitoring/compel-monitoring-server/db"
	model "github.com/adhuri/Compel-Monitoring/compel-monitoring-server/model"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func handleConnectMessage(conn net.Conn, server *model.Server) {
	// When everything is done close the connection
	defer conn.Close()

	// Read the ConnectRequest
	connectMessage := monitorProtocol.ConnectRequest{}
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&connectMessage)
	//err := binary.Read(conn, binary.LittleEndian, &connectMessage)
	if err != nil {
		// If failure in parsing, close the connection and return
		fmt.Println("ERROR : Bad Message From Client" + err.Error())
		return
	} else {
		// If success, print the message received
		fmt.Println("INFO: Connect Request Received")
		fmt.Printf("%+v\n", connectMessage)
	}

	server.UpdateState(connectMessage.AgentIP)

	// Create a ConnectAck Message
	connectAck := monitorProtocol.ConnectReply{
		MessageId:     connectMessage.MessageId,
		AgentIP:       connectMessage.AgentIP,
		IsSuccessfull: 1,
	}

	// Send Connect Ack back to the client
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(connectAck)
	//err = binary.Write(conn, binary.LittleEndian, connectAck)
	if err != nil {
		// If failure in parsing, close the connection and return
		return
	}

}

func tcpListener(wg *sync.WaitGroup, server *model.Server) {
	defer wg.Done()
	// Server listens on all interfaces for TCP connestion
	addr := ":" + server.GetTcpPort()

	listener, err := net.Listen("tcp", addr)
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
		go handleConnectMessage(conn, server)
	}
}

func handleMonitorMessage(conn *net.UDPConn, server *model.Server) {
	var buf [10000]byte

	n, _, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		fmt.Println("Error Reading from UDP socket")
		return
	}
	//fmt.Println(string(buf[0:n]))
	var statsMessage monitorProtocol.StatsMessage
	if err := gob.NewDecoder(bytes.NewReader(buf[0:n])).Decode(&statsMessage); err != nil {
		// handle error
		fmt.Println("Error Decoding at Server")
		return
	}
	// fmt.Printf("%q: {%s,%v}\n", statsMessage.MessageId, utils.IpToString(statsMessage.AgentIP[0:]), statsMessage.Data)
	// fmt.Println(statsMessage.MessageId)
	// fmt.Println(utils.IpToString(statsMessage.AgentIP[0:]))
	fmt.Println(statsMessage.Data)
	// fmt.Println(addr)
	agentIp := statsMessage.AgentIP
	if server.IsAgentConnected(agentIp) {
		// save in the DB
		//statsMessage.Data
		db.StoreData(agentIp.String(), statsMessage.Data)
		//influx.AddPoint(agentIp.String(), containerId, cpuUsage, memoryUsage, timestamp)
		fmt.Println("Valid Agent : ")
		server.UpdateState(agentIp)
	}
	//conn.WriteToUDP([]byte("Hello Client"), addr)

}

func udpListener(wg *sync.WaitGroup, server *model.Server) {
	defer wg.Done()

	addr := ":" + server.GetUdpPort()

	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		fmt.Println("Error in Resolving Address " + err.Error())
		panic("Unable to Start UDP Service on server")
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic("Unable to Start UDP Service on server")
	}

	for {
		handleMonitorMessage(conn, server)
	}

}

func main() {

	serverUdpPort := flag.String("udpport", "7071", "udp port on the server")
	serverTcpPort := flag.String("tcpport", "8081", "tcp port of the server")

	flag.Parse()

	server := model.NewServer(*serverTcpPort, *serverUdpPort)
	var wg sync.WaitGroup
	wg.Add(2)

	go tcpListener(&wg, server)
	go udpListener(&wg, server)

	wg.Wait()

}
