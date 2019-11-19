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

func (m *telegrafJsonMetric) compareTo(value telegrafJsonMetric) bool {
	//ignore timestamp and Fields[file_age] but we want to make sure that Fields[file_age] exists and that it is a timestamp type
	if m.Name != value.Name {
		return false
	}

	if m.Tags != nil && value.Tags != nil {
		//iterate through
		return false
	}

	return true
}

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

type monitorOutput func([]string, string)

type dispatchProcessing func(string) []string


