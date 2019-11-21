package main

import (
	"time"
)


func updateConf() {
	//grab an instance of the ticker associated with the config.
	//stop the ticker
	//create a new ticker with the new configuration information.
}



func setupTimer(interval int, fileName string, dispatch dispatchProcessing, dispatchOutput monitorOutput) time.Ticker {
	ticker := time.NewTicker( time.Duration(interval) * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				var output = readFile(fileName, dispatch)
				dispatchOutput(output, fileName)
			}
		}
	}()

	return *ticker
}
