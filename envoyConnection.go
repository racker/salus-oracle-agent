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
	"errors"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var host = "localhost"
var port = 8094
var retryLimit = 5


type iconnection interface {
	WriteToEnvoy(input string)
	Retry() error
}

type connection struct {
	conn net.Conn
	mux sync.Mutex
}


func (c *connection)  connect() error {

	if c.conn == nil {
		var err error
		c.conn, err = net.Dial("tcp", host+":"+strconv.FormatInt(int64(port), 10))
		if err != nil {
			log.Printf("Unable to create connection to Envoy: %s\n", err)
			return err
		}
	}

	return nil
}


func (c *connection) Retry() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	for i := 0; i < retryLimit; i++ {

		time.Sleep(10*time.Second)
		log.Println("Retrying connection to Envoy")
		err := c.connect()

		if err == nil {
			log.Println("Succeeded in connecting to Envoy")
			return nil
		}
	}
	return errors.New("unable to Connect to Envoy")
}


func (c *connection) WriteToEnvoy(input string) {
	errFlag := false
	c.mux.Lock()
	if c.conn != nil {
		_, err := c.conn.Write(append([]byte(input), []byte("\r\n")...))
		if err != nil {
			log.Printf("Could not write to Envoy: %s", err)
			err := c.Retry()
			if err != nil {
				log.Fatalf("Could not connect to Envoy: %s\n", err)
			}
			errFlag = true
		}
	} else {

		log.Println("Failed to send to Envoy: No Connection. Attempting to recreate connection.")
		err := c.Retry()
		if err != nil {
			log.Fatalf("Could not connect to Envoy: %s\n", err)
		}
		errFlag = true
	}
	c.mux.Unlock()
	if errFlag {
		c.WriteToEnvoy(input)
	}
}


