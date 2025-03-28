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
	ctx := context.Background()
	args := config.ParseServerArgs()

	logger.Init(args.IsDev)
	defer logger.Destroy()

	server, cleanUp, err := getMetrixServer(ctx, args)
	if err != nil {
		panic(err)
	}
	defer cleanUp(ctx)

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

func getMetrixServer(ctx context.Context, args config.ServerArgs) (*server.MetrixServer, cleanUpFunc, error) {
	if args.DatabaseConnStr != "" {
		psqlClient, err := postgres.NewPostgresClient(ctx, args.DatabaseConnStr)
		if err != nil {
			return nil, nil, err
		}

		return server.NewMetrixServer(server.MetrixServerOptions{
			Addr:          args.Addr,
			PrivateKey:    args.PrivateKey,
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
		PrivateKey:    args.PrivateKey,
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
