package protocol

import (
	"encoding/json"
	"fmt"
	"time"
)

// The stats message to be sent to
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
