package main

import (
	"flag"
	"net"
	"os"
	"time"

	logrus "github.com/Sirupsen/logrus"

	docker "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/docker"
	model "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/model"
	runc "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc"
	stats "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc/stats"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

var (
	log *logrus.Logger
)

func init() {

	log = logrus.New()

	// Output logging to stdout
	log.Out = os.Stdout

	// Only log the info severity or above.
	log.Level = logrus.InfoLevel

	// Microseconds level logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	customFormatter.FullTimestamp = true

	log.Formatter = customFormatter

}

func worker(client *model.Client, containerId string, containerStats chan monitorProtocol.ContainerStats, currentCounter uint64) {
	stats := runc.GetContainerStats(client, containerId, log)
	containerStats <- stats
}

func sendStats(client *model.Client, counter uint64) {

	//Set SystemCPU usage
	sysCPUusage, err := stats.GetSystemCPU(log)
	if err != nil {
		log.Errorln("Cannot Get System CPU")
	} else {
		client.SetTotalCPU(sysCPUusage)
	}
	//Set Memory Limit
	sysMemoryLimit, err := stats.GetSystemMemory(log)
	if err != nil {
		log.Errorln("Cannot Get System Memory")
	} else {
		client.SetTotalMemory(sysMemoryLimit)
	}

	var containers []string = runc.GetRunningContainers(log)
	numOfWorkers := len(containers)
	containerStats := make(chan monitorProtocol.ContainerStats, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		client.UpdateContainerCounter(containers[i], counter)
		go worker(client, containers[i], containerStats, counter)
	}

	//var buffer bytes.Buffer
	var statsToSend = make([]monitorProtocol.ContainerStats, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		//buffer.WriteString(<-containerStats)
		statsToSend[i] = <-containerStats
	}
	//stringToSend := buffer.String()

	monitorProtocol.SendContainerStatistics(statsToSend, client.GetServerIp(), client.GetServerUdpPort(), log)
}

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
	// Read command line arguments
	serverIp := flag.String("server", "127.0.0.1", "ip of the monitoring server")
	serverUdpPort := flag.String("udpport", "7071", "udp port on the server")
	serverTcpPort := flag.String("tcpport", "8081", "tcp port of the server")
	flag.Parse()
	log.WithFields(logrus.Fields{
		"serverIp":      *serverIp,
		"serverUdpPort": *serverUdpPort,
		"serverTcpPort": *serverTcpPort,
	}).Info("Inputs from command line")

	// Connect to monitoring server
	client := model.NewClient(*serverIp, *serverTcpPort, *serverUdpPort)
	monitorProtocol.ConnectToServer(client.GetServerIp(), client.GetServerTcpPort(), log)

	// After successful connection update flag on client
	client.UpdateServerStatus(true)

	// Choosing RuncStats or DockerStats
	statsObject := DockerStats{dockerContainerStats: docker.NewDockerContainerStats()}

	// Initialise Stats Timer
	statsTimer := time.NewTicker(time.Second * 2).C
	aliveTimer := time.NewTicker(time.Second * 10).C
	var counter uint64 = 0

	for {
		select {
		case <-statsTimer:
			{
				// Refresh object
				//*statsObject.ClearDockerContainerList()
				statsObject.dockerContainerStats.ClearDockerContainerList()

				if client.GetServerStatus() {
					counter++
					statsObject.sendStats(client, counter)
				} else {
					log.Warnln("Server Offline .... Trying to Reconnect")
					monitorProtocol.ConnectToServer(client.GetServerIp(), client.GetServerTcpPort(), log)
					client.UpdateServerStatus(true)
				}
			}
		case <-aliveTimer:
			{
				isAlive := checkIfServerIsAlive(client)
				if !isAlive {
					// update the server status
					log.Errorln("Server Dead")
					client.UpdateServerStatus(false)
				} else {
					log.Infoln("Server is still Alive")
				}
			}
		}
	}

}

func (rcs *RuncStats) worker(client *model.Client, containerId string, containerStats chan monitorProtocol.ContainerStats, currentCounter uint64) {
	stats := runc.GetContainerStats(client, containerId, log)
	containerStats <- stats
}

func (rcs *RuncStats) sendStats(client *model.Client, counter uint64) {

	//Set SystemCPU usage
	sysCPUusage, err := stats.GetSystemCPU(log)
	if err != nil {
		log.Errorln("Cannot Get System CPU")
	} else {
		client.SetTotalCPU(sysCPUusage)
	}
	//Set Memory Limit
	sysMemoryLimit, err := stats.GetSystemMemory(log)
	if err != nil {
		log.Errorln("Cannot Get System Memory")
	} else {
		client.SetTotalMemory(sysMemoryLimit)
	}

	var containers []string = runc.GetRunningContainers(log)
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

	monitorProtocol.SendContainerStatistics(statsToSend, client.GetServerIp(), client.GetServerUdpPort(), log)
}

func (dcs *DockerStats) worker(client *model.Client, containerId string, containerStats chan monitorProtocol.ContainerStats, currentCounter uint64) {
	stats := docker.GetContainerStats(dcs.dockerContainerStats, containerId, log)
	containerStats <- stats
}

func (dcs *DockerStats) sendStats(client *model.Client, counter uint64) {

	dcs.dockerContainerStats.GetDockerStats(log)
	var containers []string = docker.GetRunningContainers(dcs.dockerContainerStats, log)
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

	monitorProtocol.SendContainerStatistics(statsToSend, client.GetServerIp(), client.GetServerUdpPort(), log)

}
