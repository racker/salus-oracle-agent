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
	WriteToEnvoy(input []byte)
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
			log.Println("Succeeded in reconnecting to Envoy")
			return nil
		}
	}

	return errors.New("unable to Connect to Envoy")
	//return backoff.Retry(c.connect, backoff.NewExponentialBackOff())
}


func (c *connection) WriteToEnvoy(input []byte) {

	c.mux.Lock()
	if c.conn != nil {

		c.conn.Write([]byte("\r\n"))
		_, err := c.conn.Write(input)
		if err != nil {
			log.Printf("Could not write to Envoy: %s", err)
			err := c.Retry()
			if err != nil {
				log.Fatalf("Could not connect to Envoy: %s\n", err)
			}
		}
	}else {

		log.Println("Failed to send to Envoy: No Connection. Attempting to recreate connection")
		err := c.Retry()
		if err != nil {
			log.Fatalf("Could not connect to Envoy: %s\n", err)
		}
	}
	c.mux.Unlock()


}


