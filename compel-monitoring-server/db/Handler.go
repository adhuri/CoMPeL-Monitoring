package db

import (
	"time"

	"github.com/Sirupsen/logrus"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func StoreData(agentIp string, dataReceived []monitorProtocol.ContainerStats, influxServerIp string, influxPort string, log *logrus.Logger) {

	if len(dataReceived) != 0 {
		startTime := time.Now()
		conn := GetConnection(influxServerIp, influxPort)
		for _, containerStat := range dataReceived {

			containerId := containerStat.ContainerID
			cpuUsage := containerStat.MetricData.CPU
			memoryUsage := containerStat.MetricData.Memory
			timestamp := containerStat.Timestamp

			dateTime := time.Unix(timestamp, 0)
			AddPoint(agentIp, containerId, cpuUsage, memoryUsage, dateTime, conn)
		}
		conn.Close()
		elapsed := time.Since(startTime)
		log.Infoln("Time taken to Save Container Data in INFLUX-DB: ", elapsed)
	} else {
		log.Infoln("No Data to Save in DB")
	}

}
