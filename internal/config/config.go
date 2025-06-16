// Package config содержит структуры для конфигурации
// агента и сервера.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/xantinium/metrix/internal/tools"
)

// ServerArgs структура, описывающая аргументы сервера.
type ServerArgs struct {
	Addr               string
	StoragePath        string
	PrivateKey         string
	CryptoPrivateKey   string
	DatabaseConnStr    string
	StoreInterval      time.Duration
	IsDev              bool
	IsProfilingEnabled bool
	RestoreStorage     bool
}

// ParseServerArgs парсит агрументы командной строки в ServerArgs.
func ParseServerArgs() ServerArgs {
	return mergeServerArgs(
		parseServerArgsFromFile(),
		parseServerArgsFromFlags(),
		parseServerArgsFromEnv(),
	)
}

// AgentArgs структура, описывающая аргументы агента.
type AgentArgs struct {
	Addr               string
	PrivateKey         string
	CryptoPublicKey    string
	PollInterval       int
	ReportInterval     time.Duration
	ReportRateLimit    int
	IsDev              bool
	IsProfilingEnabled bool
}

// ParseAgentArgs парсит агрументы командной строки в AgentArgs.
func ParseAgentArgs() AgentArgs {
	return mergeAgentArgs(
		parseAgentArgsFromFile(),
		parseAgentArgsFromFlags(),
		parseAgentArgsFromEnv(),
	)
}

// netAddress кастомная структура для обработки флага -a.
type netAddress struct {
	Host   string
	Port   int
	Parsed bool
}

// String возращает сериализованную строку.
func (addr netAddress) String() string {
	if addr.Parsed {
		return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
	}

	return "localhost:8080"
}

// Set парсит структуру из сырой строки.
func (addr *netAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("invalid address format")
	}

	host := hp[0]

	port, err := tools.StrToInt(hp[1])
	if err != nil {
		return err
	}

	addr.Host = host
	addr.Port = port
	addr.Parsed = true

	return nil
}

type optionalServerArgs struct {
	Addr               *string
	StoragePath        *string
	PrivateKey         *string
	CryptoPrivateKey   *string
	DatabaseConnStr    *string
	StoreInterval      *time.Duration
	IsDev              *bool
	IsProfilingEnabled *bool
	RestoreStorage     *bool
}

func mergeServerArgs(argsToMerge ...optionalServerArgs) ServerArgs {
	result := ServerArgs{}

	for _, args := range argsToMerge {
		if args.Addr != nil {
			result.Addr = *args.Addr
		}
		if args.StoragePath != nil && fs.ValidPath(*args.StoragePath) {
			result.StoragePath = *args.StoragePath
		}
		if args.PrivateKey != nil {
			result.PrivateKey = *args.PrivateKey
		}
		if args.CryptoPrivateKey != nil {
			result.CryptoPrivateKey = *args.CryptoPrivateKey
		}
		if args.DatabaseConnStr != nil {
			result.DatabaseConnStr = *args.DatabaseConnStr
		}
		if args.StoreInterval != nil && *args.StoreInterval >= 0 {
			result.StoreInterval = *args.StoreInterval
		}
		if args.IsDev != nil {
			result.IsDev = *args.IsDev
		}
		if args.IsProfilingEnabled != nil {
			result.IsProfilingEnabled = *args.IsProfilingEnabled
		}
		if args.RestoreStorage != nil {
			result.RestoreStorage = *args.RestoreStorage
		}
	}

	return result
}

type optionalAgentArgs struct {
	Addr               *string
	PrivateKey         *string
	CryptoPublicKey    *string
	PollInterval       *int
	ReportInterval     *time.Duration
	ReportRateLimit    *int
	IsDev              *bool
	IsProfilingEnabled *bool
}

func mergeAgentArgs(argsToMerge ...optionalAgentArgs) AgentArgs {
	result := AgentArgs{}

	for _, args := range argsToMerge {
		if args.Addr != nil {
			result.Addr = *args.Addr
		}
		if args.PrivateKey != nil {
			result.PrivateKey = *args.PrivateKey
		}
		if args.CryptoPublicKey != nil {
			result.CryptoPublicKey = *args.CryptoPublicKey
		}
		if args.PollInterval != nil && *args.PollInterval > 0 {
			result.PollInterval = *args.PollInterval
		}
		if args.ReportInterval != nil && *args.ReportInterval > 0 {
			result.ReportInterval = *args.ReportInterval
		}
		if args.ReportRateLimit != nil && *args.ReportRateLimit >= 0 {
			result.ReportRateLimit = *args.ReportRateLimit
		}
		if args.IsDev != nil {
			result.IsDev = *args.IsDev
		}
		if args.IsProfilingEnabled != nil {
			result.IsProfilingEnabled = *args.IsProfilingEnabled
		}
	}

	return result
}
