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
	ReportInterval  int  `env:"REPORT_INTERVAL"`
	PollInterval    int  `env:"POLL_INTERVAL"`
	SendMeticsBatch bool `env:"SEND_METRICS_BATCH"`
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
	loggerEncoding  string
	loggerLevel     string
	serializeType   string
	enableGzip      bool
	sendMeticsBatch bool
)

func init() {
	flag.StringVar(&address, "a", "localhost:8080", "location http server")
	flag.IntVar(&reportInterval, "r", 4, "interval for run metrics")
	flag.IntVar(&pollInterval, "p", 2, "interval for run metrics")
	flag.BoolVar(&sendMeticsBatch, "mb", true, "set gzip for agent and server")

	//	Logger
	flag.StringVar(&loggerEncoding, "logen", "full", "set logger config encoding")
	flag.StringVar(&loggerLevel, "loglv", "InfoLevel", "set logger config level")

	//Serialize Type
	flag.StringVar(&serializeType, "sertype", "json", "set logger config encoding")

	//Serialize Type
	flag.BoolVar(&enableGzip, "gzip", false, "set gzip for agent and server")
}

// Разбираем конфигурацию по структурам
func ParseConfig() (*Config, error) {
	flag.Parse()

	var config Config

	config.Metrics.PollInterval = pollInterval
	config.Metrics.ReportInterval = reportInterval
	config.Metrics.SendMeticsBatch = sendMeticsBatch

	config.Server.Address = address

	config.Logger.Encoding = loggerEncoding
	config.Logger.Level = loggerLevel

	config.Serializer.SerType = serializeType

	config.Gzip.Enable = enableGzip

	//Init by environment variables
	env.Parse(&config.Metrics)
	env.Parse(&config.Server)
	env.Parse(&config.Logger)
	env.Parse(&config.Serializer)
	env.Parse(&config.Gzip)

	return &config, nil
}
