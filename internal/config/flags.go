package config

import (
	"flag"
	"log"
	"net"
	"time"
)

// parseServerArgsFromFlags парсит флаги в optionalServerArgs.
func parseServerArgsFromFlags() optionalServerArgs {
	address := new(netAddress)
	flag.Var(address, "a", "address of metrix server in form <host:port>")
	isDev := flag.Bool("dev", false, "is metrix server running in development mode")
	isProfilingEnabled := flag.Bool("profile", false, "is profiling via pprof enabled")
	privateKey := flag.String("k", "", "key for hash funcs")
	cryptoPrivateKey := flag.String("crypto-key", "", "private key for crypto funcs in HEX")
	storeInterval := flag.Int("i", 300, "interval (in seconds) of writing metrics into file")
	storagePath := flag.String("f", "./metrix.db", "path to file for metrics writing")
	restoreStorage := flag.Bool("r", true, "read metrics from file on start")
	databaseConnStr := flag.String("d", "", "connection string for postgresql")
	trustedSubnet := flag.String("t", "", "defines the allowed IP subnet in CIDR notation (e.g., \"192.168.1.0/24\")")

	flag.Parse()

	args := optionalServerArgs{
		StoragePath:        storagePath,
		PrivateKey:         privateKey,
		CryptoPrivateKey:   cryptoPrivateKey,
		DatabaseConnStr:    databaseConnStr,
		IsDev:              isDev,
		IsProfilingEnabled: isProfilingEnabled,
		RestoreStorage:     restoreStorage,
	}

	{
		tmp := address.String()
		args.Addr = &tmp
	}

	if storeInterval != nil {
		tmp := time.Duration(*storeInterval) * time.Second
		args.StoreInterval = &tmp
	}

	if trustedSubnet != nil {
		var err error
		_, args.TrustedSubnet, err = net.ParseCIDR(*trustedSubnet)
		if err != nil {
			log.Printf("failed to parse trusted subnet: %v\n", err)
		}
	}

	return args
}

// parseAgentArgsFromFlags парсит флаги в optionalAgentArgs.
func parseAgentArgsFromFlags() optionalAgentArgs {
	address := new(netAddress)
	flag.Var(address, "a", "address of metrix server in form <host:port>")
	privateKey := flag.String("k", "", "key for hash funcs")
	cryptoPublicKey := flag.String("crypto-key", "", "public key for crypto funcs in HEX")
	pollInterval := flag.Int("p", 2, "poll interval (in sec)")
	reportInterval := flag.Int("r", 2, "report interval (in sec)")
	reportRateLimit := flag.Int("l", 0, "rate limit for simultaneous reports (0 = no limit)")
	isDev := flag.Bool("dev", false, "is metrix agent running in development mode")
	isProfilingEnabled := flag.Bool("profile", false, "is profiling via pprof enabled")

	flag.Parse()

	args := optionalAgentArgs{
		PrivateKey:         privateKey,
		CryptoPublicKey:    cryptoPublicKey,
		PollInterval:       pollInterval,
		ReportRateLimit:    reportRateLimit,
		IsDev:              isDev,
		IsProfilingEnabled: isProfilingEnabled,
	}

	{
		tmp := address.String()
		args.Addr = &tmp
	}

	if reportInterval != nil {
		tmp := time.Duration(*reportInterval) * time.Second
		args.ReportInterval = &tmp
	}

	return args
}
