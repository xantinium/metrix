// Пакет config содержит структуры для конфигурации
// агента и сервера.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xantinium/metrix/internal/tools"
)

// ServerArgs структура, описывающая аргументы сервера.
type ServerArgs struct {
	Addr  string
	IsDev bool
}

// ParseServerArgs парсит агрументы командной строки в ServerArgs.
func ParseServerArgs() ServerArgs {
	address := new(NetAddress)
	flag.Var(address, "a", "address of metrix server in form <host:port>")
	isDev := flag.Bool("dev", false, "is metrix server running in development mode")

	flag.Parse()

	args := ServerArgs{
		Addr:  address.String(),
		IsDev: *isDev,
	}

	envArgs := parseServerArgsFromEnv()

	if envArgs.Addr != "" {
		args.Addr = envArgs.Addr
	}

	return args
}

// parseServerArgsFromEnv парсит переменные окружения в ServerArgs.
func parseServerArgsFromEnv() ServerArgs {
	return ServerArgs{
		Addr: os.Getenv("ADDRESS"),
	}
}

// AgentArgs структура, описывающая аргументы агента.
type AgentArgs struct {
	Addr           string
	PollInterval   int
	ReportInterval int
	IsDev          bool
}

// ParseAgentArgs парсит агрументы командной строки в AgentArgs.
func ParseAgentArgs() AgentArgs {
	address := new(NetAddress)
	flag.Var(address, "a", "address of metrix server in form <host:port>")
	pollInterval := flag.Int("p", 2, "poll interval (in sec)")
	reportInterval := flag.Int("r", 2, "report interval (in sec)")
	isDev := flag.Bool("dev", false, "is metrix agent running in development mode")

	flag.Parse()

	args := AgentArgs{
		Addr:           address.String(),
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
		IsDev:          *isDev,
	}

	envArgs := parseAgentArgsFromEnv()

	if envArgs.Addr != "" {
		args.Addr = envArgs.Addr
	}
	if envArgs.PollInterval != 0 {
		args.PollInterval = envArgs.PollInterval
	}
	if envArgs.ReportInterval != 0 {
		args.ReportInterval = envArgs.ReportInterval
	}

	return args
}

// parseAgentArgsFromEnv парсит переменные окружения в AgentArgs.
func parseAgentArgsFromEnv() AgentArgs {
	var err error

	args := AgentArgs{
		Addr: os.Getenv("ADDRESS"),
	}

	pollIntervalStr := os.Getenv("POLL_INTERVAL")
	reportIntervalStr := os.Getenv("REPORT_INTERVAL")

	args.PollInterval, err = tools.StrToInt(pollIntervalStr)
	if err != nil {
		args.PollInterval = 0
	}

	args.ReportInterval, err = tools.StrToInt(reportIntervalStr)
	if err != nil {
		args.ReportInterval = 0
	}

	return args
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
