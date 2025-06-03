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
	builder := server.NewMetrixServerBuilder().
		SetAddr(args.Addr).
		SetPrivateKey(args.PrivateKey).
		SetStoreInterval(args.StoreInterval)

	// Если строка подключения к БД отсутствует,
	// используем in-memory хранилище и моковый DBChecker.
	if args.DatabaseConnStr == "" {
		memStorage, err := memstorage.NewMemStorage(args.StoragePath, args.RestoreStorage)
		if err != nil {
			return nil, nil, err
		}

		builder.SetStorage(memStorage, new(emptyDBChecker))

		return builder.Build(), memStorage.Destroy, nil
	}

	psqlClient, err := postgres.NewPostgresClient(ctx, args.DatabaseConnStr)
	if err != nil {
		return nil, nil, err
	}

	builder.SetStorage(psqlClient, psqlClient)

	return builder.Build(), psqlClient.Destroy, nil
}

func waitForStopSignal() <-chan os.Signal {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	return stopChan
}
