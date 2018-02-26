package main

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	myDB     = "db0"
	username = "apiChecker"
	password = "SundioApiChecker1234"
)

func main() {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://10.10.18.167:8086",
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(5000 * time.Millisecond)
	for t := range ticker.C {
		fmt.Print(t)
		checks := LoadChecks("checks.json")
		checksOutcomes := Execute(checks)

		writeToInflux(c, checksOutcomes)
		fmt.Printf(" %s\n", checksOutcomes)
	}
}

func writeToInflux(c client.Client, results []ExecutionOutcome) {

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  myDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"env": "ACPT"}
	fields := map[string]interface{}{
		results[0].Name: int64(results[0].Duration / time.Millisecond),
		results[1].Name: int64(results[1].Duration / time.Millisecond),
	}

	pt, err := client.NewPoint("Player", tags, fields, time.Now())

	if err != nil {
		log.Fatal(err)
	}

	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
}
