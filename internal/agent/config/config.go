package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
)

// Структура которая необходима для конфигурационных данных, читаемых из файл
type ConfigJSONStruct struct {
	Address        string `json:"address,omitempty"`
	ReportInterval int    `json:"report_interval,omitempty"`
	PollInterval   int    `json:"poll_interval,omitempty"`
	CryptoKey      string `json:"crypto_key,omitempty"`
}

type ConfigJSON struct {
	ConfigJson string `env:"CONFIG"`
}

type PathEncrypt struct {
	PathEncryptKey   string `env:"CRYPTO_KEY"`
	KeyEncryptEnbled bool
}

type Workers struct {
	LimitWorkers int `env:"RATE_LIMIT"`
}

type SHA256 struct {
	Key string `env:"KEY"`
}

type Gzip struct {
	Enable bool `env:"GZIP,omitempty"`
}

type Serializer struct {
	SerType string `env:"SER_TYPE,omitempty"`
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
	Server      Server
	Metrics     Metrics
	Logger      Logger
	Serializer  Serializer
	Gzip        Gzip
	SHA256      SHA256
	Workers     Workers
	PathEncrypt PathEncrypt
	ConfigJSON  ConfigJSON
}

var (
	address          string
	reportInterval   int
	pollInterval     int
	loggerEncoding   string
	loggerLevel      string
	serializeType    string
	enableGzip       bool
	sendMeticsBatch  bool
	keySHA256        string
	limitWorkersPool int
	pathEncryptKey   string
	configJson       string
)

/*
{
    "address": "localhost:8080", // аналог переменной окружения ADDRESS или флага -a
    "report_interval": "1s", // аналог переменной окружения REPORT_INTERVAL или флага -r
    "poll_interval": "1s", // аналог переменной окружения POLL_INTERVAL или флага -p
    "crypto_key": "/path/to/key.pem" // аналог переменной окружения CRYPTO_KEY или флага -crypto-key
}
*/

var result = ConfigJSONStruct{
	Address:        "localhost:8080",
	CryptoKey:      "",
	PollInterval:   2,
	ReportInterval: 10,
}

func InitFlag(flagInit *flag.FlagSet) {
	//Config
	flagInit.StringVar(&configJson, "c", "", "path to config file by json")
	flagInit.StringVar(&configJson, "config", "", "path to config file by json")

	//Encrypt
	flagInit.StringVar(&pathEncryptKey, "crypto-key", result.CryptoKey, "path encrypt key")

	flagInit.StringVar(&address, "a", result.Address, "location http server")
	flagInit.IntVar(&reportInterval, "r", result.ReportInterval, "interval for run metrics")
	flagInit.IntVar(&pollInterval, "p", result.PollInterval, "interval for run metrics")
	flagInit.BoolVar(&sendMeticsBatch, "mb", true, "set gzip for agent and server")

	//	Logger
	flagInit.StringVar(&loggerEncoding, "logen", "full", "set logger config encoding")
	flagInit.StringVar(&loggerLevel, "loglv", "InfoLevel", "set logger config level")

	//Serialize Type
	flagInit.StringVar(&serializeType, "sertype", "json", "set logger config encoding")

	//Serialize Type
	flagInit.BoolVar(&enableGzip, "gzip", false, "set gzip for agent and server")

	//sha 256 key
	flagInit.StringVar(&keySHA256, "k", "invalidkey", "set gzip for agent and server")

	//Works
	flagInit.IntVar(&limitWorkersPool, "l", 5, "limit workers send to server metrics")
}

// Это просто треш, а не библиотека, такой процесс повторной инициализации флагов могли придумать только в golang - жесть
// Такое выполнение повтороне позволяетс реализовать инициализацию флаго: по умолчанию, из командной строки и из файла json,
// при условии, что у конфигов из файла приоритет самый низкий, тоесть инициализация значения происходит только из файла в том
// случаи если нет значение в командной строке и в переменных окружения.
func ParseFlag() {
	flags1 := flag.NewFlagSet("myapp1", flag.ExitOnError)
	InitFlag(flags1)
	InitFlag(flag.CommandLine)
	flag.Parse()

	err := flags1.Parse(nil)
	if err != nil {
		fmt.Println(err)
	}

	//Проверяем наличие конфигурационного файла, если он есть то выставляем значения по умолчанию из него
	if configJson != "" {
		errConfig := ConfigFileRead(configJson)
		//Если ошибка то возвращаемся из функции, все флаги проинициализированы
		if errConfig != nil {
			fmt.Println("Error Set Path File Config", errConfig)
			return
		}

		flags1.VisitAll(func(f *flag.Flag) {
			if f.Name == "crypto-key" {

				fmt.Printf("Flag name: %s\n", f.Name)
				fmt.Printf("Flag default value: %v\n", f.DefValue)
				fmt.Printf("Flag usage: %s\n", f.Usage)
				fmt.Println("Flag Value:", f.Value)

				fmt.Println()
			}
		})

		flags2 := flag.NewFlagSet("myapp2", flag.ExitOnError)
		InitFlag(flags2)
		flags2.Parse(nil)

		flag.Parse()

		flags2.VisitAll(func(f *flag.Flag) {
			if f.Name == "crypto-key" {

				fmt.Printf("Flag2 name: %s\n", f.Name)
				fmt.Printf("Flag2 default value: %v\n", f.DefValue)
				fmt.Printf("Flag2 usage: %s\n", f.Usage)
				fmt.Println("Flag2 Value:", f.Value)

				fmt.Println()
			}
		})

		fmt.Println("pathEncryptKey=>", pathEncryptKey)
	}
	//os.Exit(1)
}

// Разбираем конфигурацию по структурам
func ParseConfig() (*Config, error) {
	ParseFlag()

	var config Config

	config.Metrics.PollInterval = pollInterval
	config.Metrics.ReportInterval = reportInterval
	config.Metrics.SendMeticsBatch = sendMeticsBatch

	config.Server.Address = address

	config.Logger.Encoding = loggerEncoding
	config.Logger.Level = loggerLevel

	config.Serializer.SerType = serializeType

	config.Gzip.Enable = enableGzip

	config.SHA256.Key = keySHA256

	config.Workers.LimitWorkers = limitWorkersPool

	config.PathEncrypt.PathEncryptKey = pathEncryptKey
	config.PathEncrypt.KeyEncryptEnbled = false

	//Init by environment variables
	env.Parse(&config.Metrics)
	env.Parse(&config.Server)
	env.Parse(&config.Logger)
	env.Parse(&config.Serializer)
	env.Parse(&config.Gzip)
	env.Parse(&config.SHA256)
	env.Parse(&config.Workers)
	env.Parse(&config.PathEncrypt)
	env.Parse(&config.ConfigJSON)

	return &config, nil
}

func ConfigFileRead(path string) error {
	// Получаем текущую рабочую директорию
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении текущей рабочей директории:", err)
		return err
	}

	// Путь к файлу
	filePath := filepath.Join(wd, path)

	// Проверяем существование файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Файл не существует")
		return err
	}

	// Читаем содержимое файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return err
	}

	// Выполняем маршалинг данных в структуру
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("Ошибка при маршалинге данных:", err)
		return err
	}

	return nil
}
