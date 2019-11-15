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
	generalInfo monitor
}

type TablespaceConfiguration struct  {
	generalInfo monitor
}

type RmanConfiguration struct {
	generalInfo monitor
	errorCodeWhitelist []string
}

type monitorOutput func([]string, string)

type dispatchProcessing func(string) []string


