package stats

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/model"
	"github.com/adhuri/Compel-Monitoring/utils"
	"github.com/opencontainers/runc/libcontainer/system"
)

func main() {
	//	fmt.Print("Memory Used for container 1 : ", getContainerMemory("container1"), "\n")
	//fmt.Print("Total System Memory : ", getSystemMemory(), "\n")

}

//Get Memory used for a container used using cgroups from /sys/fs/cgroup/memory/user.slice/<containerName>/memory.stat
func getContainerMemory(container string) (memoryused int64) {
	// Memory Used = total_cache + total_rss - https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt Section 5.5
	memoryused = 0
	if len(container) <= 0 {
		fmt.Printf("No container name defined")
		return
	}
	contents, err := ioutil.ReadFile("/sys/fs/cgroup/memory/user.slice/" + container + "/memory.stat")
	if err != nil {
		fmt.Print("Error : ioutil Read fail for getContainerMemory")
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 0 {
			//Parsing Total_rss

			if fields[0] == "total_rss" {

				// If any issues with ParseInt read Kernel documentation https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
				val, err := strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					fmt.Println("Error: total_rss", fields[1], err)
				}
				memoryused += val // Adding the total_rss to memory used
				//fmt.Printf("total_rss %d",val)
			}
			// Parsing Total_cache
			if fields[0] == "total_cache" {

				// If any issues with ParseInt read Kernel documentation https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
				val, err := strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					fmt.Println("Error: total_cache", fields[1], err)
				}
				memoryused += val // Adding the total_rss to memory used
				//fmt.Printf("total_cache %d",val)

			}

		}

	}

	//fmt.Print("memoryused ", memoryused)

	return

}

//GetSystemMemory Total Memory for the system using cgroups from /sys/fs/cgroup/memory/user.slice/<containerName>/memory.stat
func GetSystemMemory() (uint64, error) {
	// We get System Memory from /proc/meminfo MemTotal
	defer utils.TimeTrack(time.Now(), "Stats.go-GetSystemMemory")
	var totalmemory uint64
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Print("ERROR : Unable to read /proc/meminfo")
		return 0, err
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 0 {
			//Parsing Total_rss

			if fields[0] == "MemTotal:" {
				val, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					fmt.Println("Error: MemTotal:", fields[1], err)
				}
				totalmemory = val * 1024 // KB to Bytes
				return totalmemory, nil
			}
		}
	}
	//fmt.Print("totalmemory", totalmemory)
	return totalmemory, nil
}

//CalculateMemoryUsed %Memory Used by the container
func CalculateMemoryUsed(client *model.Client, containerID string) (memorypercent float64) {

	//cmemory := getContainerMemory(containerID)
	cmemory := getContainerMemory(containerID)
	tmemory := client.GetTotalMemory() // Might return 0 if not set due to some issue. Log printed
	if tmemory == 0 {
		fmt.Print("Error : Get System Memory returned 0 ")
		return
	}

	memorypercent = float64(cmemory) / float64(tmemory)
	fmt.Printf("cmemory %d , tmemory %d \n", cmemory, tmemory)
	fmt.Printf("memorypercent %f %%\n", memorypercent)
	return

}

//Get CPU used for a container used using cgroups from /sys/fs/cgroup/cpu,cpuacct/user.slice/<containername>/cpuacct.usage
//This we get by looking at /proc/<container pid>/cgroups

func getContainerCPU(container string) (cpuused int64) {

	cpuused = 0
	if len(container) <= 0 {
		fmt.Printf("No container name defined")
		return
	}
	contents, err := ioutil.ReadFile("/sys/fs/cgroup/cpu,cpuacct/user.slice/" + container + "/cpuacct.usage")
	if err != nil {
		fmt.Print("Error : ioutil Read fail for getContainerMemory")
		return
	}
	lines := strings.Split(string(contents), "\n")
	val, err := strconv.ParseInt(lines[0], 10, 64)
	if err != nil {
		fmt.Println("Error: cpuacct.usage not giving out integer usage", lines[0], err)
	}
	cpuused = val // cpu used from cpuacct.usag
	fmt.Printf("cpuused %d", cpuused)
	return

}

//GetSystemCPU Total CPU for the system using cgroups from /sys/fs/cgroup/memory/user.slice/<containerName>/memory.stat
func GetSystemCPU() (uint64, error) {
	//Using Docker code from https://github.com/docker/docker/blob/cd6a61f1b17830464250406244ed8ef113db8a3c/daemon/stats/collector_unix.go
	defer utils.TimeTrack(time.Now(), "Stats.go-GetSystemCPU")
	const nanoSecondsPerSecond = 1e9

	clockTicksPerSecond := uint64(system.GetClockTicks())
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		fmt.Print("ERROR : Unable to read /proc/meminfo")
		return 0, err
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 0 {
			switch parts[0] {
			case "cpu":
				if len(parts) < 8 {
					return 0, fmt.Errorf("invalid number of cpu fields")
				}
				var totalClockTicks uint64
				for _, i := range parts[1:8] {
					v, err := strconv.ParseUint(i, 10, 64)
					if err != nil {
						return 0, fmt.Errorf("Unable to convert value %s to int: %s", i, err)
					}
					totalClockTicks += v
				}
				return (totalClockTicks * nanoSecondsPerSecond) / clockTicksPerSecond, nil
			}
		}
	}
	return 0, fmt.Errorf("invalid stat format. Error trying to parse the '/proc/stat' file")
}
