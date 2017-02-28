package model

import "sync"

type Server struct {
	sync.Mutex
	clients map[int64]string
}

func New() *Server {
	return &Server{
		clients: make(map[int64]string),
	}
}

func retrieveAllClients() []string {

}

func addClient() {

}
