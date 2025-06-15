package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/xantinium/metrix/internal/agent"
	"github.com/xantinium/metrix/internal/config"
	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/tools"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	tools.PrintBuildInfo(tools.BuildInfo{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	})

	args := config.ParseAgentArgs()

	logger.Init(args.IsDev)
	defer logger.Destroy()

	agent := agent.NewMetrixAgent(agent.MetrixAgentOptions{
		ServerAddr:         args.Addr,
		PrivateKey:         args.PrivateKey,
		CryptoPublicKey:    args.CryptoPublicKey,
		PollInterval:       args.PollInterval,
		ReportInterval:     args.ReportInterval,
		ReportRateLimit:    args.ReportRateLimit,
		IsProfilingEnabled: args.IsProfilingEnabled,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agent.Run(ctx)

	<-waitForStopSignal()
}

func waitForStopSignal() <-chan os.Signal {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	return stopChan
}
