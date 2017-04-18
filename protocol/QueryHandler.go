package protocol

import model "github.com/adhuri/Compel-Monitoring/compel-monitoring-server/model"

type QueryResponse struct {
	Clients []Client `json:"clients,omitempty"`
}

type Client struct {
	ClientIp   string   `json:"agentId,omitempty"`
	Containers []string `json:"containers,omitempty"`
}

func GenerateQueryResponse(server *model.Server) *QueryResponse {
	activeAgents := make([]string, 0)
	clients := make([]Client, 0)

	server.RetrieveAllActiveClients(&activeAgents)
	for _, agent := range activeAgents {
		containerList := server.RetrieveAllActiveContainers(agent)
		clients = append(clients, Client{ClientIp: agent, Containers: containerList})
	}

	// fmt.Println(activeAgents)
	// fmt.Println()

	return &QueryResponse{Clients: clients}

}
