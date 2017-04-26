package db

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Monitoring/compel-monitoring-server/model"
	monitorProtocol "github.com/adhuri/Compel-Monitoring/protocol"
)

func StoreData(agentIp string, dataReceived []monitorProtocol.ContainerStats, server *model.Server, log *logrus.Logger) {

	influxServerIp := server.GetInfluxServer()
	influxPort := server.GetInfluxPort()

	if len(dataReceived) != 0 {
		startTime := time.Now()
		conn := GetConnection(influxServerIp, influxPort)
		for _, containerStat := range dataReceived {

			containerId := containerStat.ContainerID
			cpuUsage := containerStat.MetricData.CPU
			memoryUsage := containerStat.MetricData.Memory
			timestamp := containerStat.Timestamp

			dateTime := time.Unix(timestamp, 0)
			byteArray1 := []byte(containerId)
			size1 := len(byteArray1)
			//fmt.Println(size1)

			byteArray2 := []byte(agentIp)
			size2 := len(byteArray2)
			//fmt.Println(size2)
			AddPoint(agentIp, containerId, cpuUsage, memoryUsage, dateTime, conn)

			server.IncrementPointsSavedInDBCounterCounter()
			server.UpdateTotalDataWrittenToDB(int64(24 + size2 + size1))
		}
		conn.Close()
		elapsed := time.Since(startTime)
		server.UpdateDBWriteTime(elapsed)
		log.Infoln("Time taken to Save Container Data in INFLUX-DB: ", elapsed)
	} else {
		log.Infoln("No Data to Save in DB")
	}

}
