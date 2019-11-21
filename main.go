package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var tickers = make(map[string]time.Ticker)
var conn iconnection = &connection{}
var timestamp iTimeInformation = &TimeInformation{}


type iTimeInformation interface {
	Now() time.Time
	getFileInformation(string) time.Time
}

type TimeInformation struct {}

func (t *TimeInformation) Now() time.Time {
	return time.Now()
}

func (t *TimeInformation) getFileInformation(fileName string) time.Time {
	fileStat, err := os.Stat(fileName)
	if err != nil {
		log.Fatal("Unable to read file: ", err)
	}
	return fileStat.ModTime()
}


func main() {

	err := conn.Retry()
	if err != nil {
		log.Fatalf("Failed to connect to Envoy: %s", err)
	}

	// This is here only till we get configuration management up to show that each one works
	tickers["dataguard"] = setupTimer(5, "./testdata/dataguard.txt", processDataguard, createDataguardOutput)

	tickers["rman"]  = setupTimer(5, "./testdata/RMAN.txt", processRMAN, createRMANOutput)

	tickers["tablespace"] = setupTimer(5, "./testdata/tablespace.txt", processTablespace, createTablespaceOutput)


	select { } // make sure the application continues to run

}


func generateJSON(input telegrafJsonMetric) []byte {

	returnValue, err := json.Marshal(input)
	if err != nil {
		log.Fatal("Could not Marshal the telegrafJsonMetrics: ", err)
	}
	return returnValue
}

var createRMANOutput monitorOutput = func(processedData []string, fileName string) {
	var output telegrafJsonMetric
	output.Fields = make(map[string]interface{})
	output.Fields["error_codes"] = processedData
	output.Timestamp = timestamp.Now()
	output.Name = "RMAN"
	if processedData == nil {
		output.Fields["status"] = "missing"
	}else {
		output.Fields["file_age"] = timestamp.getFileInformation(fileName)
		output.Fields["status"] = "success"
	}
	conn.WriteToEnvoy(generateJSON(output))
}

//for tablespace we need to emit for every line in the file
var createTablespaceOutput monitorOutput = func(processedData []string, fileName string) {
	var output telegrafJsonMetric
	output.Fields = make(map[string]interface{})
	output.Tags = make(map[string]string)
	output.Timestamp = timestamp.Now()
	output.Name = "Tablespace"

	if processedData == nil {
		output.Fields["status"] = "missing"
		conn.WriteToEnvoy(generateJSON(output))
	} else {
		output.Fields["file_age"] = timestamp.getFileInformation(fileName)
		output.Fields["status"] = "success"
		for index, element := range processedData {
			if index%2 == 0 { // even should be tablespace name
				output.Tags["tablespace_name"] = element
			} else {
				output.Fields["usage"] = element
				// We need to make sure we are sending it after both have been set
				conn.WriteToEnvoy(generateJSON(output))
			}
		}
	}
}

var createDataguardOutput monitorOutput = func(processedData []string, fileName string) {
	var output telegrafJsonMetric
	output.Fields = make(map[string]interface{})
	output.Timestamp = timestamp.Now()
	output.Name = "dataguard"
	if processedData == nil {
		output.Fields["status"] = "missing"
	}else {
		output.Fields["file_age"] = timestamp.getFileInformation(fileName)
		output.Fields["status"] = "success"
		output.Fields["replication"], _ = strconv.Atoi(processedData[0])
	}
	conn.WriteToEnvoy(generateJSON(output))
}

func readFile(fileName string, dispatch dispatchProcessing) []string {
	file, err := os.Open(fileName)
	if err != nil {
		// We need to continue to send metrics through the system even if we can't read the file
		log.Printf("log file %s does not exist", fileName)
		return nil

	}
	defer file.Close()

	var output []string
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Println("Unable to read file: ", err)
		return nil
	}
	for scanner.Scan() {
		var temp = dispatch(scanner.Text())

		for _, element := range temp {
			output = append(output, element)
		}

	}

	return output
}


var processRMAN dispatchProcessing = func(fileLine string) []string {
	var errorCode = regexp.MustCompile(`ORA-[0-9]+|RMAN-[0-9]+`)

	var returnValues = errorCode.FindAllString(fileLine, -1)

	return returnValues
}

var processTablespace dispatchProcessing = func(fileLine string) []string {
	values := strings.Split(fileLine, ":")
	for index, element := range values {
		values[index] = strings.TrimSpace(element)
	}

	return values
}

var processDataguard dispatchProcessing = func(fileLine string) []string {
	return []string{fileLine}
}

