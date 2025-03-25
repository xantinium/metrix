package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xantinium/metrix/internal/agent"
	"github.com/xantinium/metrix/internal/config"
	"github.com/xantinium/metrix/internal/logger"
)

func main() {
	args := config.ParseAgentArgs()

	logger.Init(args.IsDev)
	defer logger.Destroy()

	agent := agent.NewMetrixAgent(agent.MetrixAgentOptions{
		ServerAddr:     args.Addr,
		PrivateKey:     args.PrivateKey,
		PollInterval:   args.PollInterval,
		ReportInterval: args.ReportInterval,
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
