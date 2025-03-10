package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/xantinium/metrix/internal/config"
	"github.com/xantinium/metrix/internal/infrastructure/memstorage"
	"github.com/xantinium/metrix/internal/infrastructure/postgres"
	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/server"
)

func main() {
	args := config.ParseServerArgs()

	logger.Init(args.IsDev)
	defer logger.Destroy()

	server, cleanUp, err := getMetrixServer(args)
	if err != nil {
		panic(err)
	}
	defer cleanUp(context.TODO())

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

type cleanUpFunc = func(context.Context)

type emptyDBChecker struct{}

func (emptyDBChecker) Ping(_ context.Context) error {
	return nil
}

func getMetrixServer(args config.ServerArgs) (*server.MetrixServer, cleanUpFunc, error) {
	if args.DatabaseConnStr != "" {
		psqlClient, err := postgres.NewPostgresClient(args.DatabaseConnStr)
		if err != nil {
			return nil, nil, err
		}

		return server.NewMetrixServer(server.MetrixServerOptions{
			Addr:          args.Addr,
			StoreInterval: args.StoreInterval,
			Storage:       psqlClient,
			DBChecker:     psqlClient,
		}), psqlClient.Destroy, nil
	}

	memStorage, err := memstorage.NewMemStorage(args.StoragePath, args.RestoreStorage)
	if err != nil {
		return nil, nil, err
	}

	return server.NewMetrixServer(server.MetrixServerOptions{
		Addr:          args.Addr,
		StoreInterval: args.StoreInterval,
		Storage:       memStorage,
		DBChecker:     new(emptyDBChecker),
	}), memStorage.Destroy, nil
}

func waitForStopSignal() <-chan os.Signal {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	return stopChan
}
