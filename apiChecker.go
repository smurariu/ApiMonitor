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
	groupedResults := make(map[string]map[string][]ExecutionOutcome)

	for _, result := range results {
		if groupedResults[result.Environment] == nil {
			groupedResults[result.Environment] = make(map[string][]ExecutionOutcome)
		}

		groupedResults[result.Environment][result.APIName] = append(groupedResults[result.Environment][result.APIName], result)
	}
	// Create a new points batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  myDB,
		Precision: "s",
	})

	if err != nil {
		log.Fatal(err)
	}

	for env, envCheckResults := range groupedResults {
		tags := map[string]string{"env": env}
		for apiName, apiCheckResults := range envCheckResults {
			fields := make(map[string]interface{})

			for _, apiCheckResult := range apiCheckResults {
				// Create a point and add to batch
				fields[apiCheckResult.Name] = int64(apiCheckResult.Duration / time.Millisecond)
			}

			pt, err := client.NewPoint(apiName, tags, fields, time.Now())

			if err != nil {
				log.Fatal(err)
			}

			bp.AddPoint(pt)
		}
	}

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
}
