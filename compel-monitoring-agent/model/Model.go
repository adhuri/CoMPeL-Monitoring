package model

import (
	"errors"
	"sync"
)

type statType struct {
	statName string
}

var CPU_STATS = statType{"CPU"}
var MEM_STATS = statType{"MEMORY"}
var BLKIO_STATS = statType{"BLKIO"}

type Client struct {
	sync.RWMutex
	containerStats  map[string]string
	containerStatus map[string]uint64
	totalMemory     string
	totalCPU        string
	serverAlive     bool
}

func NewClient() *Client {
	return &Client{
		containerStats:  make(map[string]string),
		containerStatus: make(map[string]uint64),
		totalMemory:     "",
		totalCPU:        "",
		serverAlive:     false,
	}
}

// This method accepts container ID and statType
// It returns the value of the stats for the contianer
// It aquires a reader lock before reading the map
func (client *Client) GetStats(containerId string, stat statType) (string, error) {
	key := containerId + stat.statName
	client.RLock()
	defer client.RUnlock()
	value, present := client.containerStats[key]
	if present {
		return value, nil
	} else {
		return "", errors.New("Value not present")
	}
}

// This method accepts stat type, container id and value as input
// It acquires a writer lock before updating the map
func (client *Client) SetStats(stat statType, containerId string, value string) {
	key := containerId + stat.statName
	client.Lock()
	defer client.Unlock()
	client.containerStats[key] = value
}

// returns total memory used by all the containers
func (client *Client) GetTotalMemStats() string {
	client.RLock()
	defer client.RUnlock()
	return client.totalMemory
}

// return total CPU cycles usued by all containers
func (client *Client) GetTotalCPUStats() string {
	client.RLock()
	defer client.RUnlock()
	return client.totalCPU
}

// sets total memory used by all the containers
func (client *Client) SetTotalMemStats(value string) {
	client.Lock()
	defer client.Unlock()
	client.totalMemory = value
}

// sets total CPU cycles used by all the containers
func (client *Client) SetTotalCPUStats(value string) {
	client.Lock()
	defer client.Unlock()
	client.totalCPU = value
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
