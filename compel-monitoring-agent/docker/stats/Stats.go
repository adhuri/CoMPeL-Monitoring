package stats

//
// import (
// 	"fmt"
// 	"time"
//
// 	"github.com/adhuri/Compel-Monitoring/utils"
// )
//
// // func main() {
// // 	//	fmt.Print("Memory Used for container 1 : ", getContainerMemory("container1"), "\n")
// // 	//fmt.Print("Total System Memory : ", getSystemMemory(), "\n")
// // 	//fmt.Println(getDockerStatsFromCommandLine("9189afaebbea", "CPUPerc"))
// // 	//fmt.Println(getDockerStatsFromCommandLine("9189afaebbea", "MemPerc"))
// //
// // 	fmt.Println(CalculateMemoryPercentage("9189afaebbea"))
// // 	fmt.Println(CalculateCPUUsedPercentage("9189afaebbea"))
// // }
//
// //CalculateMemoryPercentage %Memory Used by the container
// func CalculateMemoryPercentage(containerID string) (memoryPercent float64) {
// 	defer utils.TimeTrack(time.Now(), "Stats.go- Docker CalculateMemoryPercentage")
// 	memoryPercent = 0.0
//
// 	if stat, exists := ds.Stats[containerID]; exists {
// 		//do something here
// 		memoryPercent = stat.MemoryPercent
// 		return
// 	}
// 	// Else if no key exists
// 	fmt.Println("CalculateMemoryPercentage() - Could not find key in datastructure DockerContainerStats for containerID", containerID)
// 	return
//
// }
//
// //CalculateCPUUsedPercentage %CPU Used by the container
// func (ds DockerCalculateCPUUsedPercentage(containerID string) (cpuPercent float64) {
//
// 	defer utils.TimeTrack(time.Now(), "stats.go- Docker CalculateCPUUsedPercentage")
// 	cpuPercent = 0.0
// 	if stat, exists := ds.Stats[containerID]; exists {
// 		//do something here
// 		cpuPercent = stat.CpuPercent
// 		return
// 	}
// 	// Else if no key exists
// 	fmt.Println("CalculateCPUUsedPercentage() - Could not find key in datastructure DockerContainerStats for containerID", containerID)
// 	return
//
// }
