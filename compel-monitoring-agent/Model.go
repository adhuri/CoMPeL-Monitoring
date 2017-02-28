package model

import "sync"

type statType struct {
	string statName
}

var CPU_STATS = statType{"CPU"}
var MEM_STATS = statType{"MEMORY"}
var BLKIO_STATS = statType{"BLKIO"}

type Client struct {
	sync.Mutex
	contianerStats map[string]string
}

func (client *Client) GetStats(containerId string, stat statType) (string, error) {
	key := containerId + stat.string
	client.Lock()
	defer client.Unlock()
	value, present = contianerStats[key]
	if present {
		return value, nil
	} else {
		return nil, error.New("Value not present")
	}
}

func (*Client) SetStats(stat statType, containerId string, value string) {
	key := containerId + stat.statName
	client.Lock()
	defer client.Unlock()
	contianerStats[key] = value
}
