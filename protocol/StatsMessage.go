package protocol

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adhuri/Compel-Monitoring/utils"
)

type StatsMessage struct {
	MessageId int64
	AgentIP   [4]byte
	Data      []string
}

// The stats to be sent to
type StatsJSON struct {
	ContainerID string `json:"containerID"`
	Timestamp   int64  `json:"timestamp"`
	Data        `json:"data"`
}

type Data struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

func EncodeStatsJSON(cID, cpu, memory string) []byte {

	message := &StatsJSON{
		ContainerID: cID,
		Timestamp:   time.Now().Unix(),
		Data: Data{
			CPU:    cpu,
			Memory: memory,
		},
	}
	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error : Marshalling JSON in StatsMessage", err)
		return []byte(`{}`)
	}
	return b

}

func NewStatsMessage(dataToSend []string) *StatsMessage {

	// Get External IP of host
	var hostIP [4]byte
	err := utils.GetIPAddressOfHost(hostIP[0:])
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
