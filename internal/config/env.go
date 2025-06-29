package config

import (
	"log"
	"net"
	"time"

	"github.com/xantinium/metrix/internal/tools"
)

// parseServerArgsFromEnv парсит переменные окружения в optionalServerArgs.
func parseServerArgsFromEnv() optionalServerArgs {
	address := tools.GetStrFromEnv("ADDRESS")
	privateKey := tools.GetStrFromEnv("KEY")
	cryptoPrivateKey := tools.GetStrFromEnv("CRYPTO_KEY")
	storeInterval := tools.GetIntFromEnv("STORE_INTERVAL")
	storagePath := tools.GetStrFromEnv("FILE_STORAGE_PATH")
	restoreStorage := tools.GetBoolFromEnv("RESTORE")
	databaseConnStr := tools.GetStrFromEnv("DATABASE_DSN")
	shutdownTimeout := tools.GetIntFromEnv("SHUTDOWN_TIMEOUT")
	trustedSubnet := tools.GetStrFromEnv("TRUSTED_SUBNET")

	args := optionalServerArgs{}

	if address.Exists {
		args.Addr = &address.Value
	}
	if privateKey.Exists {
		args.PrivateKey = &privateKey.Value
	}
	if cryptoPrivateKey.Exists {
		args.CryptoPrivateKey = &cryptoPrivateKey.Value
	}
	if storeInterval.Exists {
		tmp := time.Duration(storeInterval.Value) * time.Second
		args.StoreInterval = &tmp
	}
	if storagePath.Exists {
		args.StoragePath = &storagePath.Value
	}
	if restoreStorage.Exists {
		args.RestoreStorage = &restoreStorage.Value
	}
	if databaseConnStr.Exists {
		args.DatabaseConnStr = &databaseConnStr.Value
	}
	if shutdownTimeout.Exists {
		tmp := time.Duration(shutdownTimeout.Value) * time.Second
		args.ShutdownTimeout = &tmp
	}
	if trustedSubnet.Exists {
		var err error
		_, args.TrustedSubnet, err = net.ParseCIDR(trustedSubnet.Value)
		if err != nil {
			log.Printf("failed to parse trusted subnet: %v\n", err)
		}
	}

	return args
}

// parseAgentArgsFromEnv парсит переменные окружения в optionalAgentArgs.
func parseAgentArgsFromEnv() optionalAgentArgs {
	address := tools.GetStrFromEnv("ADDRESS")
	privateKey := tools.GetStrFromEnv("KEY")
	cryptoPublicKey := tools.GetStrFromEnv("CRYPTO_KEY")
	pollInterval := tools.GetIntFromEnv("POLL_INTERVAL")
	reportInterval := tools.GetIntFromEnv("REPORT_INTERVAL")
	reportRateLimit := tools.GetIntFromEnv("RATE_LIMIT")
	shutdownTimeout := tools.GetIntFromEnv("SHUTDOWN_TIMEOUT")

	args := optionalAgentArgs{}

	if address.Exists {
		args.Addr = &address.Value
	}
	if privateKey.Exists {
		args.PrivateKey = &privateKey.Value
	}
	if cryptoPublicKey.Exists {
		args.CryptoPublicKey = &cryptoPublicKey.Value
	}
	if pollInterval.Exists {
		args.PollInterval = &pollInterval.Value
	}
	if reportInterval.Exists {
		tmp := time.Duration(reportInterval.Value) * time.Second
		args.ReportInterval = &tmp
	}
	if reportRateLimit.Exists {
		args.ReportRateLimit = &reportRateLimit.Value
	}
	if shutdownTimeout.Exists {
		tmp := time.Duration(shutdownTimeout.Value) * time.Second
		args.ShutdownTimeout = &tmp
	}

	return args
}
