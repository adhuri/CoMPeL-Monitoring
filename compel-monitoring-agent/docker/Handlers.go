package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
	"github.com/adhuri/Compel-Monitoring/utils"
)

func main() {
	DS := NewDockerContainerStats()

	DS.GetDockerStats()

	fmt.Println(DS.GetRunningDockerContainers())

}

func (ds *DockerContainerStats) GetDockerStats() {

	//Defining byte buffer to store the output
	var (
		cmdOut []byte
		err    error
	)

	// Getting all Running Docker containers

	// Command to process the list of containers - returns each container name with \n seperated
	//command := "docker stats --format \"{{.Container}}: {{.CPUPerc}} : {{.MemPerc}} : {{.NetIO}} : {{.BlockIO}}\" --no-stream"
	command := "docker stats --format \"{{.Container}}:{{.CPUPerc}}:{{.MemPerc}}\" --no-stream"

	//Requires /bin/sh due to sudo permissions
	cmd := exec.Command("/bin/sh", "-c", command)

	if cmdOut, err = cmd.Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error in dockerstats.go-GetRunningDockerContainers()- ", err)
	}
	containerDataList := strings.Split(string(cmdOut), "\n")

	if len(containerDataList) == 0 {
		fmt.Println("Handlers.go - GetDockerStats() No containers running ")
	}

	for _, el := range containerDataList {
		if el != "" {
			// All elements should be parseable
			containerId, cpuPercent, memoryPecent := parseContainerDetails(el)
			ds.Stats[containerId] = StatType{CpuPercent: cpuPercent, MemoryPercent: memoryPecent}
		}
	}

	//fmt.Println(" Containers running ", containerDataList, len(containerDataList)-1)
	//return containerDataList[0 : len(containerDataList)-1]

	//return make([]string, 4)

}

func parseContainerDetails(line string) (containerID string, cpuPercent float64, memoryPercent float64) {
	// Since the format is colon seperated
	containerDetails := strings.Split(line, ":")
	fmt.Println(containerDetails)

	if len(containerDetails) > 3 {
		fmt.Println("parseContainerDetails() - Seems you added NetBlock or Disk IO but forgot to parse it")
	} else if len(containerDetails) < 3 {
		fmt.Println("parseContainerDetails() - Did you delete CPU Percentage or Disk IO but forgot to unparse it")
	}

	containerID = containerDetails[0]
	cpuPercent, _ = strconv.ParseFloat(strings.Trim(containerDetails[1], "%"), 64)
	memoryPercent, _ = strconv.ParseFloat(strings.Trim(containerDetails[2], "%"), 64)

	return

}

func (ds *DockerContainerStats) GetRunningDockerContainers() []string {

	//Track time using utils
	defer utils.TimeTrack(time.Now(), "dockerstats.go-GetRunningDockerContainers")

	containerDataList := make([]string, 0, len(ds.Stats))
	for k := range ds.Stats {
		containerDataList = append(containerDataList, k)
	}

	//Empty list
	if len(containerDataList) == 0 {
		fmt.Println(" No container running ")
		return make([]string, 0)
	}
	// Since it contains "\n"

	fmt.Println(" Containers running ", containerDataList, len(containerDataList))
	return containerDataList[0 : len(containerDataList)-1]

	//return make([]string, 4)
}

func (ds *DockerContainerStats) GetContainerStats(containerID string) monitorProtocol.ContainerStats {

	//Timing this function
	defer utils.TimeTrack(time.Now(), "Handlers.go-GetContainerStats")

	//Calculating Memory Used
	memoryPercentage := CalculateMemoryPercentage(ds, containerID)

	//Calculating CPU Used
	cpuPercentage := CalculateCPUUsedPercentage(ds, containerID)

	message := monitorProtocol.GetContainerStats(containerID, cpuPercentage, memoryPercentage)
	return message
}
