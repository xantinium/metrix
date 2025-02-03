package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xantinium/metrix/internal/agent"
)

func main() {
	agent := agent.NewMetrixAgent(agent.MetrixAgentOptions{
		ServerAddr:     "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
	})

	agent.Run()

	<-waitForStopSignal()
	agent.Stop()
}

func waitForStopSignal() <-chan os.Signal {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	return stopChan
}
