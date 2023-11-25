package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/unserialize"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/storage"
	"github.com/jackc/pgx/v5"
)

type Decoretor repository.Decorator

type MetricsRepository interface {
	repository.MetricsRepository
	GetMetricsTypeIDByName(ctx context.Context, nameMetric string) (int, error)
}

type metricsRepository struct {
	db    *sql.DB
	mutex *sync.Mutex
}

// Contruct
func NewMetricsRepository(db *sql.DB) MetricsRepository {
	var mutex sync.Mutex
	return &metricsRepository{
		db:    db,
		mutex: &mutex,
	}
}

func (connect *metricsRepository) GetMetricGauge(ctx context.Context, nameMetric string) (float64, error) {
	var valueMetric float64
	// Query for a value based on a single row.
	if err := connect.db.QueryRowContext(ctx, "SELECT value from metrics_gauge where name = $1",
		nameMetric).Scan(&valueMetric); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("canPurchase %s: unknown metrics", nameMetric)
		}
		return 0, fmt.Errorf("canPurchase %s: %v", nameMetric, err)
	}
	return valueMetric, nil
}

func (connect *metricsRepository) UpdateMetricGauge(ctx context.Context, nameMetric string, value float64) error {

	connect.mutex.Lock()

	defer connect.mutex.Unlock()

	fmt.Println("DB UpdateMetricGauge")

	metricsTypeID, err := connect.GetMetricsTypeIDByName(ctx, "gauge")
	if err != nil {
		fmt.Println("error get metrics type by name:", err)
		return fmt.Errorf("error get metrics type by name: %v", err)
	}

	fmt.Println("DB UpdateMetricGauge: ", 1)

	sqlStatement := `UPDATE metrics_gauge SET value = $2 WHERE type_id = $3 AND name = $1;`
	res, err := connect.db.ExecContext(ctx, sqlStatement, nameMetric, value, metricsTypeID)
	if err != nil {
		return fmt.Errorf("error update metrics: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error get rows affected metrics: %v", err)
	}
	if !(count > 0) {
		_, err = connect.db.ExecContext(ctx,
			"INSERT INTO metrics_gauge(value, type_id, name) VALUES (@value, @type_id, @name)",
			pgx.NamedArgs{"value": value, "type_id": metricsTypeID, "name": nameMetric})
	}

	if err != nil {
		return fmt.Errorf("error Update Metrics: %v", err)
	}

	return nil
}

func (connect *metricsRepository) GetMetricCounter(ctx context.Context, nameMetric string) (int64, error) {
	var deltaMetric int64
	// Query for a value based on a single row.
	if err := connect.db.QueryRow("SELECT delta FROM metrics_counter WHERE name = $1",
		nameMetric).Scan(&deltaMetric); err != nil {
		return 0, err
	}

	return deltaMetric, nil
}

func (connect *metricsRepository) UpdateMetricCounter(ctx context.Context, nameMetric string, value int64) error {

	connect.mutex.Lock()

	defer connect.mutex.Unlock()

	fmt.Println("DB UpdateMetricCounter")

	var deltaMetricCalc int64

	metricsTypeID, err := connect.GetMetricsTypeIDByName(ctx, "counter")

	if err != nil {
		fmt.Println("error get metrics type by name:", err)
		return fmt.Errorf("error get metrics type by name: %v", err)
	}

	deltaMetric, err := connect.GetMetricCounter(ctx, nameMetric)

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
	res, err := connect.db.ExecContext(ctx, sqlStatement, nameMetric, deltaMetricCalc, metricsTypeID)
	if err != nil {
		return fmt.Errorf("error update metrics: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error get rows affected metrics: %v", err)
	}
	if !(count > 0) {
		_, err = connect.db.ExecContext(ctx,
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
func (connect *metricsRepository) GetAllMetrics(ctx context.Context) (*models.MemStorage, error) {
	modelsMemStorage := storage.NewMemStorage()

	rowsMetrics, err := connect.db.QueryContext(ctx, "SELECT value, name FROM metrics_gauge")
	if err != nil {
		return nil, fmt.Errorf("error select all metrics gauge: %v", err)
	}

	// Закрываем rowsMetrics
	defer func() {
		_ = rowsMetrics.Close()
		_ = rowsMetrics.Err()
	}()

	var valueMetricGauge float64
	var nameMetric string

	for rowsMetrics.Next() {
		err = rowsMetrics.Scan(&valueMetricGauge, &nameMetric)
		if err != nil {
			return nil, err
		}

		modelsMemStorage.Gauge[nameMetric] = valueMetricGauge
	}

	rowsMetrics, err = connect.db.QueryContext(ctx, "SELECT delta, name FROM metrics_counter")
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

func (connect *metricsRepository) GetMetricsTypeIDByName(ctx context.Context, nameMetric string) (int, error) {
	var id int

	fmt.Println("GetMetricsTypeIdByName Metrics => ", nameMetric)

	//rowArray := Row{}
	// Query for a value based on a single row.
	sqlStatement := `SELECT id FROM metrics_type WHERE name = $1`
	//row := db.QueryRow(sqlStatement, id).scan
	if err := connect.db.QueryRowContext(ctx, sqlStatement, nameMetric).Scan(&id); err != nil {
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

func (connect *metricsRepository) PingDatabase(ctx context.Context) error {
	return connect.db.Ping()
}

// Можно разбить на две части этот кода для большей детализации и ясности
func (connect *metricsRepository) SaveMetricsBatch(ctx context.Context, metrics []unserialize.Metrics) error {
	connect.mutex.Lock()

	tx, err := connect.db.Begin()

	if err != nil {
		return err
	}
	// можно вызвать Rollback в defer,
	// если Commit будет раньше, то откат проигнорируется
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		connect.mutex.Unlock()
	}()

	//var countUpdate = 0
	_, err = updateBatch(ctx, tx, metrics)

	if err != nil {
		fmt.Printf("error Update metrics: %v", err)
		return err
	}

	return nil
}

func insertBatch(ctx context.Context, tx *sql.Tx, insertMetrisAfterUpdate []unserialize.Metrics) error {

	if len(insertMetrisAfterUpdate) > 0 {

		gaugeInsert, err := tx.PrepareContext(context.TODO(),
			"INSERT INTO metrics_gauge(value, type_id, name)"+
				" VALUES($1,$2,$3)")
		if err != nil {
			fmt.Printf("error Set insert gauge metrics: %v", err)
			return err
		}

		counterInsert, err := tx.PrepareContext(context.TODO(),
			"INSERT INTO metrics_counter(delta, type_id, name)"+
				" VALUES($1,$2,$3)")
		if err != nil {
			fmt.Printf("error Set insert counter metrics: %v", err)
			return err
		}
		defer gaugeInsert.Close()
		defer counterInsert.Close()

		for _, v := range insertMetrisAfterUpdate {
			if v.MType == "gauge" {
				//fmt.Println("GAUGE VALUE SAVE METRICS", v.Value)
				_, err := gaugeInsert.ExecContext(context.TODO(), v.Value, 1, v.ID)
				if err != nil {
					fmt.Printf("error insert gauge metrics: %v", err)
					return err
				}
			}

			if v.MType == "counter" {
				_, err := counterInsert.ExecContext(context.TODO(), *v.Delta, 2, v.ID)
				if err != nil {
					fmt.Printf("error insert counter metrics: %v", err)
					return err
				}
			}

		}
	}

	return nil
}

func updateBatch(ctx context.Context, tx *sql.Tx, metrics []unserialize.Metrics) ([]unserialize.Metrics, error) {
	//var insertMetrisAfterUpdate []unserialize.Metrics

	gaugeUpdate, err := tx.PrepareContext(context.TODO(),
		"UPDATE metrics_gauge SET value = $1 WHERE type_id = $2 AND name = $3;")
	if err != nil {
		fmt.Printf("error SET update gauge metrics: %v", err)
		return nil, err
	}
	counterUpdate, err := tx.PrepareContext(context.TODO(),
		"UPDATE metrics_counter SET delta = delta + $1 WHERE type_id = $2 AND name = $3;")
	if err != nil {
		fmt.Printf("error Set update counter metrics: %v", err)
		return nil, err
	}
	//fmt.Println("Подготовили запросы:")

	defer gaugeUpdate.Close()
	defer counterUpdate.Close()

	for _, v := range metrics {
		if v.MType == "gauge" {
			res, err := gaugeUpdate.ExecContext(context.TODO(), *v.Value, 1, v.ID)
			if err != nil {
				return nil, err
			}
			count, err := res.RowsAffected()
			if err != nil {
				fmt.Printf("error update gauge metrics: %v", err)
				return nil, fmt.Errorf("error get rows affected metrics: %v", err)
			}
			if !(count > 0) {
				//insertMetrisAfterUpdate = append(insertMetrisAfterUpdate, v)
				insertBatch(ctx, tx, []unserialize.Metrics{v})
			}
		}

		if v.MType == "counter" {
			res, err := counterUpdate.ExecContext(context.TODO(), *v.Delta, 2, v.ID)
			if err != nil {
				return nil, err
			}
			count, err := res.RowsAffected()
			if err != nil {
				fmt.Printf("error update counter metrics: %v", err)
				return nil, fmt.Errorf("error get rows affected metrics: %v", err)
			}
			if !(count > 0) {
				//insertMetrisAfterUpdate = append(insertMetrisAfterUpdate, v)
				insertBatch(ctx, tx, []unserialize.Metrics{v})
			}
		}

	}

	return nil, nil
}
