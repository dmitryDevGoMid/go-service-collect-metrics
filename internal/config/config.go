package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Gzip struct {
	Enable bool `json:"GZIP,omitempty"`
}

type Serializer struct {
	SerType string `json:"SER_TYPE,omitempty"`
}

type Logger struct {
	Encoding string `env:"LOG_ENCODING,omitempty"`
	Level    string `env:"LOG_LEVEL,omitempty"`
}

type Metrics struct {
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

type Server struct {
	Address string `env:"ADDRESS"`
}

type Config struct {
	Server     Server
	Metrics    Metrics
	Logger     Logger
	Serializer Serializer
	Gzip       Gzip
}

var (
	address         string
	reportInterval  int
	pollInterval    int
	logger_encoding string
	logger_level    string
	serialize_type  string
	enable_gzip     bool
)

func init() {
	flag.StringVar(&address, "a", "localhost:8080", "location http server")
	flag.IntVar(&reportInterval, "r", 4, "interval for run metrics")
	flag.IntVar(&pollInterval, "p", 2, "interval for run metrics")

	//	Logger
	flag.StringVar(&logger_encoding, "logen", "full", "set logger config encoding")
	flag.StringVar(&logger_level, "loglv", "InfoLevel", "set logger config level")

	//Serialize Type
	flag.StringVar(&serialize_type, "sertype", "json", "set logger config encoding")

	//Serialize Type
	flag.BoolVar(&enable_gzip, "gzip", true, "set gzip for agent and server")
}

// Разбираем конфигурацию по структурам
func ParseConfig() (*Config, error) {
	flag.Parse()

	var config Config

	config.Metrics.PollInterval = pollInterval
	config.Metrics.ReportInterval = reportInterval

	config.Server.Address = address

	config.Logger.Encoding = logger_encoding
	config.Logger.Level = logger_level

	config.Serializer.SerType = serialize_type

	config.Gzip.Enable = enable_gzip

	//Init by environment variables
	env.Parse(&config.Metrics)
	env.Parse(&config.Server)
	env.Parse(&config.Logger)
	env.Parse(&config.Serializer)
	env.Parse(&config.Gzip)

	return &config, nil
}
