package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xantinium/metrix/internal/server"
)

func main() {
	var err error

	server := server.NewMetrixServer(8080)

	for {
		select {
		case err = <-server.Run():
			if err != nil {
				log.Fatal(fmt.Errorf("failed to run metrix server: %v", err))
				return
			}

			return
		case <-waitForStopSignal():
			err = server.Stop()
			if err != nil {
				log.Println(fmt.Errorf("failed to gracefully stop metrix server: %v", err))
			}
		}
	}
}

func waitForStopSignal() <-chan os.Signal {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	return stopChan
}
