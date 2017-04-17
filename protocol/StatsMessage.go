package protocol

import (
	"fmt"
	"net"
	"time"
	"unsafe"

	"github.com/adhuri/Compel-Monitoring/utils"
	"github.com/mitchellh/hashstructure"
)

type StatsMessage struct {
	MessageId int64
	AgentIP   net.IP
	HashCode  uint64
	Data      []ContainerStats
}

// The stats to be sent to
type ContainerStats struct {
	ContainerID string
	Timestamp   int64
	MetricData  Data
}

type Data struct {
	CPU    float64
	Memory float64
}

func (stats *ContainerStats) Size() int {
	size := int(unsafe.Sizeof(*stats))
	size += len(stats.ContainerID)
	return size
}

func GetContainerStats(cID string, cpu float64, memory float64) ContainerStats {

	message := ContainerStats{
		ContainerID: cID,
		Timestamp:   time.Now().Unix(),
		MetricData: Data{
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

	hash, err := hashstructure.Hash(dataToSend, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Hash : ", hash)

	return &StatsMessage{
		MessageId: time.Now().Unix(),
		AgentIP:   hostIP,
		HashCode:  hash,
		Data:      dataToSend,
	}
}
