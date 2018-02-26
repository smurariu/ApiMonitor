package playerCheck

import (
	"log"
	"net/http"
	"os"
	"time"
)

const (
	apiURL      = "http://10.110.22.16:81/Bewo.Hub.Service/BewotecHubService.svc"
	contentType = "text/xml; charset=utf-8"
	soapAction  = "http://www.bewotec.de/bewotecws/Schema/BewotecHubService/GetProductOffers"
)

// func main() {
// 	executionTimes := make(chan Execution)

// 	go timeExecution(doPing, executionTimes, "ping")
// 	go timeExecution(doRequest, executionTimes, "getPrices")

// 	execution1 := <-executionTimes
// 	execution2 := <-executionTimes

// 	fmt.Printf("%s: %s\n", execution1.name, execution1.duration)
// 	fmt.Printf("%s: %s\n", execution2.name, execution2.duration)
// }

// ExecuteChecks runs the checks on the external system
func ExecuteChecks() (results [2]Execution) {
	executionTimes := make(chan Execution)

	go timeExecution(doPing, executionTimes, "ping")
	go timeExecution(doRequest, executionTimes, "getPrices")

	execution1 := <-executionTimes
	execution2 := <-executionTimes

	return [2]Execution{execution1, execution2}
}

type executionTime func()

// Execution represents the outcome of a check
type Execution struct {
	Name     string
	Duration time.Duration
}

func timeExecution(d executionTime, executionTimes chan Execution, name string) {
	t1 := time.Now()
	d()
	executionTimes <- Execution{Name: name, Duration: time.Now().Sub(t1)}
}

func doRequest() {
	hc := http.Client{}
	body, err := os.Open("playerCheck/request.xml")
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("POST", apiURL, body)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("SOAPAction", soapAction)

	hc.Do(req)
	defer body.Close()
}

func doPing() {
	hc := http.Client{}
	req, _ := http.NewRequest("POST", apiURL, nil)
	hc.Do(req)
	// resp, _ := hc.Do(req)
	// fmt.Println(resp.StatusCode)
}
