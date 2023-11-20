package postgre

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories/filestorage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

type PGStorage struct {
	conn *sqlx.DB
	log  *zap.Logger
}

func NewPGStorage(dsn string, log *zap.Logger) (*PGStorage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(1 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	pg := &PGStorage{
		conn: db,
		log:  log,
	}

	if err := pg.migrate(); err != nil {
		return nil, fmt.Errorf("pg.migrate: %w", err)
	}

	return pg, nil
}

func (pg *PGStorage) Stop() {
	if err := pg.conn.Close(); err != nil {
		pg.log.Error("g.conn.Close", zap.Error(err))
		return
	}

	pg.log.Info("Database success closed")
}

func (pg *PGStorage) Ping() error {
	if pg == nil {
		return fmt.Errorf("database not connected")
	}
	return pg.conn.Ping()
}

func (pg *PGStorage) migrate() error {
	tx, err := pg.conn.Begin()
	if err != nil {
		return fmt.Errorf("pg.conn.Begin: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			pg.log.Error("tx.Rollback", zap.Error(err))
		}
	}()

	if _, err := tx.Exec(migrationQuery); err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit:%w", err)
	}

	pg.log.Info("success migration")

	return nil
}

func (pg *PGStorage) SetMetric(ctx context.Context, metrics ...*models.Metric) error {
	query := `
INSERT INTO metrics (name,
                     type,
                     delta,
                     value)
VALUES ($1,
        $2,
        $3,
        $4)
ON CONFLICT (name, type)
    DO UPDATE SET name  = EXCLUDED.name,
                  type= EXCLUDED.type,
                  delta= EXCLUDED.delta,
                  value = EXCLUDED.value
`
	tx, err := pg.conn.Beginx()
	if err != nil {
		return fmt.Errorf("pg.conn.Beginx: %w", err)
	}

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("tx.PreparexContext: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			pg.log.Error("tx.Rollback()", zap.Error(err))
		}
		if err := stmt.Close(); err != nil {
			pg.log.Error("stmt.Close", zap.Error(err))
		}
	}()

	for i := range metrics {
		metric := metrics[i]

		if err := checkMetricType(metric.MType); err != nil {
			return err
		}

		if metric.MType == models.MetricTypeCounter {
			oldMetric, err := pg.GetMetricValue(ctx, metric.MType, metric.ID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("pg.GetMetricValue: %w", err)
			}

			oldValue := int64(0)

			if oldMetric != nil {
				oldValue = *oldMetric.Delta
			}

			v := oldValue + *metric.Delta
			metric.Delta = &v
		}

		_, err = stmt.ExecContext(ctx, metric.ID, metric.MType, metric.Delta, metric.Value)

		if err != nil {
			return fmt.Errorf("stmt.ExecContext: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}

func (pg *PGStorage) GetMetricValue(ctx context.Context,
	metricType models.MetricType, metricName models.MetricName) (*models.Metric, error) {
	if err := checkMetricType(metricType); err != nil {
		return nil, fmt.Errorf("checkMetricType: %w", err)
	}

	query := `
select name, type, delta, value
from metrics
where name = $1
  and type = $2
`

	result := models.Metric{}

	row := pg.conn.QueryRowxContext(ctx, query, metricName, metricType)
	if err := row.StructScan(&result); err != nil {
		return nil, fmt.Errorf("row.StructScant: %w", err)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("row.Err: %w", err)
	}

	return &result, nil
}

func (pg *PGStorage) GetAllMetrics(ctx context.Context) ([]*models.Metric, error) {
	query := `
select name, type, delta, value
from metrics
`

	result := make([]*models.Metric, 0)
	if err := pg.conn.SelectContext(ctx, &result, query); err != nil {
		return nil, fmt.Errorf("pg.conn.SelectContext: %w", err)
	}
	return result, nil
}

func (pg *PGStorage) SetAllMetrics(ctx context.Context, metrics []*models.Metric) error {
	for i := range metrics {
		if err := pg.SetMetric(ctx, metrics[i]); err != nil {
			return fmt.Errorf("pg.SetMetric: %w", err)
		}
	}
	return nil
}

func (pg *PGStorage) GetConfig() *filestorage.StorageConf {
	return &filestorage.StorageConf{
		StoreInterval:   0,
		FileStorageFile: "",
		Restore:         false,
		Enabled:         false,
	}
}

func checkMetricType(metricType models.MetricType) error {
	switch metricType {
	case models.MetricTypeCounter:
		return nil
	case models.MetricTypeGauge:
		return nil
	default:
		return fmt.Errorf("%w: %s", models.ErrBadMetricType, metricType)
	}
}
