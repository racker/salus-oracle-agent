package main

import (
	"path/filepath"
)

type telegrafJsonMetric struct {
	Timestamp int64
	Name      string
	Tags      map[string]string
	Fields    map[string]interface{}
} // this is our output format

type Configuration struct {
	interval int
	configType string `json:"type"`
	databaseName string `json:"databaseNames"`
	filePath string `json:"filePath"`
	errorCodeWhitelist []string
}

type InputConfiguration struct {
	Type          string   		`json:"type"`
	DatabaseNames []string 		`json:"databaseNames"`
	FilePath      string   		`json:"filePath"`
	Interval      int      		`json:"interval"`
	ErrorCodeWhitelist []string `json:errorCodeWhitelist`
}

func (c *Configuration) resolvePath() string {
	return filepath.Join(c.filePath,c.databaseName+".txt")
}

type monitorOutput func(processedData []string, fileName string, err error)

type dispatchProcessing func(fileLine string, conf Configuration) []string


