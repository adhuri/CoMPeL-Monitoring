package main

import (
	"bytes"

	runc "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func worker(containerId string, containerStats chan string) {
	stats := runc.GetContainerStats(containerId)
	containerStats <- stats
}

func main() {
	done := make(chan bool)
	go monitorProtocol.ConnectToServer(done)
	<-done

	var containers []string = runc.GetRunningContaiers()
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
