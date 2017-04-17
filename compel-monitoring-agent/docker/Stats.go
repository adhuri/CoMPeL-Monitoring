package docker

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Monitoring/utils"
)

//CalculateMemoryPercentage %Memory Used by the container
func CalculateMemoryPercentage(ds *DockerContainerStats, containerID string, log *logrus.Logger) (memoryPercent float64) {
	defer utils.TimeTrack(time.Now(), "Stats.go- Docker CalculateMemoryPercentage")
	memoryPercent = 0.0

	if stat, exists := ds.Stats[containerID]; exists {
		//do something here
		memoryPercent = stat.MemoryPercent
		return
	}
	// Else if no key exists
	log.Errorln("CalculateMemoryPercentage() - Could not find key in datastructure DockerContainerStats for containerID", containerID)
	return

}

//CalculateCPUUsedPercentage %CPU Used by the container
func CalculateCPUUsedPercentage(ds *DockerContainerStats, containerID string, log *logrus.Logger) (cpuPercent float64) {

	defer utils.TimeTrack(time.Now(), "stats.go- Docker CalculateCPUUsedPercentage")
	cpuPercent = 0.0
	if stat, exists := ds.Stats[containerID]; exists {
		//do something here
		cpuPercent = stat.CpuPercent
		return
	}
	// Else if no key exists
	log.Errorln("CalculateCPUUsedPercentage() - Could not find key in datastructure DockerContainerStats for containerID", containerID)
	return

}
