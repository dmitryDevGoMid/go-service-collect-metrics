package db

import (
	"fmt"
	"sync"

	"database/sql"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var connectPostgres *sql.DB

// Used during creation of singleton client object
var err error

// Used to execute client creation procedure only once.
var postgresOnce sync.Once

// Интерфейс которыq будет реализован в структуре conn
type Connection interface {
	Close() error
	DB() *sql.DB
	Ping() error
}

// Структура которая будет возвращать
type conn struct {
	connection *sql.DB
	cfg        *config.Config
}

// Выполняем соединение с базой данных и возвращаем коннект вместе с ошибкой
func NewConnection(cfg *config.Config) Connection {

	connectPostgres, err := GetPostgresConnection(cfg)
	if err != nil {
		fmt.Println("Error opening database: ", err)
	}

	return &conn{connection: connectPostgres, cfg: cfg}
}

// func CloseClientDB() {
func (c *conn) Close() error {
	c.connection.Close()
	if err != nil {
		return err
	}

	return nil
}

// Реализуем функции на структуре conn для того чтобы она соответствовала интерфейсу Connection
func (c *conn) DB() *sql.DB {
	return c.connection
}

func (c *conn) Ping() error {
	err = c.connection.Ping()
	if err != nil {
		return err
	}

	return nil
}

// Получаем одно соединение для базы данных
func GetPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	postgresOnce.Do(func() {

		connectPostgres, err = sql.Open("pgx", cfg.DataBase.DatabaseURL)
		if err != nil {
			fmt.Println("Error opening database: ", err)
		}

		connectPostgres.Exec(`set search_path='public'`)

		//defer db.Close()

	})

	return connectPostgres, err
}
