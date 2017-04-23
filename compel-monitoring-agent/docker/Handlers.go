package docker

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
	"github.com/adhuri/Compel-Monitoring/utils"
)

func (ds *DockerContainerStats) GetDockerStats(log *logrus.Logger) {

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
		log.Errorln(os.Stderr, "There was an error in dockerstats.go-GetRunningDockerContainers()- ", err)
	}
	containerDataList := strings.Split(string(cmdOut), "\n")

	if len(containerDataList) == 0 {
		log.Warnln("Handlers.go - GetDockerStats() No containers running ")
	}
	//log.Debugln("Refreshing DockerStats object for issue #5")
	//ds = NewDockerContainerStats()

	for _, el := range containerDataList {
		if el != "" {
			// All elements should be parseable
			containerId, cpuPercent, memoryPercent := parseContainerDetails(el, log)
			//ds.Stats[containerId] = NewStatType(cpuPercent, memoryPercent)
			// To avoid lock issues
			ds.SetContainerStat(containerId, NewStatType(cpuPercent, memoryPercent))
		}
	}

	//fmt.Println(" Containers running ", containerDataList, len(containerDataList)-1)
	//return containerDataList[0 : len(containerDataList)-1]

	//return make([]string, 4)

}

func parseContainerDetails(line string, log *logrus.Logger) (containerID string, cpuPercent float64, memoryPercent float64) {
	// Since the format is colon seperated
	containerDetails := strings.Split(line, ":")
	log.Infoln(containerDetails)

	if len(containerDetails) > 3 {
		log.Errorln("parseContainerDetails() - Seems you added NetBlock or Disk IO but forgot to parse it")
	} else if len(containerDetails) < 3 {
		log.Errorln("parseContainerDetails() - Did you delete CPU Percentage or Disk IO but forgot to unparse it")
	}

	containerID = containerDetails[0]
	cpuPercent, _ = strconv.ParseFloat(strings.Trim(containerDetails[1], "%"), 64)
	memoryPercent, _ = strconv.ParseFloat(strings.Trim(containerDetails[2], "%"), 64)

	return

}

func GetRunningContainers(ds *DockerContainerStats, log *logrus.Logger) []string {

	//Track time using utils
	defer utils.TimeTrack(time.Now(), "dockerstats.go-GetRunningDockerContainers")

	containerDataList := make([]string, 0, len(ds.Stats))
	for k := range ds.GetAllContainerStat() {
		containerDataList = append(containerDataList, k)
	}

	//Empty list
	if len(containerDataList) == 0 {
		log.Warnln(" No container running ")
		return make([]string, 0)
	}
	// Since it contains "\n"

	log.Infoln("Containers running ", len(containerDataList), containerDataList)
	//return containerDataList[0 : len(containerDataList)-1]
	return containerDataList[:]
	//return make([]string, 4)
}

func GetContainerStats(ds *DockerContainerStats, containerID string, log *logrus.Logger) monitorProtocol.ContainerStats {

	//Timing this function
	defer utils.TimeTrack(time.Now(), "Handlers.go-GetContainerStats")

	//Calculating Memory Used
	memoryPercentage := CalculateMemoryPercentage(ds, containerID, log)
	log.Debugln("memoryPercentage for container ", containerID, " - ", memoryPercentage)

	//Calculating CPU Used
	cpuPercentage := CalculateCPUUsedPercentage(ds, containerID, log)
	log.Debugln("cpuPercentage for container ", containerID, " - ", cpuPercentage)

	message := monitorProtocol.GetContainerStats(containerID, cpuPercentage, memoryPercentage)
	return message
}
