package model

import "sync"

type statType struct {
	string statName
}

var CPU_STATS = statType{"CPU"}
var MEM_STATS = statType{"MEMORY"}
var BLKIO_STATS = statType{"BLKIO"}

type Client struct {
	sync.RWMutex
	contianerStats map[string]string
	totalMemory    string
	totalCPU       string
}

// This method accepts container ID and statType
// It returns the value of the stats for the contianer
// It aquires a reader lock before reading the map
func (client *Client) GetStats(containerId string, stat statType) (string, error) {
	key := containerId + stat.statName
	client.RLock()
	defer client.RUnlock()
	value, present = contianerStats[key]
	if present {
		return value, nil
	} else {
		return nil, error.New("Value not present")
	}
}

// This method accepts stat type, container id and value as input
// It acquires a writer lock before updating the map
func (client *Client) SetStats(stat statType, containerId string, value string) {
	key := containerId + stat.statName
	client.Lock()
	defer client.Unlock()
	contianerStats[key] = value
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
