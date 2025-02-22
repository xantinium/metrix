package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xantinium/metrix/internal/config"
	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/server"
)

func main() {
	var err error

	args := config.ParseServerArgs()

	logger.Init(args.IsDev)
	defer logger.Destroy()

	server := server.NewMetrixServer(server.MetrixServerOptions{
		Addr:           args.Addr,
		StoragePath:    args.StoragePath,
		StoreInterval:  args.StoreInterval,
		RestoreMetrics: args.RestoreStorage,
	})

	for {
		select {
		case err = <-server.Run():
			if err != nil {
				logger.Errorf("failed to run metrix server: %v", err)
				return
			}

			return
		case <-waitForStopSignal():
			err = server.Stop()
			if err != nil {
				logger.Errorf("failed to gracefully stop metrix server: %v", err)
			}

			return
		}
	}
}

func waitForStopSignal() <-chan os.Signal {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	return stopChan
}
