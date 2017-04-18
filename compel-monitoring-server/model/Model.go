package model

import (
	"net"
	"sync"
	"time"
)

type Server struct {
	sync.Mutex
	connectedClients map[string]int64
	activeContainers map[string][]string
	udpPort          string
	tcpPort          string
	restPort         string
	influxServer     string
}

func NewServer(tcpPort, udpPort, restPort, influxServer string) *Server {
	return &Server{
		connectedClients: make(map[string]int64),
		activeContainers: make(map[string][]string),
		udpPort:          udpPort,
		tcpPort:          tcpPort,
		restPort:         restPort,
		influxServer:     influxServer,
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

func (server *Server) GetRestPort() string {
	server.Lock()
	defer server.Unlock()

	return server.restPort
}
