package model

import (
	"errors"
	"sync"
	"time"
)

type statType struct {
	statName string
}

var CPU_STATS = statType{"CPU"}

type Client struct {
	sync.RWMutex
	containerStats         map[string]int64
	containerStatus        map[string]uint64
	oldTotalCPU            int64 // Since Difference of CPU is stored for system
	newTotalCPU            int64 // Since Difference of CPU is stored for system
	totalMemory            uint64
	serverAlive            bool
	serverIp               string
	serverUdpPort          string
	serverTcpPort          string
	connectionTime         time.Duration
	totalStatsMessagesSent int64
	totalAmountOfDataSent  int64
}

func NewClient(serverIp, tcpPort, udpPort string) *Client {
	return &Client{
		containerStats:         make(map[string]int64),
		containerStatus:        make(map[string]uint64),
		oldTotalCPU:            -1, // Set as -1 if first time CPU calculate or Stale CPU
		serverAlive:            false,
		serverIp:               serverIp,
		serverUdpPort:          udpPort,
		serverTcpPort:          tcpPort,
		totalStatsMessagesSent: 0,
		totalAmountOfDataSent:  0,
	}
}

// This method accepts container ID and statType
// It returns the value of the stats for the contianer
// It aquires a reader lock before reading the map
func (client *Client) GetStats(containerId string, stat statType) (int64, error) {
	key := containerId + stat.statName
	client.RLock()
	defer client.RUnlock()
	value, present := client.containerStats[key]
	if present {
		return value, nil
	} else {
		return -1, errors.New("Value not present")
	}
}

// This method accepts stat type, container id and value as input
// It acquires a writer lock before updating the map
func (client *Client) SetStats(stat statType, containerId string, value int64) {
	key := containerId + stat.statName
	client.Lock()
	defer client.Unlock()
	client.containerStats[key] = value
}

// return total CPU cycles usued by all containers
func (client *Client) GetTotalCPU() (int64, int64, error) {
	client.RLock()
	defer client.RUnlock()
	// Invalid CPU  - First Run or killed Container Stale Value
	if client.oldTotalCPU == -1 {
		return client.oldTotalCPU, client.newTotalCPU, errors.New("oldTotalCPU is -1 (Stale or First Run)")
	}
	return client.oldTotalCPU, client.newTotalCPU, nil
}

// sets total CPU cycles used by all the containers
func (client *Client) SetTotalCPU(value int64) {
	client.Lock()
	defer client.Unlock()

	client.oldTotalCPU = client.newTotalCPU // Since it needs swapping of new to Old
	client.newTotalCPU = value
}

// return total Memory limit
func (client *Client) GetTotalMemory() uint64 {
	client.RLock()
	defer client.RUnlock()
	return client.totalMemory
}

// sets total Memory Limit
func (client *Client) SetTotalMemory(value uint64) {
	client.Lock()
	defer client.Unlock()
	client.totalMemory = value
}

// Return True if the container stats were recorded in previous cycle
// We match the current counter with the counter of the container
func (client *Client) IsContainerAlive(containerId string, currentCounter uint64) bool {
	client.RLock()
	defer client.RUnlock()
	return (currentCounter-client.containerStatus[containerId] <= 1)
}

// We update the counter value everytime we send the stats for the given container to the server
func (client *Client) UpdateContainerCounter(containerId string, currentCounter uint64) {
	client.Lock()
	defer client.Unlock()
	client.containerStatus[containerId] = currentCounter
}

func (client *Client) UpdateServerStatus(isAlive bool) {
	client.Lock()
	defer client.Unlock()
	client.serverAlive = isAlive
}

func (client *Client) GetServerStatus() bool {
	client.RLock()
	defer client.RUnlock()
	return client.serverAlive
}

func (client *Client) GetServerTcpPort() string {
	client.RLock()
	defer client.RUnlock()
	return client.serverTcpPort

}

func (client *Client) GetServerUdpPort() string {
	client.RLock()
	defer client.RUnlock()
	return client.serverUdpPort
}

func (client *Client) GetServerIp() string {
	client.RLock()
	defer client.RUnlock()
	return client.serverIp

}

func (client *Client) GetConnectionTime() time.Duration {
	client.RLock()
	defer client.RUnlock()
	return client.connectionTime
}

func (client *Client) SetConnectionTime(timeTaken time.Duration) {
	client.Lock()
	defer client.Unlock()
	client.connectionTime = timeTaken
}

func (client *Client) GetTotalPacketsSent() int64 {
	client.RLock()
	defer client.RUnlock()
	return client.totalStatsMessagesSent
}

func (client *Client) IncrementTotalPacketsSent() {
	client.Lock()
	defer client.Unlock()
	client.totalStatsMessagesSent += 1
}

func (client *Client) GetTotalAmountDataSent() int64 {
	client.RLock()
	defer client.RUnlock()
	return client.totalAmountOfDataSent
}

func (client *Client) UpdateTotalAmountDataSent(dataSize int64) {
	client.Lock()
	defer client.Unlock()
	client.totalAmountOfDataSent += dataSize
}
