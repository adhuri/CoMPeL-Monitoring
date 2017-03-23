package model

import (
	"net"
	"sync"
	"time"
)

type Server struct {
	sync.Mutex
	connectedClients map[string]int64
	udpPort          string
	tcpPort          string
}

func NewServer(tcpPort, udpPort string) *Server {
	return &Server{
		connectedClients: make(map[string]int64),
		udpPort:          udpPort,
		tcpPort:          tcpPort,
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

func (server *Server) RetrieveAllActiveClients(activeAgents []string) {
	server.Lock()
	defer server.Unlock()

	currentTime := time.Now().Unix()
	for key := range server.connectedClients {
		if (currentTime - server.connectedClients[key]) < 40 {
			activeAgents = append(activeAgents, key)
		} else {
			delete(server.connectedClients, key)
		}
	}
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
