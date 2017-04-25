package model

import (
	"net"
	"sync"
	"time"
)

type Server struct {
	sync.Mutex
	connectedClients          map[string]int64
	activeContainers          map[string][]string
	udpPort                   string
	tcpPort                   string
	restPort                  string
	influxServer              string
	influxPort                string
	totalPacketsReceived      int64
	statsMap                  map[string]time.Time
	pointsSavedInDB           int64
	totalTimeTakenToStoreInDB time.Duration
}

func NewServer(tcpPort, udpPort, restPort, influxServer string, influxPort string) *Server {
	return &Server{
		connectedClients:          make(map[string]int64),
		activeContainers:          make(map[string][]string),
		udpPort:                   udpPort,
		tcpPort:                   tcpPort,
		restPort:                  restPort,
		influxServer:              influxServer,
		influxPort:                influxPort,
		totalPacketsReceived:      0,
		statsMap:                  make(map[string]time.Time),
		totalTimeTakenToStoreInDB: time.Nanosecond,
		pointsSavedInDB:           0,
	}
}

func (server *Server) IsAgentConnected(agentIp net.IP) bool {
	server.Lock()
	defer server.Unlock()
	currentTime := time.Now().Unix()
	value, present := server.connectedClients[agentIp.String()]
	if present {
		return (currentTime - value) < 40
	}
	return false
}

func (server *Server) RetrieveAllActiveClients(activeAgents *[]string) {
	server.Lock()
	defer server.Unlock()

	currentTime := time.Now().Unix()
	for key := range server.connectedClients {
		if (currentTime - server.connectedClients[key]) < 40 {
			*activeAgents = append(*activeAgents, key)
		} else {
			delete(server.connectedClients, key)
		}
	}
}

func (server *Server) RetrieveAllActiveContainers(agentIp string) []string {
	server.Lock()
	defer server.Unlock()

	containers := server.activeContainers[agentIp]
	actvContainers := make([]string, 0)
	for _, containerId := range containers {
		actvContainers = append(actvContainers, containerId)
	}
	return actvContainers
}

func (server *Server) SetActiveContainersForAgent(agentIp string, containerList []string) {
	server.Lock()
	defer server.Unlock()

	server.activeContainers[agentIp] = containerList
}

func (server *Server) addClient(ip string) {
	server.Lock()
	defer server.Unlock()

	server.connectedClients[ip] = time.Now().Unix()
}

func (server *Server) UpdateState(agentIp net.IP) {
	ip := agentIp.String()
	server.addClient(ip)
}

func (server *Server) GetUdpPort() string {
	server.Lock()
	defer server.Unlock()

	return server.udpPort
}

func (server *Server) GetTcpPort() string {
	server.Lock()
	defer server.Unlock()

	return server.tcpPort
}

func (server *Server) GetInfluxServer() string {
	server.Lock()
	defer server.Unlock()

	return server.influxServer
}

func (server *Server) GetInfluxPort() string {
	server.Lock()
	defer server.Unlock()

	return server.influxPort
}

func (server *Server) GetRestPort() string {
	server.Lock()
	defer server.Unlock()

	return server.restPort
}

func (server *Server) GetPacketReceivedCounter() int64 {
	server.Lock()
	defer server.Unlock()

	return server.totalPacketsReceived
}

func (server *Server) IncrementPacketReceivedCounter() {
	server.Lock()
	defer server.Unlock()

	server.totalPacketsReceived += 1
}

func (server *Server) GetConectionTime(agentIp string) time.Time {
	server.Lock()
	defer server.Unlock()

	return server.statsMap[agentIp]
}

func (server *Server) UpdateStatsMap(agentIp string, connectionTime time.Time) {
	server.Lock()
	defer server.Unlock()

	server.statsMap[agentIp] = connectionTime
}

func (server *Server) GetPointsSavedInDBCounter() int64 {
	server.Lock()
	defer server.Unlock()

	return server.pointsSavedInDB
}

func (server *Server) IncrementPointsSavedInDBCounterCounter() {
	server.Lock()
	defer server.Unlock()

	server.pointsSavedInDB += 1
}

func (server *Server) GetDBWriteTime() time.Duration {
	server.Lock()
	defer server.Unlock()

	return server.totalTimeTakenToStoreInDB
}

func (server *Server) UpdateDBWriteTime(writeTime time.Duration) {
	server.Lock()
	defer server.Unlock()

	server.totalTimeTakenToStoreInDB += writeTime
}
