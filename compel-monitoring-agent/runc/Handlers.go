package runc

import (

"fmt"
"os/exec"
"os"
"strings"
"time"
utils "github.com/adhuri/Compel-Monitoring/utils"
)


// Function to get running containers ; Returns empty list if no container running

func GetRunningContainers() []string {
	
	//Track time using utils

	defer utils.TimeTrack(time.Now(), "Handlers.go-GetRunningContainers")

	//Defining byte buffer to store the output
	var (
			cmdOut []byte
			err    error
	)

	// Getting all containers having status = running

	status:="running"
	// Command to process the list of containers - returns each container name with \n seperated
	command := "runc list|grep "+ status +"| cut -d\" \" -f1"
	
	//Requires /bin/sh due to sudo permissions
	cmd := exec.Command("/bin/sh","-c",command) 
	

	if cmdOut, err = cmd.Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error in GetRunningContainers()- run list ", err)
	}
	containerList := strings.Split(string(cmdOut),"\n")

	//Empty list
	if len(containerList) == 1 {
		fmt.Println( " No container running ")
		return make([]string , 0 ) }
	// Since it contains "\n" 

	fmt.Println(" Containers running " , containerList, len(containerList)-1)
	return containerList[0:len(containerList)]

	//return make([]string, 4)
}

func GetContainerStats(containerId string) string {
	return containerId
}
