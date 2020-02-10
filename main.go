package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var tickers = make(map[string]time.Ticker)
var conn iconnection = &connection{}
var timestamp iTimeInformation = &TimeInformation{}


func main() {
	configArg := flag.String("config",
		"./config.d",
		"Path to the configuration directory")
	flag.Parse()

	go exitApplication()

	err := conn.Retry()
	if err != nil {
		log.Fatalf("Failed to connect to Envoy: %s", err)
	}

	err = readConfigsFromPath(*configArg)
	if err != nil {
		log.Fatalf("Failed to read Config Files: %s", err)
		return
	}

	select { } // make sure the application continues to run

}


func generateJSON(input telegrafJsonMetric) string {

	returnValue, err := json.Marshal(input)
	if err != nil {
		log.Fatal("Could not Marshal the telegrafJsonMetrics: ", err)
	}
	return string(returnValue)
}

var createRMANOutput monitorOutput = func(processedData []string, fileName string, err error) {
	var output telegrafJsonMetric

	output.Fields = make(map[string]interface{})
	output.Fields["error_codes"] = processedData
	output.Timestamp = timestamp.Now()
	output.Name = TypeRMAN
	if err != nil {
		output.Fields["status"] = StatusMissing
	}else {
		output.Fields["file_age"] = timestamp.getFileInformation(fileName)
		output.Fields["status"] = StatusSuccess
	}
	conn.WriteToEnvoy(generateJSON(output))
}

//for tablespace we need to emit for every line in the file
var createTablespaceOutput monitorOutput = func(processedData []string, fileName string, err error) {
	var output telegrafJsonMetric
	output.Fields = make(map[string]interface{})
	output.Tags = make(map[string]string)
	output.Timestamp = timestamp.Now()
	output.Name = TypeTablespace

	if err != nil {
		output.Fields["status"] = StatusMissing
		conn.WriteToEnvoy(generateJSON(output))
		return
	} else if processedData == nil {
		// this is potentially a formatting issue with the database script we rely on to write the files we are monitoring
		output.Fields["status"] = "malformed"
		conn.WriteToEnvoy(generateJSON(output))
		return
	} else {
		output.Fields["file_age"] = timestamp.getFileInformation(fileName)
		output.Fields["status"] = StatusSuccess
		for index, element := range processedData {
			if index%2 == 0 { // even should be tablespace name
				output.Tags["tablespace_name"] = element
			} else {
				output.Fields["usage"] = element
				// We need to make sure we are sending it after both have been set
				conn.WriteToEnvoy(generateJSON(output))
			}
		}
		return
	}
}

var createDataguardOutput monitorOutput = func(processedData []string, fileName string, err error) {
	var output telegrafJsonMetric
	output.Fields = make(map[string]interface{})
	output.Timestamp = timestamp.Now()
	output.Name = TypeDataguard
	if err != nil {
		output.Fields["status"] = StatusMissing
	} else if processedData == nil {
		// this is potentially a formatting issue with the database script we rely on to write the files we are monitoring
		output.Fields["status"] = StatusMalformed
	} else {
		output.Fields["file_age"] = timestamp.getFileInformation(fileName)
		output.Fields["status"] = StatusSuccess
		output.Fields["replication"], _ = strconv.Atoi(processedData[0])
	}
	conn.WriteToEnvoy(generateJSON(output))
}

func readFile(fileName string, config Configuration, dispatch dispatchProcessing) ([]string, error ){
	file, err := os.Open(fileName)
	if err != nil {
		// We need to continue to send metrics through the system even if we can't read the file
		log.Printf("log file %s does not exist with error: %s", fileName, err)
		return nil, err

	}
	defer file.Close()

	var output []string
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Println("Unable to read file: ", err)
		return nil, err
	}
	for scanner.Scan() {
		var temp = dispatch(scanner.Text(), config)

		for _, element := range temp {
			output = append(output, element)
		}

	}

	return output, nil
}


var processRMAN dispatchProcessing = func(fileLine string, conf Configuration) []string {
	var errorCode = regexp.MustCompile(`ORA-[0-9]+|RMAN-[0-9]+`)

	var foundValues = errorCode.FindAllString(fileLine, -1)
	var returnValues = []string{}
	for _, i := range foundValues {
		var flag = true
		for _, i2 := range conf.errorCodeWhitelist {
			if strings.Compare(i, i2) == 0 {
				flag = false
			}
		}
		if flag {
			returnValues = append(returnValues, i)
		}
	}

	return returnValues
}

var processTablespace dispatchProcessing = func(fileLine string, conf Configuration) []string {
	values := strings.Split(fileLine, ":")
	for index, element := range values {
		values[index] = strings.TrimSpace(element)
	}

	return values
}

var processDataguard dispatchProcessing = func(fileLine string, conf Configuration) []string {
	return []string{fileLine}
}


func exitApplication() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)


	// Block until a signal is received.
	v := <- c
	log.Println("Successfully received signal: ", v)
	// stop running monitors

	log.Println("Gracefully shutting down monitors")
	stopMonitors()
	log.Println("Monitors successfully shut down")
	os.Exit(0)
}


func stopMonitors() {
	for _, value := range tickers {
		value.Stop()
	}
}


func readConfigsFromPath(configDir string) error {
	return filepath.Walk(configDir, readConfig)

}

func readConfig(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("Error reading file: ", err)
		return err
	}

	if info.IsDir() {
		log.Println("Oracle Agent doesn't recursively check directories for configurations: ", path)
		return err
	}

	fileContents, _ := ioutil.ReadFile(path)

	var inputConfig InputConfiguration
	err = json.Unmarshal(fileContents, &inputConfig)
	if(err != nil) {
		log.Println("Could not unmarshal configuration: ", err)
	}

	var config Configuration
	config.interval = inputConfig.Interval
	config.filePath = inputConfig.FilePath
	config.configType = inputConfig.Type
	var dispatch dispatchProcessing
	var output monitorOutput
	switch inputConfig.Type {
		case TypeDataguard:
			dispatch = processDataguard
			output = createDataguardOutput
		case TypeTablespace:
			dispatch = processTablespace
			output = createTablespaceOutput
		case TypeRMAN:
			dispatch = processRMAN
			output = createRMANOutput
			config.errorCodeWhitelist = inputConfig.ErrorCodeWhitelist
	default:
		log.Println("Unable to determine type of configuration: ", inputConfig.Type)
		return nil
	}

	for _, database := range inputConfig.DatabaseNames {
		config.databaseName = database
		tickers[info.Name()+":"+database] = setupTimer(config, dispatch, output)
	}



	return nil
}

