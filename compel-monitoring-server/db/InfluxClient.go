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

func GetConnection() influx.Client {

	// Create a new HTTPClient
	conn, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     "http://localhost:10090",
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

// func main() {
// 	i, err := strconv.ParseInt("1490057610", 10, 64)
// 	if err != nil {
// 		panic(err)
// 	}
// 	tm := time.Unix(i, 0)
// 	fmt.Println(tm)
//
// 	c, err := influx.NewHTTPClient(influx.HTTPConfig{
// 		Addr:     "http://localhost:10090",
// 		Username: username,
// 		Password: password,
// 	})
//
// 	defer c.Close()
//
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	//q := fmt.Sprintf("SELECT * FROM %s", "container_data")
// 	q := fmt.Sprintf("select * from container_data where agent = '192.168.0.26' ORDER BY time DESC LIMIT 6")
// 	res, err := queryDB(c, q)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//fmt.Println(res[0].Series[0].Values)
// 	//fmt.Println(res[0].Series[0].Values[0])
//
// 	for _, value := range res[0].Series[0].Values {
// 		fmt.Printf("%s : %s : %s \n", value[0], value[1], value[2])
// 	}
//
// 	//AddPoint("192.168.12.1", "mycontainer", 0, 0.00012064271, tm)
// }
