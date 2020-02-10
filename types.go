/*
 * Copyright 2020 Rackspace US, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	StatusMissing 	= "missing"
	StatusSuccess 	= "success"
	StatusMalformed = "malformed"
	TypeDataguard	= "oracle_dataguard"
	TypeTablespace	= "oracle_tablespace"
	TypeRMAN		= "oracle_rman"
)

type telegrafJsonMetric struct {
	Timestamp int64
	Name      string
	Tags      map[string]string
	Fields    map[string]interface{}
} // this is our output format

type Configuration struct {
	interval int
	configType string
	databaseName string
	filePath string
	errorCodeWhitelist []string
}

func (c *Configuration) resolvePath() string {
	return filepath.Join(c.filePath,c.databaseName+".txt")
}

type InputConfiguration struct {
	Type          string   		`json:"type"`
	DatabaseNames []string 		`json:"databaseNames"`
	FilePath      string   		`json:"filePath"`
	Interval      int      		`json:"interval"`
	ErrorCodeWhitelist []string `json:errorCodeWhitelist`
}

type monitorOutput func(processedData []string, fileName string, err error)

type dispatchProcessing func(fileLine string, conf Configuration) []string

type iTimeInformation interface {
	// the telegraf data format is expecting timestamps as int64
	Now() int64
	getFileInformation(string) int64
}

type TimeInformation struct {}

func (t *TimeInformation) Now() int64 {
	// Our system is expecting timestamps in milliseconds
	return time.Now().Unix() * 1000
}

func (t *TimeInformation) getFileInformation(fileName string) int64 {
	fileStat, err := os.Stat(fileName)
	if err != nil {
		log.Fatal("Unable to read file: ", err)
	}
	// Our system is expecting timestamps in milliseconds
	return fileStat.ModTime().Unix() * 1000
}
