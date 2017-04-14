package stats

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/adhuri/Compel-Monitoring/utils"
)

func main() {
	//	fmt.Print("Memory Used for container 1 : ", getContainerMemory("container1"), "\n")
	//fmt.Print("Total System Memory : ", getSystemMemory(), "\n")
	//fmt.Println(getDockerStatsFromCommandLine("9189afaebbea", "CPUPerc"))
	//fmt.Println(getDockerStatsFromCommandLine("9189afaebbea", "MemPerc"))

	fmt.Println(CalculateMemoryPercentage("9189afaebbea"))
	fmt.Println(CalculateCPUUsedPercentage("9189afaebbea"))
}

//CalculateMemoryPercentage %Memory Used by the container
func CalculateMemoryPercentage(containerID string) (memoryPercent float64) {
	defer utils.TimeTrack(time.Now(), "Stats.go- Docker CalculateMemoryPercentage")
	memoryPercent = 0.0

	memory, err := getDockerStatsFromCommandLine(containerID, "MemPerc")
	if err != nil {
		fmt.Println("Could not get Memory Percentage statistics for ", containerID)
		return memoryPercent
	}
	memoryPercent, _ = strconv.ParseFloat(strings.Trim(memory, "%"), 64)
	return

}

//CalculateCPUUsedPercentage %CPU Used by the container
func CalculateCPUUsedPercentage(containerID string) (cpuPercent float64) {

	defer utils.TimeTrack(time.Now(), "stats.go- Docker CalculateCPUUsedPercentage")
	cpuPercent = 0.0

	cpu, err := getDockerStatsFromCommandLine(containerID, "CPUPerc")
	if err != nil {
		fmt.Println("Could not get Memory Percentage statistics for ", containerID)
		return
	}
	cpuPercent, _ = strconv.ParseFloat(strings.Trim(cpu, "%"), 64)

	return

}

// Gets direct stats of container from commandLine
// Available stats Type {{.Container}}: {{.CPUPerc}} : {{.MemPerc}} : {{.NetIO}} : {{.BlockIO
func getDockerStatsFromCommandLine(containerID string, statsType string) (string, error) {
	//Defining byte buffer to store the output
	var (
		cmdOut []byte
		err    error
	)
	//
	//command := "docker stats --format \"{{.Container}}: {{.CPUPerc}} : {{.MemPerc}} : {{.NetIO}} : {{.BlockIO}}\" --no-stream"
	command := "docker stats " + containerID + " --format \"{{." + statsType + "}}\" --no-stream"
	fmt.Println(command)
	//Requires /bin/sh due to sudo permissions
	cmd := exec.Command("/bin/sh", "-c", command)

	if cmdOut, err = cmd.Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error in stats.go-getDockerStatsFromCommandLine()- ", err)
		return "", err
	}
	containerStats := strings.Split(string(cmdOut), "\n")

	return containerStats[0], nil

}
