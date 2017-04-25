package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	logrus "github.com/Sirupsen/logrus"
	db "github.com/adhuri/Compel-Monitoring/compel-monitoring-server/db"
	model "github.com/adhuri/Compel-Monitoring/compel-monitoring-server/model"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
	"github.com/gorilla/mux"
	"github.com/mitchellh/hashstructure"
)

var (
	log    *logrus.Logger
	server *model.Server
)

func init() {

	log = logrus.New()

	// Output logging to stdout
	log.Out = os.Stdout

	// Only log the info severity or above.
	log.Level = logrus.InfoLevel

	// Microseconds level logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	customFormatter.FullTimestamp = true

	log.Formatter = customFormatter

}

func handleConnectMessage(conn net.Conn, server *model.Server) {
	// When everything is done close the connection
	// defer conn.Close()

	// Read the ConnectRequest
	connectMessage := monitorProtocol.ConnectRequest{}
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&connectMessage)
	//err := binary.Read(conn, binary.LittleEndian, &connectMessage)
	if err != nil {
		// If failure in parsing, close the connection and return
		log.Errorln("Bad Message From Client" + err.Error())
		return
	} else {
		// If success, print the message received
		log.Infoln("Connect Request Received")
		log.Debugln("Connect Request Content : ", connectMessage)
	}

	if server.IsAgentConnected(connectMessage.AgentIP) {
		server.UpdateState(connectMessage.AgentIP)
	} else {
		server.UpdateState(connectMessage.AgentIP)
		server.UpdateStatsMap(connectMessage.AgentIP.String(), time.Now())
	}

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
		log.Errorln("Connect Ack Failed")
		return
	}
	log.Infoln("Connect Ack Sent")

}

func tcpListener(wg *sync.WaitGroup, server *model.Server) {
	defer wg.Done()
	// Server listens on all interfaces for TCP connestion
	addr := ":" + server.GetTcpPort()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Server Failed To Start ")
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
		log.Errorln("Error Reading from UDP socket")
		return
	}
	//fmt.Println(string(buf[0:n]))
	var statsMessage monitorProtocol.StatsMessage
	if err := gob.NewDecoder(bytes.NewReader(buf[0:n])).Decode(&statsMessage); err != nil {
		// handle error
		log.Errorln("Error Decoding at Server")
		return
	}
	log.Infoln("Stats Message Received From Agent " + statsMessage.AgentIP.String())
	// fmt.Printf("%q: {%s,%v}\n", statsMessage.MessageId, utils.IpToString(statsMessage.AgentIP[0:]), statsMessage.Data)
	// fmt.Println(statsMessage.MessageId)
	// fmt.Println(utils.IpToString(statsMessage.AgentIP[0:]))
	log.Debugln(statsMessage.Data)
	// fmt.Println(addr)
	agentIp := statsMessage.AgentIP
	if server.IsAgentConnected(agentIp) {
		// save in the DB
		//statsMessage.Data
		hash, err := hashstructure.Hash(statsMessage.Data, nil)
		if err != nil {
			log.Errorln("Hash Calculation Failed")
			return
		}

		if hash != statsMessage.HashCode {
			log.Errorln("Hash Didn't Match")
			return
		}

		containerList := make([]string, 0)
		for _, containerStat := range statsMessage.Data {
			containerId := containerStat.ContainerID
			containerList = append(containerList, containerId)
		}

		db.StoreData(agentIp.String(), statsMessage.Data, server, log)
		//influx.AddPoint(agentIp.String(), containerId, cpuUsage, memoryUsage, timestamp)
		log.Infoln("Agent " + agentIp.String() + " Validated")
		server.UpdateState(agentIp)
		server.SetActiveContainersForAgent(agentIp.String(), containerList)
		server.IncrementPacketReceivedCounter()
	}
	//conn.WriteToUDP([]byte("Hello Client"), addr)

}

func udpListener(wg *sync.WaitGroup, server *model.Server) {
	defer wg.Done()

	addr := ":" + server.GetUdpPort()

	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		log.Fatalln("Error in Resolving Address " + err.Error())
		// panic("Unable to Start UDP Service on server")
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalln("Unable to Start UDP Service on server")
		// panic("Unable to Start UDP Service on server")
	}

	for {
		handleMonitorMessage(conn, server)
	}

}

func HandleQuery(w http.ResponseWriter, req *http.Request) {
	queryResponse := monitorProtocol.GenerateQueryResponse(server)
	json.NewEncoder(w).Encode(queryResponse)
}

func predictionQueryListener(wg *sync.WaitGroup, server *model.Server) {
	defer wg.Done()
	router := mux.NewRouter()
	router.HandleFunc("/query", HandleQuery).Methods("GET")
	err := http.ListenAndServe(":"+server.GetRestPort(), router)
	if err != nil {
		log.Fatalln("Unable to Start REST Server")
	}
}

func PrintStatisticsUtility(wg *sync.WaitGroup, server *model.Server) {
	defer wg.Done()
	statsTimer := time.NewTicker(time.Second * 15).C
	for {
		select {
		case <-statsTimer:
			{
				log.Infoln("")
				fmt.Println("")
				fmt.Println("\t\t Server Statistics")
				influxServerIp := server.GetInfluxServer()
				totalWriteTime := server.GetDBWriteTime()
				totalPointsSaved := server.GetPointsSavedInDBCounter()
				avgWriteTime := totalWriteTime.Seconds() / float64(totalPointsSaved)
				fmt.Println("\t\t Average DB Write Time :     \t", avgWriteTime, " seconds")
				fmt.Println("\t\t Connected To Influx Server: \t", influxServerIp)
				totalPacketReceived := server.GetPacketReceivedCounter()
				activeAgents := []string{}
				server.RetrieveAllActiveClients(&activeAgents)
				for _, agent := range activeAgents {
					connectedAt := server.GetConectionTime(agent)
					fmt.Println("\t\t Connected Agent IP:         \t", agent)
					elapsed := time.Since(connectedAt)
					fmt.Println("\t\t Agent Active Since:         \t", elapsed.String())
				}
				fmt.Println("\t\t Total Packets Received:     \t", totalPacketReceived)
				fmt.Println("")
			}
		}
	}
}

func main() {

	serverUdpPort := flag.String("udpport", "7071", "udp port on the server")
	serverTcpPort := flag.String("tcpport", "8081", "tcp port of the server")
	restPort := flag.String("restPort", "9091", "tcp port of the server")
	influxServer := flag.String("influxServer", "127.0.0.1", "ip of influx server")
	influxPort := flag.String("influxPort", "8086", "port of influx server")

	flag.Parse()

	log.WithFields(logrus.Fields{
		"serverUdpPort": *serverUdpPort,
		"serverTcpPort": *serverTcpPort,
		"restPort":      *restPort,
		"influxServer":  *influxServer,
		"influxPort":    *influxPort,
	}).Info("Inputs from command line")

	server = model.NewServer(*serverTcpPort, *serverUdpPort, *restPort, *influxServer, *influxPort)
	var wg sync.WaitGroup
	wg.Add(4)

	go tcpListener(&wg, server)
	go udpListener(&wg, server)
	go predictionQueryListener(&wg, server)
	go PrintStatisticsUtility(&wg, server)

	wg.Wait()

}
