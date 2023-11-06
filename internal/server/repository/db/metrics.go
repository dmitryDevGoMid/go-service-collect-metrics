package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/storage"
	"github.com/jackc/pgx/v5"
)

type MetricsRepository interface {
	repository.MetricsRepository
	GetMetricsTypeIDByName(nameMetric string) (int, error)
}

type metricsRepository struct {
	db *sql.DB
}

// Contruct
func NewMetricsRepository(db *sql.DB) MetricsRepository {
	return &metricsRepository{
		db: db,
	}
}

func (connect *metricsRepository) GetMetricGauge(nameMetric string) (float64, error) {
	var valueMetric float64
	// Query for a value based on a single row.
	if err := connect.db.QueryRow("SELECT value from metrics_gauge where name = $1",
		nameMetric).Scan(&valueMetric); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("canPurchase %s: unknown metrics", nameMetric)
		}
		return 0, fmt.Errorf("canPurchase %s: %v", nameMetric, err)
	}
	return valueMetric, nil
}

func (connect *metricsRepository) UpdateMetricGauge(nameMetric string, value float64) error {
	fmt.Println("DB UpdateMetricGauge")

	metricsTypeID, err := connect.GetMetricsTypeIDByName("gauge")
	if err != nil {
		fmt.Println("error get metrics type by name:", err)
		return fmt.Errorf("error get metrics type by name: %v", err)
	}

	fmt.Println("DB UpdateMetricGauge: ", 1)

	sqlStatement := `UPDATE metrics_gauge SET value = $2 WHERE type_id = $3 AND name = $1;`
	res, err := connect.db.Exec(sqlStatement, nameMetric, value, metricsTypeID)
	if err != nil {
		return fmt.Errorf("error update metrics: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error get rows affected metrics: %v", err)
	}
	if !(count > 0) {
		_, err = connect.db.ExecContext(context.TODO(),
			"INSERT INTO metrics_gauge(value, type_id, name) VALUES (@value, @type_id, @name)",
			pgx.NamedArgs{"value": value, "type_id": metricsTypeID, "name": nameMetric})
	}

	if err != nil {
		return fmt.Errorf("error Update Metrics: %v", err)
	}

	return nil
}

func (connect *metricsRepository) GetMetricCounter(nameMetric string) (int64, error) {
	var deltaMetric int64
	// Query for a value based on a single row.
	if err := connect.db.QueryRow("SELECT delta FROM metrics_counter WHERE name = $1",
		nameMetric).Scan(&deltaMetric); err != nil {
		return 0, err
	}

	return deltaMetric, nil
}

func (connect *metricsRepository) UpdateMetricCounter(nameMetric string, value int64) error {
	fmt.Println("DB UpdateMetricCounter")

	var deltaMetricCalc int64

	metricsTypeID, err := connect.GetMetricsTypeIDByName("counter")

	if err != nil {
		fmt.Println("error get metrics type by name:", err)
		return fmt.Errorf("error get metrics type by name: %v", err)
	}

	deltaMetric, err := connect.GetMetricCounter(nameMetric)

	if err != nil {
		if err == sql.ErrNoRows {
			deltaMetricCalc = value
		} else {
			return fmt.Errorf("select counter metrics: %v", err)
		}
	} else {
		deltaMetricCalc = deltaMetric + value
	}

	fmt.Println("DB UpdateMetricGauge: ", 1)

	sqlStatement := `UPDATE metrics_counter SET delta = $2 WHERE type_id = $3 AND name = $1;`
	res, err := connect.db.Exec(sqlStatement, nameMetric, deltaMetricCalc, metricsTypeID)
	if err != nil {
		return fmt.Errorf("error update metrics: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error get rows affected metrics: %v", err)
	}
	if !(count > 0) {
		_, err = connect.db.ExecContext(context.TODO(),
			"INSERT INTO metrics_counter(delta, type_id, name) VALUES (@delta, @type_id, @name)",
			pgx.NamedArgs{
				"delta":   deltaMetricCalc,
				"type_id": metricsTypeID,
				"name":    nameMetric,
			})
	}

	if err != nil {
		return fmt.Errorf("error Update Metrics: %v", err)
	}

	return nil
}

// Return all Metrics for Page
func (connect *metricsRepository) GetAllMetrics() (*models.MemStorage, error) {
	modelsMemStorage := storage.NewMemStorage()

	rowsMetrics, err := connect.db.QueryContext(context.TODO(), "SELECT value, name FROM metrics_gauge")
	if err != nil {
		return nil, fmt.Errorf("error select all metrics gauge: %v", err)
	}

	defer func() {
		_ = rowsMetrics.Close()
		_ = rowsMetrics.Err()
	}() // or modify return value}

	var valueMetricGauge float64
	var nameMetric string

	for rowsMetrics.Next() {
		err = rowsMetrics.Scan(&valueMetricGauge, &nameMetric)
		if err != nil {
			return nil, err
		}

		modelsMemStorage.Gauge[nameMetric] = valueMetricGauge
	}

	rowsMetrics, err = connect.db.QueryContext(context.TODO(), "SELECT delta, name FROM metrics_counter")
	if err != nil {
		return nil, fmt.Errorf("error select all metrics gauge: %v", err)
	}

	var valueMetricCounter int64

	for rowsMetrics.Next() {
		err = rowsMetrics.Scan(&valueMetricCounter, &nameMetric)
		if err != nil {
			return nil, err
		}

		modelsMemStorage.Counter[nameMetric] = valueMetricCounter
	}

	fmt.Println(modelsMemStorage)

	return modelsMemStorage, nil
}

func (connect *metricsRepository) GetMetricsTypeIDByName(nameMetric string) (int, error) {
	var id int

	fmt.Println("GetMetricsTypeIdByName Metrics => ", nameMetric)

	//rowArray := Row{}
	// Query for a value based on a single row.
	sqlStatement := `SELECT id FROM metrics_type WHERE name = $1`
	//row := db.QueryRow(sqlStatement, id).scan
	if err := connect.db.QueryRow(sqlStatement, nameMetric).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("GetMetricsTypeIdByName => ", id)
			return 0, fmt.Errorf("can purchase %s: unknown metrics", nameMetric)
		}
		fmt.Println("GetMetricsTypeIdByName => ", id)
		return 0, fmt.Errorf("can metrics purchase %s: %v", nameMetric, err)
	}

	fmt.Println("GetMetricsTypeIdByName => ", id)
	return id, nil
}

func (connect *metricsRepository) PingDatabase() error {
	return connect.db.Ping()
}
