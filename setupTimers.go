package main

import (
	"log"
	"time"
)


func setupTimer(config Configuration, dispatch dispatchProcessing, dispatchOutput monitorOutput) time.Ticker {
	ticker := time.NewTicker( time.Duration(config.interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println(config.resolvePath())
				output, err := readFile(config.resolvePath(), config, dispatch)
				dispatchOutput(output, config.resolvePath(), err)
			}
		}
	}()

	return *ticker
}
