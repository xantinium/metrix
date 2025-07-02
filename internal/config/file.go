package config

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/xantinium/metrix/internal/tools"
)

//easyjson:json
type serverConfig struct {
	Addr               *string `json:"address"`
	RPCAddr            *string `json:"rpc_address"`
	StoragePath        *string `json:"store_file"`
	PrivateKey         *string `json:"key"`
	CryptoPrivateKey   *string `json:"crypto_key"`
	DatabaseConnStr    *string `json:"database_dsn"`
	StoreInterval      *string `json:"store_interval"`
	IsDev              *bool   `json:"is_dev"`
	IsProfilingEnabled *bool   `json:"is_profiling"`
	RestoreStorage     *bool   `json:"restore"`
	EnableRPC          *bool   `json:"enable_rpc"`
	ShutdownTimeout    *string `json:"shutdown_timeout"`
	TrustedSubnet      *string `json:"trusted_subnet"`
}

// parseServerArgsFromFile парсит файл в optionalServerArgs.
func parseServerArgsFromFile() optionalServerArgs {
	result := optionalServerArgs{}

	confPath := getConfigFilePath()
	if confPath != "" {
		confJSON, err := os.ReadFile(confPath)
		if err != nil {
			log.Printf("failed to read conf: %v\n", err)
			return result
		}

		var conf serverConfig
		err = json.Unmarshal(confJSON, &conf)
		if err != nil {
			log.Printf("failed to read conf: %v\n", err)
			return result
		}

		result.Addr = conf.Addr
		result.RPCAddr = conf.RPCAddr
		result.StoragePath = conf.StoragePath
		result.PrivateKey = conf.PrivateKey
		result.CryptoPrivateKey = conf.CryptoPrivateKey
		result.DatabaseConnStr = conf.DatabaseConnStr
		result.IsDev = conf.IsDev
		result.IsProfilingEnabled = conf.IsProfilingEnabled
		result.RestoreStorage = conf.RestoreStorage
		result.EnableRPC = conf.EnableRPC

		if conf.StoreInterval != nil {
			result.StoreInterval = parseDuration(*conf.StoreInterval)
		}

		if conf.ShutdownTimeout != nil {
			result.ShutdownTimeout = parseDuration(*conf.ShutdownTimeout)
		}

		if conf.TrustedSubnet != nil {
			_, result.TrustedSubnet, err = net.ParseCIDR(*conf.TrustedSubnet)
			if err != nil {
				log.Printf("failed to parse trusted subnet: %v\n", err)
				return result
			}
		}
	}

	return result
}

//easyjson:json
type agentConfig struct {
	Addr               *string `json:"address"`
	PrivateKey         *string `json:"key"`
	CryptoPublicKey    *string `json:"crypto_key"`
	RequestMethod      *string `json:"request_method"`
	PollInterval       *string `json:"poll_interval"`
	ReportInterval     *string `json:"report_interval"`
	ReportRateLimit    *int    `json:"report_rate_limit"`
	IsDev              *bool   `json:"is_dev"`
	IsProfilingEnabled *bool   `json:"is_profiling"`
	ShutdownTimeout    *string `json:"shutdown_timeout"`
}

// parseAgentArgsFromFile парсит файл в optionalAgentArgs.
func parseAgentArgsFromFile() optionalAgentArgs {
	result := optionalAgentArgs{}

	confPath := getConfigFilePath()
	if confPath != "" {
		confJSON, err := os.ReadFile(confPath)
		if err != nil {
			log.Printf("failed to read conf: %v\n", err)
			return result
		}

		var conf agentConfig
		err = json.Unmarshal(confJSON, &conf)
		if err != nil {
			log.Printf("failed to read conf: %v\n", err)
			return result
		}

		result.Addr = conf.Addr
		result.PrivateKey = conf.PrivateKey
		result.CryptoPublicKey = conf.CryptoPublicKey
		result.RequestMethod = conf.RequestMethod
		result.ReportRateLimit = conf.ReportRateLimit
		result.IsDev = conf.IsDev
		result.IsProfilingEnabled = conf.IsProfilingEnabled

		if conf.ReportInterval != nil {
			result.ReportInterval = parseDuration(*conf.ReportInterval)
		}

		if conf.PollInterval != nil {
			pollInterval := parseDuration(*conf.PollInterval)
			if pollInterval != nil {
				tmp := int(pollInterval.Seconds())
				result.PollInterval = &tmp
			}
		}

		if conf.ShutdownTimeout != nil {
			result.ShutdownTimeout = parseDuration(*conf.ShutdownTimeout)
		}
	}

	return result
}

func getConfigFilePath() string {
	var confPath string

	flag.StringVar(&confPath, "c", "", "path to config in JSON")
	flag.StringVar(&confPath, "config", "", "path to config in JSON")

	flag.Parse()

	confPathEnv := tools.GetStrFromEnv("CONFIG")
	if confPathEnv.Exists {
		confPath = confPathEnv.Value
	}

	return confPath
}

func parseDuration(rawDuration string) *time.Duration {
	duration, err := time.ParseDuration(rawDuration)
	if err != nil {
		log.Printf("failed to parse duration: %v\n", err)
		return nil
	}

	return &duration
}
