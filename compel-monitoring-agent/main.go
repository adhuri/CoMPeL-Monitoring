package main

import (
	"fmt"
	"net"
	"time"

	model "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/model"
	runc "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc"
	stats "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc/stats"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func worker(client *model.Client, containerId string, containerStats chan string, currentCounter uint64) {
	stats := runc.GetContainerStats(client, containerId)
	containerStats <- stats
}

func sendStats(client *model.Client, counter uint64) {

	//Set SystemCPU usage
	sysCPUusage, err := stats.GetSystemCPU()
	if err != nil {
		fmt.Println("Error : cannot GetSystemCPU")
	} else {
		client.SetTotalCPU(sysCPUusage)
	}
	//Set Memory Limit
	sysMemoryLimit, err := stats.GetSystemMemory()
	if err != nil {
		fmt.Println("Error : cannot GetSystemCPU")
	} else {
		client.SetTotalMemory(sysMemoryLimit)
	}

	var containers []string = runc.GetRunningContainers()
	numOfWorkers := len(containers)
	containerStats := make(chan string, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		client.UpdateContainerCounter(containers[i], counter)
		go worker(client, containers[i], containerStats, counter)
	}

	//var buffer bytes.Buffer
	var statsToSend = make([]string, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		//buffer.WriteString(<-containerStats)
		statsToSend[i] = <-containerStats
	}
	//stringToSend := buffer.String()

	monitorProtocol.SendContainerStatistics(statsToSend)
}

func checkIfServerIsAlive() bool {
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func main() {
	client := model.NewClient()
	var counter uint64 = 0
	monitorProtocol.ConnectToServer()
	client.UpdateServerStatus(true)
	statsTimer := time.NewTicker(time.Second * 2).C
	aliveTimer := time.NewTicker(time.Second * 10).C
	for {
		select {
		case <-statsTimer:
			{
				if client.GetServerStatus() {
					counter++
					sendStats(client, counter)
				} else {
					fmt.Println("Server Offline .... Trying to Reconnect")
					monitorProtocol.ConnectToServer()
					client.UpdateServerStatus(true)
				}
			}
		case <-aliveTimer:
			{
				isAlive := checkIfServerIsAlive()
				if !isAlive {
					// update the server status
					fmt.Println("Server Dead")
					client.UpdateServerStatus(false)
				} else {
					fmt.Println("Server is still Alive")
				}
			}
		}
	}

}
