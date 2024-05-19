package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
)

type ConfigJsonStruct struct {
	Address        string `json:"address,omitempty"`
	Restore        bool   `json:"restore,omitempty"`
	Store_interval int    `json:"store_interval,omitempty"`
	Store_file     string `json:"store_file,omitempty"`
	Database_dsn   string `json:"database_dsn,omitempty"`
	Crypto_key     string `json:"crypto_key,omitempty"`
}

type ConfigJson struct {
	ConfigJson string `env:"CONFIG"`
}

type PathEncrypt struct {
	PathEncryptKey   string `env:"CRYPTO_KEY"`
	KeyEncryptEnbled bool
}

type HashSHA256 struct {
	Key string `env:"KEY"`
}

type DataBase struct {
	DatabaseURL string `env:"DATABASE_DSN"`
}

type File struct {
	Restore         bool   `env:"RESTORE"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
}

type Gzip struct {
	Enable bool `env:"GZIP"`
}

type Serializer struct {
	SerType string `env:"SER_TYPE"`
}

type Logger struct {
	Encoding string `env:"LOG_ENCODING"`
	Level    string `env:"LOG_LEVEL"`
}

type Metrics struct {
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
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
	File        File
	DataBase    DataBase
	HashSHA256  HashSHA256
	PathEncrypt PathEncrypt
	ConfigJson  ConfigJson
}

var (
	address string
	//reportInterval int
	pollInterval   int
	loggerEncoding string
	loggerLevel    string
	serializeType  string
	enableGzip     bool

	restoreFile       bool
	storeIntervalFile int
	fileStoragePath   string

	databaseURL string

	keySHA256 string

	pathEncryptKey string

	configJson string
)

var result = ConfigJsonStruct{
	Address:        "localhost:8080",
	Crypto_key:     "",
	Database_dsn:   "postgres://manager:M45fgMetr@localhost:5432/metrics?sslmode=disable",
	Restore:        true,
	Store_interval: 300,
	Store_file:     "/tmp/metrics-db.json",
}

/*

{
    "address": "localhost:8080", // аналог переменной окружения ADDRESS или флага -a
    "restore": true, // аналог переменной окружения RESTORE или флага -r
    "store_interval": "1s", // аналог переменной окружения STORE_INTERVAL или флага -i
    "store_file": "/path/to/file.db", // аналог переменной окружения STORE_FILE или -f
    "database_dsn": "", // аналог переменной окружения DATABASE_DSN или флага -d
    "crypto_key": "/path/to/key.pem" // аналог переменной окружения CRYPTO_KEY или флага -crypto-key
}

*/

func InitFlag(flagInit *flag.FlagSet) {
	//Config
	flagInit.StringVar(&configJson, "c", "", "path to config file by json")
	flagInit.StringVar(&configJson, "config", "", "path to config file by json")

	//Encrypt
	flagInit.StringVar(&pathEncryptKey, "crypto-key", result.Crypto_key, "path encrypt key")

	flagInit.StringVar(&address, "a", result.Address, "location http server")
	//flag.IntVar(&reportInterval, "r", 400, "interval for run metrics")
	flagInit.IntVar(&pollInterval, "p", 200, "interval for run metrics")

	//	Logger
	flagInit.StringVar(&loggerEncoding, "logen", "console", "set logger config encoding")
	flagInit.StringVar(&loggerLevel, "loglv", "InfoLevel", "set logger config level")

	//Serialize Type
	flagInit.StringVar(&serializeType, "sertype", "json", "set logger config encoding")

	//Serialize Type
	flagInit.BoolVar(&enableGzip, "gzip", false, "set gzip for agent and server")

	//File
	flagInit.BoolVar(&restoreFile, "r", result.Restore, "restore file")
	//flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "path file")
	flagInit.StringVar(&fileStoragePath, "f", result.Store_file, "path file")
	flagInit.IntVar(&storeIntervalFile, "i", result.Store_interval, "store interval file")

	//Connection Database
	/*
	  - POSTGRES_PASSWORD=M45fgMetr
	  - POSTGRES_USER=manager
	  - POSTGRES_DB=metrics
	*/
	flagInit.StringVar(&databaseURL, "d", result.Database_dsn, "database url for conection postgress")

	//sha 256 key
	flagInit.StringVar(&keySHA256, "k", "invalidkey", "set key for calc SHA256")
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

	}

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

// Разбираем конфигурацию по структурам
func ParseConfig() (*Config, error) {
	ParseFlag()

	var config Config

	config.Metrics.PollInterval = pollInterval
	//config.Metrics.ReportInterval = reportInterval

	config.Server.Address = address

	config.Logger.Encoding = loggerEncoding
	config.Logger.Level = loggerLevel

	config.Serializer.SerType = serializeType

	config.Gzip.Enable = enableGzip

	config.File.FileStoragePath = fileStoragePath
	config.File.Restore = restoreFile
	config.File.StoreInterval = storeIntervalFile

	config.DataBase.DatabaseURL = databaseURL

	config.HashSHA256.Key = keySHA256

	config.PathEncrypt.PathEncryptKey = pathEncryptKey
	config.PathEncrypt.KeyEncryptEnbled = false

	//Init by environment variables
	env.Parse(&config.Metrics)
	env.Parse(&config.Server)
	env.Parse(&config.Logger)
	env.Parse(&config.Serializer)
	env.Parse(&config.Gzip)
	env.Parse(&config.File)
	env.Parse(&config.DataBase)
	env.Parse(&config.HashSHA256)
	env.Parse(&config.PathEncrypt)

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
