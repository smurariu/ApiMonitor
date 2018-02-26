package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//Header defines a HTTP header to add to the request
type Header struct {
	Name  string
	Value string
}

//Check defines a check to be performed
type Check struct {
	Name       string
	TargetURL  string
	HTTPMethod string
	Headers    []Header
	Body       string
}

//LoadChecks reads the check to be performed from a file
func LoadChecks(filename string) []Check {
	var objs []Check
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Printf("File Error: [%v]\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, &objs)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return objs
}
