package main

import (
	"bytes"
	"time"

	model "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/model"
	runc "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func worker(client Client, containerId string, containerStats chan string) {
	stats := runc.GetContainerStats(containerId)
	containerStats <- stats
}

func sendStats(client Client) {
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

func main() {
	client := new(model.Client)
	monitorProtocol.ConnectToServer()
	statsTimer := time.NewTicker(time.Second * 2).C
	for {
		select {
		case <-statsTimer:
			sendStats(client)
		}
	}

}
