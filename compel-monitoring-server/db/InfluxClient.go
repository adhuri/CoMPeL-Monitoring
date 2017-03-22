package influx

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	MyDB     = "square_holes"
	username = "bubba"
	password = "bumblebeetuna"
)

func AddPoint(agentIp string, containerId string, cpuUsage float32, memoryUsage float32, timestamp time.Time) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:10090",
		Username: username,
		Password: password,
	})

	defer c.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
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

	pt, err := client.NewPoint("container_data", tags, fields, timestamp)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
}

// queryDB convenience function to query the database
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
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

func main() {
	i, err := strconv.ParseInt("1490057610", 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	fmt.Println(tm)

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:10090",
		Username: username,
		Password: password,
	})

	defer c.Close()

	if err != nil {
		log.Fatal(err)
	}

	//q := fmt.Sprintf("SELECT * FROM %s", "container_data")
	q := fmt.Sprintf("select * from container_data")
	res, err := queryDB(c, q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res[0].Series[0].Values[0][0])

	//AddPoint("192.168.12.1", "mycontainer", 0, 0.00012064271, tm)
}
