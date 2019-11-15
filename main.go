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
var conn connection


func main() {

	err := conn.connect()
	if err != nil {
		conn.retry()
	}

	// This is here only till we get configuration management up to show that each one works
	tickers["dataguard"] = setupTimer(5, "dataguard.txt", processDataguard, createDataguardOutput)

	tickers["rman"]  = setupTimer(5, "RMAN.txt", processRMAN, createRMANOutput)

	tickers["tablespace"] = setupTimer(5, "tablespace.txt", processTablespace, createTablespaceOutput)


	select { } // make sure the application continues to run

}


func generateJSON(input telegrafJsonMetric) []byte {

	returnValue, err := json.Marshal(input)
	if err != nil {
		log.Fatal("Could not Marshal the telegrafJsonMetrics: ", err)
	}
	return returnValue
}

func createRMANOutput(input []string, fileName string) {
	var output telegrafJsonMetric
	output.Fields = make(map[string]interface{})
	output.Fields["error_codes"] = input
	output.Timestamp = time.Now()
	output.Name = "RMAN"
	if input == nil {
		output.Fields["status"] = "missing"
	}else {
		output.Fields["file_age"] = getFileInformation(fileName)
		output.Fields["status"] = "success"
	}
	conn.WriteToEnvoy(generateJSON(output))
}

//for tablespace we need to emit for every line in the file
func createTablespaceOutput(input []string, fileName string) {
	var output telegrafJsonMetric
	output.Fields = make(map[string]interface{})
	output.Tags = make(map[string]string)
	output.Timestamp = time.Now()
	output.Name = "Tablespace"

	if input == nil {
		output.Fields["status"] = "missing"
		conn.WriteToEnvoy(generateJSON(output))
	} else {
		output.Fields["file_age"] = getFileInformation(fileName)
		output.Fields["status"] = "success"
		for index, element := range input {
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

func createDataguardOutput(input []string, fileName string) {
	var output telegrafJsonMetric
	output.Fields = make(map[string]interface{})
	output.Timestamp = time.Now()
	output.Name = "dataguard"
	if input == nil {
		output.Fields["status"] = "missing"
	}else {
		output.Fields["file_age"] = getFileInformation(fileName)
		output.Fields["status"] = "success"
		output.Fields["replication"], _ = strconv.Atoi(input[0])
	}
	conn.WriteToEnvoy(generateJSON(output))
}

func getFileInformation(fileName string) time.Time {
	fileStat, err := os.Stat(fileName)
	if err != nil {
		log.Fatal("Unable to read file: ", err)
	}
	return fileStat.ModTime()
}

func readFile(fileName string, dispatch func(string)[]string) []string {
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
		var value = scanner.Text()
		var temp = dispatch(value)

		for _, element := range temp {
			output = append(output, element)
		}

	}

	return output
}


func processRMAN(input string) []string {
	var errorCode = regexp.MustCompile(`ORA-[0-9]+|RMAN-[0-9]+`)

	var returnValues = errorCode.FindAllString(input, -1)

	return returnValues
}

func processTablespace(input string) []string {
	values := strings.Split(input, ":")
	return values
}

func processDataguard(input string) []string {
	return []string{input}
}
