package stats

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/adhuri/Compel-Monitoring/utils"
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
func getSystemMemory() (totalmemory int64) {
	// We get System Memory from /proc/meminfo MemTotal
	defer utils.TimeTrack(time.Now(), "Stats.go-GetSystemMemory")
	totalmemory = 0
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Print("ERROR : Unable to read /proc/meminfo")
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 0 {
			//Parsing Total_rss

			if fields[0] == "MemTotal:" {
				val, err := strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					fmt.Println("Error: MemTotal:", fields[1], err)
				}
				totalmemory = val * 1024 // KB to Bytes
				return
			}
		}
	}
	//fmt.Print("totalmemory", totalmemory)
	return
}

//CalculateMemoryUsed %Memory Used by the container
func CalculateMemoryUsed(containerID string) (memorypercent float64) {

	cmemory := getContainerMemory(containerID)
	tmemory := getSystemMemory() // Need to shift it to client getTotalMemory
	if tmemory == 0 {
		fmt.Print("Error : Get System Memory returned 0 ")
		return
	}

	memorypercent = float64(cmemory) / float64(tmemory)
	fmt.Printf("cmemory %d , tmemory %d \n", cmemory, tmemory)
	fmt.Printf("memorypercent %f %%\n", memorypercent)
	return

}

//Get Total Memory for the system using cgroups from /sys/fs/cgroup/memory/user.slice/<containerName>/memory.stat
//func getSystemCPU() (systemcpu int64) {
//
//}
