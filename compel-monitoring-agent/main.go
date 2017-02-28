package main

import (
	"bytes"
	"time"

	model "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/model"
	runc "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func worker(client model.Client, containerId string, containerStats chan string, currentCounter uint64) {
	stats := runc.GetContainerStats(containerId)
	containerStats <- stats
}

func sendStats(client model.Client, counter uint64) {
	var containers []string = runc.GetRunningContainers()
	numOfWorkers := len(containers)
	containerStats := make(chan string, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		client.UpdateContainerCounter(containers[i], counter)
		go worker(client,containers[i], containerStats, counter)
	}

	var buffer bytes.Buffer
	for i := 0; i < numOfWorkers; i++ {
		buffer.WriteString(<-containerStats)
	}
	stringToSend := buffer.String()

	monitorProtocol.SendContainerStatistics(stringToSend)
}

func main() {
	client := *new(model.Client)
	var counter uint64 = 0
	monitorProtocol.ConnectToServer()
	statsTimer := time.NewTicker(time.Second * 2).C
	for {
		select {
		case <-statsTimer:
			{
				counter++
				sendStats(client, counter)
			}
		}
	}

}
