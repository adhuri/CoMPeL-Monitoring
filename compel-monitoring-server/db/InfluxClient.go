package db

import (
	"log"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

const (
	MyDB     = "square_holes"
	username = "bubba"
	password = "bumblebeetuna"
)

func GetConnection(ip string) influx.Client {

	// Create a new HTTPClient
	conn, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     "http://" + ip + ":10090",
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func AddPoint(agentIp string, containerId string, cpuUsage float64, memoryUsage float64, timestamp time.Time, conn influx.Client) {

	// Create a new point batch
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{
		"agent":     agentIp,
		"container": containerId,
	}
	fields := map[string]interface{}{
		"cpu":    cpuUsage,
		"memory": memoryUsage,
	}

	pt, err := influx.NewPoint("container_data", tags, fields, timestamp)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := conn.Write(bp); err != nil {
		log.Fatal(err)
	}
}

// queryDB convenience function to query the database
func queryDB(clnt influx.Client, cmd string) (res []influx.Result, err error) {
	q := influx.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
