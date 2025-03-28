// Пакет config содержит структуры для конфигурации
// агента и сервера.
package config

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/xantinium/metrix/internal/tools"
)

// ServerArgs структура, описывающая аргументы сервера.
type ServerArgs struct {
	Addr            string
	IsDev           bool
	PrivateKey      string
	StoreInterval   time.Duration
	StoragePath     string
	RestoreStorage  bool
	DatabaseConnStr string
}

// ParseServerArgs парсит агрументы командной строки в ServerArgs.
func ParseServerArgs() ServerArgs {
	address := new(NetAddress)
	flag.Var(address, "a", "address of metrix server in form <host:port>")
	isDev := flag.Bool("dev", false, "is metrix server running in development mode")
	privateKey := flag.String("k", "", "key for hash funcs")
	storeInterval := flag.Int("i", 300, "interval (in seconds) of writing metrics into file")
	storagePath := flag.String("f", "./metrix.db", "path to file for metrics writing")
	restoreStorage := flag.Bool("r", true, "read metrics from file on start")
	databaseConnStr := flag.String("d", "", "connection string for postgresql")

	flag.Parse()

	args := ServerArgs{
		Addr:            address.String(),
		IsDev:           *isDev,
		PrivateKey:      *privateKey,
		StoragePath:     *storagePath,
		RestoreStorage:  *restoreStorage,
		DatabaseConnStr: *databaseConnStr,
	}
	if storeInterval != nil {
		args.StoreInterval = time.Duration(*storeInterval) * time.Second
	}

	envArgs := parseServerArgsFromEnv()

	if envArgs.Addr.Exists {
		args.Addr = envArgs.Addr.Value
	}
	if envArgs.PrivateKey.Exists {
		args.PrivateKey = envArgs.PrivateKey.Value
	}
	if envArgs.StoreInterval.Exists && envArgs.StoreInterval.Value >= 0 {
		args.StoreInterval = time.Duration(envArgs.StoreInterval.Value) * time.Second
	}
	if envArgs.StoragePath.Exists && fs.ValidPath(envArgs.StoragePath.Value) {
		args.StoragePath = envArgs.StoragePath.Value
	}
	if envArgs.RestoreStorage.Exists {
		args.RestoreStorage = envArgs.RestoreStorage.Value
	}
	if envArgs.DatabaseConnStr.Exists {
		args.DatabaseConnStr = envArgs.DatabaseConnStr.Value
	}

	return args
}

type serverEnvArgs struct {
	Addr            tools.StrEnvVar
	PrivateKey      tools.StrEnvVar
	StoreInterval   tools.IntEnvVar
	StoragePath     tools.StrEnvVar
	RestoreStorage  tools.BoolEnvVar
	DatabaseConnStr tools.StrEnvVar
}

// parseServerArgsFromEnv парсит переменные окружения в serverEnvArgs.
func parseServerArgsFromEnv() serverEnvArgs {
	return serverEnvArgs{
		Addr:            tools.GetStrFromEnv("ADDRESS"),
		PrivateKey:      tools.GetStrFromEnv("KEY"),
		StoreInterval:   tools.GetIntFromEnv("STORE_INTERVAL"),
		StoragePath:     tools.GetStrFromEnv("FILE_STORAGE_PATH"),
		RestoreStorage:  tools.GetBoolFromEnv("RESTORE"),
		DatabaseConnStr: tools.GetStrFromEnv("DATABASE_DSN"),
	}
}

// AgentArgs структура, описывающая аргументы агента.
type AgentArgs struct {
	Addr            string
	PrivateKey      string
	PollInterval    int
	ReportInterval  time.Duration
	ReportRateLimit int
	IsDev           bool
}

// ParseAgentArgs парсит агрументы командной строки в AgentArgs.
func ParseAgentArgs() AgentArgs {
	address := new(NetAddress)
	flag.Var(address, "a", "address of metrix server in form <host:port>")
	privateKey := flag.String("k", "", "key for hash funcs")
	pollInterval := flag.Int("p", 2, "poll interval (in sec)")
	reportInterval := flag.Int("r", 2, "report interval (in sec)")
	reportRateLimit := flag.Int("l", 0, "rate limit for simultaneous reports (0 = no limit)")
	isDev := flag.Bool("dev", false, "is metrix agent running in development mode")

	flag.Parse()

	args := AgentArgs{
		Addr:            address.String(),
		PrivateKey:      *privateKey,
		PollInterval:    *pollInterval,
		ReportRateLimit: *reportRateLimit,
		IsDev:           *isDev,
	}
	if reportInterval != nil {
		args.ReportInterval = time.Duration(*reportInterval) * time.Second
	}

	envArgs := parseAgentArgsFromEnv()

	if envArgs.Addr.Exists {
		args.Addr = envArgs.Addr.Value
	}
	if envArgs.PrivateKey.Exists {
		args.PrivateKey = envArgs.PrivateKey.Value
	}
	if envArgs.PollInterval.Exists && envArgs.PollInterval.Value > 0 {
		args.PollInterval = envArgs.PollInterval.Value
	}
	if envArgs.ReportInterval.Exists && envArgs.ReportInterval.Value > 0 {
		args.ReportInterval = time.Duration(envArgs.ReportInterval.Value) * time.Second
	}
	if envArgs.ReportRateLimit.Exists && envArgs.ReportRateLimit.Value >= 0 {
		args.ReportRateLimit = envArgs.ReportRateLimit.Value
	}

	return args
}

type agentEnvArgs struct {
	Addr            tools.StrEnvVar
	PrivateKey      tools.StrEnvVar
	PollInterval    tools.IntEnvVar
	ReportInterval  tools.IntEnvVar
	ReportRateLimit tools.IntEnvVar
}

// parseAgentArgsFromEnv парсит переменные окружения в agentEnvArgs.
func parseAgentArgsFromEnv() agentEnvArgs {
	return agentEnvArgs{
		Addr:            tools.GetStrFromEnv("ADDRESS"),
		PrivateKey:      tools.GetStrFromEnv("KEY"),
		PollInterval:    tools.GetIntFromEnv("POLL_INTERVAL"),
		ReportInterval:  tools.GetIntFromEnv("REPORT_INTERVAL"),
		ReportRateLimit: tools.GetIntFromEnv("RATE_LIMIT"),
	}
}

// NetAddress кастомная структура для обработки флага -a.
type NetAddress struct {
	Host   string
	Port   int
	Parsed bool
}

// String возращает сериализованную строку.
func (a NetAddress) String() string {
	if a.Parsed {
		return fmt.Sprintf("%s:%d", a.Host, a.Port)
	}

	return "localhost:8080"
}

// Set парсит структуру из сырой строки.
func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("invalid address format")
	}

	host := hp[0]

	port, err := tools.StrToInt(hp[1])
	if err != nil {
		return err
	}

	a.Host = host
	a.Port = port
	a.Parsed = true

	return nil
}
