package runc

import (
	"os/exec"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/model"
	stats "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc/stats"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
	utils "github.com/adhuri/Compel-Monitoring/utils"
)

// GetRunningContainers ... Function to get running containers ; Returns empty list if no container running
func GetRunningContainers(log *logrus.Logger) []string {

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
	command := "runc list| grep -v \"docker\" | grep " + status + "| cut -d\" \" -f1"

	//Requires /bin/sh due to sudo permissions
	cmd := exec.Command("/bin/sh", "-c", command)

	if cmdOut, err = cmd.Output(); err != nil {
		log.Errorln("There was an error in GetRunningContainers()- run list ", err)
	}
	containerList := strings.Split(string(cmdOut), "\n")

	//Empty list
	if len(containerList) == 1 {
		log.Warnln("No container running")
		return make([]string, 0)
	}
	// Since it contains "\n"

	log.Infoln("Total running containers are : ", len(containerList)-1)

	//fmt.Println(" Containers running ", containerList, len(containerList)-1)
	log.Debugln("Running Containers are: ", containerList)
	return containerList[0 : len(containerList)-1]

	//return make([]string, 4)
}

func GetContainerStats(client *model.Client, containerID string, log *logrus.Logger) monitorProtocol.ContainerStats {

	//Timing this function
	defer utils.TimeTrack(time.Now(), "Handlers.go-GetContainerStats")

	//Calculating Memory Used
	memoryPercentage := stats.CalculateMemoryPercentage(client, containerID, log)

	//Calculating CPU Used
	cpuPercentage := stats.CalculateCPUUsedPercentage(client, containerID, log)

	// Creating Stat Message
	message := monitorProtocol.GetContainerStats(containerID, cpuPercentage, memoryPercentage)
	return message
}
