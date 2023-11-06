package migration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var tableList = [...]string{"metrics_type", "metrics_gauge", "metrics_counter"}
var tableListDrop = [...]string{"metrics_gauge", "metrics_counter", "metrics_type"}

type Migration interface {
	Run(ctx context.Context)
}

type migration struct {
	db *sql.DB
}

func NewMigration(db *sql.DB) *migration {
	return &migration{db: db}
}

func (m *migration) RunCreate(ctx context.Context) {
	for _, tableName := range tableList {
		if !m.existsTable(ctx, tableName) {
			m.CreateTable(tableName)
		}
	}
}

func (m *migration) RunDrop(ctx context.Context) {
	for _, tableName := range tableListDrop {
		if m.existsTable(ctx, tableName) {
			m.DropTable(tableName)
		}
	}
}

func (m *migration) existsTable(ctx context.Context, tableName string) bool {
	var n int64
	err := m.db.QueryRow("select 1 from information_schema.tables where table_name = $1", tableName).Scan(&n)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false
		}
		return false
	}

	return true
}

/*
Example
`
	CREATE TABLE IF NOT EXISTS TEST(
			UserID int
	)
`
*/

func (m *migration) DropTable(tabelName string) {
	tableCodeDrop := ""

	switch val := tabelName; val {
	case "metrics_type":
		tableCodeDrop = m.MetricsTypeDrop()
	case "metrics_gauge":
		tableCodeDrop = m.MetricsGaugeDrop()
	case "metrics_counter":
		tableCodeDrop = m.MetricsCounterDrop()
	default:
		fmt.Println("Not condition case: ", val)
	}

	if tableCodeDrop != "" {
		db, err := m.db.Exec(tableCodeDrop)
		if err != nil {
			fmt.Printf("Error drop table: %v\n", err)
		}

		fmt.Printf("Value of db %v\n", db)
	}
}

func (m *migration) MetricsTypeDrop() string {
	return `drop table metrics_type`
}

func (m *migration) MetricsGaugeDrop() string {
	return `drop table metrics_gauge`
}

func (m *migration) MetricsCounterDrop() string {
	return `drop table metrics_counter`
}

func (m *migration) CreateTable(tabelName string) {

	tableCodeCreate := ""

	switch val := tabelName; val {
	case "metrics_type":
		tableCodeCreate = m.MetricsType()
	case "metrics_gauge":
		tableCodeCreate = m.MetricsGauge()
	case "metrics_counter":
		tableCodeCreate = m.MetricsCounter()
	default:
		fmt.Println("Not condition case: ", val)
	}

	if tableCodeCreate != "" {
		db, err := m.db.Exec(tableCodeCreate)

		if err != nil {
			fmt.Printf("Error create table: %v\n", err)
		}
		fmt.Printf("Value of db %v\n", db)
	}
}

func (m *migration) MetricsType() string {
	return `
	CREATE TABLE IF NOT EXISTS metrics_type(
		id INT GENERATED ALWAYS AS IDENTITY,
		name VARCHAR(255) NOT NULL,
		PRIMARY KEY(id) 
	);
	INSERT INTO metrics_type(name) VALUES('gauge');
	INSERT INTO metrics_type(name) VALUES('counter');
	`
}

func (m *migration) MetricsGauge() string {
	return `
	CREATE TABLE IF NOT EXISTS metrics_gauge(
		id INT GENERATED ALWAYS AS IDENTITY,
		type_id INT,
		value double precision,
		name varchar(255),
		PRIMARY KEY(id),
		CONSTRAINT fk_metrics_type
			FOREIGN KEY(type_id) 
				REFERENCES metrics_type(id)
	);
	`
}

func (m *migration) MetricsCounter() string {
	return `
	CREATE TABLE IF NOT EXISTS metrics_counter(
		id INT GENERATED ALWAYS AS IDENTITY,
		type_id INT,
		delta bigint,
		name varchar(255),
		PRIMARY KEY(id),
		CONSTRAINT fk_type_metrics
			FOREIGN KEY(type_id) 
				REFERENCES metrics_type(id)
	);
	`
}

// TODO table store for batch by date and time
func (m *migration) MetricsBatch() string {
	return `
	CREATE TABLE IF NOT EXISTS metrics_batch(
		id INT GENERATED ALWAYS AS IDENTITY,
		time_date_create TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		time_date_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY(id),
	);
	`
}
