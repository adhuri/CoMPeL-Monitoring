package model

import (
	"net"
	"sync"
	"time"
)

type Server struct {
	sync.Mutex
	connectedClients map[string]int64
}

func NewServer() *Server {
	return &Server{
		connectedClients: make(map[string]int64),
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
