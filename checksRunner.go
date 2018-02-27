package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// ExecutionOutcome represents the outcome of a check
type ExecutionOutcome struct {
	Name        string
	Duration    time.Duration
	Environment string
	APIName     string
}

// Execute runs the actual checks
func Execute(checks []Check) []ExecutionOutcome {
	checksToRun := len(checks)
	result := make([]ExecutionOutcome, checksToRun)

	executionOutcomesChannel := make(chan ExecutionOutcome)

	for i := 0; i < checksToRun; i++ {
		go executeCheck(checks[i], executionOutcomesChannel)
	}

	for i := 0; i < checksToRun; i++ {
		result[i] = <-executionOutcomesChannel
	}

	return result
}

func executeCheck(check Check, executionOutcome chan ExecutionOutcome) {
	hc := http.Client{}
	var body io.Reader

	if check.Body != "" {
		b, err := os.Open(check.Body)

		if err != nil {
			log.Fatal(err)
		}

		body = b
	}

	req, _ := http.NewRequest(check.HTTPMethod, check.TargetURL, body)

	for i := 0; i < len(check.Headers); i++ {
		req.Header.Add(check.Headers[i].Name, check.Headers[i].Value)
	}

	t1 := time.Now()
	hc.Do(req)

	executionOutcome <- ExecutionOutcome{Name: check.Name, Duration: time.Now().Sub(t1), APIName: check.APIName, Environment: check.Env}
}
