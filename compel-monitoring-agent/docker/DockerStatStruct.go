package docker

import "sync"

// Hashmap to store the docker status output only once since it takes 2 seconds to call docker stats ( for correct cpu percent calculation)
type DockerContainerStats struct {
	//sync.Mutex
	sync.RWMutex
	Stats map[string]*StatType
}

type StatType struct {
	CpuPercent    float64
	MemoryPercent float64
}

func NewDockerContainerStats() *DockerContainerStats {
	return &DockerContainerStats{
		Stats: make(map[string]*StatType),
	}
}

func NewStatType(cpuPercent float64, memoryPercent float64) *StatType {
	return &StatType{CpuPercent: cpuPercent, MemoryPercent: memoryPercent}
}

func (ds *DockerContainerStats) SetContainerStat(containerID string, st *StatType) {
	ds.Lock()
	defer ds.Unlock()
	ds.Stats[containerID] = st
}

func (ds *DockerContainerStats) GetContainerStat(containerID string) (st *StatType, exists bool) {
	ds.RLock()
	defer ds.RUnlock()
	stat, exists := ds.Stats[containerID]
	return stat, exists
}
