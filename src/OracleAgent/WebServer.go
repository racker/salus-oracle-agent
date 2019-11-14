package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

var host = "localhost"
var port = 8094
var conn net.Conn


func GetConn() net.Conn {
	if conn == nil {
		var err error

		conn, err = net.Dial("tcp", host+":"+strconv.FormatInt(int64(port), 10))
		if err != nil {
			log.Fatal("Unable to create connection to Envoy")
		}
	}

	return conn
}

/*
Here we will setup the web server to accept configuration files.
Unless we handle configuration files by just reading a particular directory of the filesystem
 */


func WriteToEnvoy(input []byte) {
	connection := GetConn()
	fmt.Fprintf(connection, "GET / HTTP/1.0\r\n\r\n")

	_, err := connection.Write(input)

	if err != nil {
		log.Fatal("Could not connect to Envoy")
	}
}


