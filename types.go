package main

import "time"

type monitor struct {
	interval int
	databaseName []string
	logFilePath string
}

type telegrafJsonMetric struct {
	Timestamp time.Time
	Name      string
	Tags      map[string]string
	Fields    map[string]interface{}
} // this is our output format

type DataguardConfiguration struct {
	Monitor monitor
}

type TablespaceConfiguration struct  {
	Monitor monitor
}

type RmanConfiguration struct {
	Monitor monitor
	errorCodeWhitelist []string
}

type monitorOutput func(processedData []string, fileName string)

type dispatchProcessing func(fileLine string) []string


