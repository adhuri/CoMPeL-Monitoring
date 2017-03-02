package model

import "sync"

type Server struct {
	sync.RWMutex
	connectedClients map[int64]string
}

func NewServer() *Server {
	return &Server{
		connectedClients: make(map[int64]string),
	}
}

func retrieveAllClients() {

}

func addClient() {

}

func (server *Server) handleStatMessage() {

}
