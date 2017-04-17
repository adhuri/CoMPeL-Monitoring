package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	docker "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/docker"
	model "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/model"
	runc "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc"
	stats "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc/stats"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

//Interface to choose Docker or RunC
type StatsInterface interface {
	worker(client *model.Client, containerId string, containerStats chan monitorProtocol.ContainerStats, currentCounter uint64)
	sendStats(client *model.Client, counter uint64)
}

type DockerStats struct {
	dockerContainerStats *docker.DockerContainerStats
}

type RuncStats struct {
}

func checkIfServerIsAlive(client *model.Client) bool {
	conn, err := net.Dial("tcp", client.GetServerIp()+":"+client.GetServerTcpPort())
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func main() {
	serverIp := flag.String("server", "127.0.0.1", "ip of the monitoring server")
	serverUdpPort := flag.String("udpport", "7071", "udp port on the server")
	serverTcpPort := flag.String("tcpport", "8081", "tcp port of the server")

	flag.Parse()

	client := model.NewClient(*serverIp, *serverTcpPort, *serverUdpPort)
	var counter uint64 = 0
	monitorProtocol.ConnectToServer(client.GetServerIp(), client.GetServerTcpPort())
	client.UpdateServerStatus(true)
	statsTimer := time.NewTicker(time.Second * 2).C
	aliveTimer := time.NewTicker(time.Second * 10).C

	// Choosing RuncStats or DockerStats
	statsObject := DockerStats{dockerContainerStats: docker.NewDockerContainerStats()}

	for {
		select {
		case <-statsTimer:
			{
				if client.GetServerStatus() {
					counter++
					statsObject.sendStats(client, counter)
				} else {
					fmt.Println("Server Offline .... Trying to Reconnect")
					monitorProtocol.ConnectToServer(client.GetServerIp(), client.GetServerTcpPort())
					client.UpdateServerStatus(true)
				}
			}
		case <-aliveTimer:
			{
				isAlive := checkIfServerIsAlive(client)
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

func (rcs *RuncStats) worker(client *model.Client, containerId string, containerStats chan monitorProtocol.ContainerStats, currentCounter uint64) {
	stats := runc.GetContainerStats(client, containerId)
	containerStats <- stats
}

func (rcs *RuncStats) sendStats(client *model.Client, counter uint64) {

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
	containerStats := make(chan monitorProtocol.ContainerStats, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		client.UpdateContainerCounter(containers[i], counter)
		go rcs.worker(client, containers[i], containerStats, counter)
	}

	//var buffer bytes.Buffer
	var statsToSend = make([]monitorProtocol.ContainerStats, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		//buffer.WriteString(<-containerStats)
		statsToSend[i] = <-containerStats
	}
	//stringToSend := buffer.String()

	monitorProtocol.SendContainerStatistics(statsToSend, client.GetServerIp(), client.GetServerUdpPort())
}

func (dcs *DockerStats) worker(client *model.Client, containerId string, containerStats chan monitorProtocol.ContainerStats, currentCounter uint64) {
	stats := docker.GetContainerStats(dcs.dockerContainerStats, containerId)
	containerStats <- stats
}

func (dcs *DockerStats) sendStats(client *model.Client, counter uint64) {

	dcs.dockerContainerStats.GetDockerStats()
	var containers []string = docker.GetRunningContainers(dcs.dockerContainerStats)
	numOfWorkers := len(containers)
	containerStats := make(chan monitorProtocol.ContainerStats, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		client.UpdateContainerCounter(containers[i], counter)
		go dcs.worker(client, containers[i], containerStats, counter)
	}

	//var buffer bytes.Buffer
	var statsToSend = make([]monitorProtocol.ContainerStats, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		//buffer.WriteString(<-containerStats)
		statsToSend[i] = <-containerStats
	}
	//stringToSend := buffer.String()

	monitorProtocol.SendContainerStatistics(statsToSend, client.GetServerIp(), client.GetServerUdpPort())

}
