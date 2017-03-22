package protocol

import (
	"net"
	"time"

	"github.com/adhuri/Compel-Monitoring/utils"
)

type StatsMessage struct {
	MessageId int64
	AgentIP   net.IP
	Data      []ContainerStats
}

// The stats to be sent to
type ContainerStats struct {
	ContainerID string
	Timestamp   int64
	Data
}

type Data struct {
	CPU    string
	Memory string
}

func GetContainerStats(cID, cpu, memory string) ContainerStats {

	message := ContainerStats{
		ContainerID: cID,
		Timestamp:   time.Now().Unix(),
		Data: Data{
			CPU:    cpu,
			Memory: memory,
		},
	}
	return message

}

func NewStatsMessage(dataToSend []ContainerStats) *StatsMessage {

	// Get External IP of host
	//var hostIP [4]byte
	hostIP, err := utils.GetIPAddressOfHost()
	// If external IP of the host is not found then panic
	if err != nil {
		panic("Error Fetching Valid IP Address")
	}

	return &StatsMessage{
		MessageId: time.Now().Unix(),
		AgentIP:   hostIP,
		Data:      dataToSend,
	}
}
