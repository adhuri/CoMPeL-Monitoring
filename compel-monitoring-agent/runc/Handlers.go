package runc

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/model"
	stats "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc/stats"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
	utils "github.com/adhuri/Compel-Monitoring/utils"
)

// GetRunningContainers ... Function to get running containers ; Returns empty list if no container running
func GetRunningContainers() []string {

	//Track time using utils

	defer utils.TimeTrack(time.Now(), "Handlers.go-GetRunningContainers")

	//Defining byte buffer to store the output
	var (
		cmdOut []byte
		err    error
	)

	// Getting all containers having status = running

	status := "running"
	// Command to process the list of containers - returns each container name with \n seperated
	command := "runc list|grep " + status + "| cut -d\" \" -f1"

	//Requires /bin/sh due to sudo permissions
	cmd := exec.Command("/bin/sh", "-c", command)

	if cmdOut, err = cmd.Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error in GetRunningContainers()- run list ", err)
	}
	containerList := strings.Split(string(cmdOut), "\n")

	//Empty list
	if len(containerList) == 1 {
		fmt.Println(" No container running ")
		return make([]string, 0)
	}
	// Since it contains "\n"

	fmt.Println(" Containers running ", containerList, len(containerList)-1)
	return containerList[0 : len(containerList)-1]

	//return make([]string, 4)
}

func GetContainerStats(client *model.Client, containerID string) string {

	//Timing this function
	defer utils.TimeTrack(time.Now(), "Handlers.go-GetContainerStats")

	//Calculating Memory Used
	memoryPercentage := strconv.FormatFloat(stats.CalculateMemoryPercentage(client, containerID), 'f', -1, 32)

	//Calculating CPU Used
	cpuPercentage := strconv.FormatFloat(stats.CalculateCPUUsedPercentage(client, containerID), 'f', -1, 32)

	message := monitorProtocol.EncodeStatsJSON(containerID, cpuPercentage, memoryPercentage)
	return string(message)
}
