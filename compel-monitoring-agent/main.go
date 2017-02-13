package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	runc "github.com/adhuri/Compel-Monitoring/compel-monitoring-agent/runc"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func worker(index int, containerStats chan string) {
	time.Sleep(time.Second * 3)

	stats := " # Container " + strconv.Itoa(index)
	fmt.Println(stats)
	containerStats <- stats
}

func main() {
	done := make(chan bool)
	go monitorProtocol.ConnectToServer(done)
	<-done

	numOfWorkers := runc.NumberOfRunningContaiers()

	containerStats := make(chan string, numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		go worker(i, containerStats)
	}

	var buffer bytes.Buffer
	for i := 0; i < numOfWorkers; i++ {
		buffer.WriteString(<-containerStats)
	}
	stringToSend := buffer.String()

	monitorProtocol.SendContainerStatistics(stringToSend)

}
