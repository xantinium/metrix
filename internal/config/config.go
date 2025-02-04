// Пакет config содержит структуры для конфигурации
// агента и сервера.
package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// ServerArgs структура, описывающая аргументы сервера.
type ServerArgs struct {
	Addr string
}

// ParseServerArgs парсит агрументы командной строки в ServerArgs.
func ParseServerArgs() ServerArgs {
	address := new(NetAddress)
	flag.Var(address, "a", "address of metrix server in form <host:port>")

	flag.Parse()

	return ServerArgs{
		Addr: address.String(),
	}
}

// AgentArgs структура, описывающая аргументы агента.
type AgentArgs struct {
	Addr           string
	PollInterval   int
	ReportInterval int
}

// ParseAgentArgs парсит агрументы командной строки в AgentArgs.
func ParseAgentArgs() AgentArgs {
	address := new(NetAddress)
	flag.Var(address, "a", "address of metrix server in form <host:port>")
	pollInterval := flag.Int("p", 2, "poll interval (in sec)")
	reportInterval := flag.Int("r", 2, "report interval (in sec)")

	flag.Parse()

	return AgentArgs{
		Addr:           address.String(),
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
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

	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}

	a.Host = host
	a.Port = port
	a.Parsed = true

	return nil
}
