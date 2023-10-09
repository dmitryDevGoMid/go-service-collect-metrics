package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Metrics struct {
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

type Server struct {
	Address string `env:"ADDRESS"`
}

type Config struct {
	Server  Server
	Metrics Metrics
}

var (
	address        string
	reportInterval int
	pollInterval   int
)

func init() {
	flag.StringVar(&address, "a", "localhost:8080", "location http server")
	flag.IntVar(&reportInterval, "r", 10, "interval for run metrics")
	flag.IntVar(&pollInterval, "p", 2, "interval for run metrics")
}

// Разбираем конфигурацию по структурам
func ParseConfig() (*Config, error) {
	flag.Parse()

	var config Config

	config.Metrics.PollInterval = pollInterval
	config.Metrics.ReportInterval = reportInterval

	config.Server.Address = address

	env.Parse(&config.Metrics)
	env.Parse(&config.Server)

	return &config, nil
}
