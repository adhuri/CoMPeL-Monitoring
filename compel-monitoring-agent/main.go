package main

import (
	"bytes"
	"time"

	runc "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func worker(containerId string, containerStats chan string) {
	stats := runc.GetContainerStats(containerId)
	containerStats <- stats
}

func sendStats() {
	var containers []string = runc.GetRunningContainers()
	numOfWorkers := len(containers)

	containerStats := make(chan string, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		go worker(containers[i], containerStats)
	}

	var buffer bytes.Buffer
	for i := 0; i < numOfWorkers; i++ {
		buffer.WriteString(<-containerStats)
	}
	stringToSend := buffer.String()

	monitorProtocol.SendContainerStatistics(stringToSend)
}

func main() {
	monitorProtocol.ConnectToServer()
	statsTimer := time.NewTicker(time.Second * 2).C
	for {
		select {
		case <-statsTimer:
			sendStats()
		}
	}

}
