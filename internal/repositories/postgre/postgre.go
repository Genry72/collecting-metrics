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
)

type PGStorage struct {
	conn *sqlx.DB
	log  *zap.Logger
}

func NewPGStorage(dsn string, log *zap.Logger) (*PGStorage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	pg := &PGStorage{
		conn: db,
		log:  log,
	}

	if err := pg.migrate(); err != nil {
		return nil, err
	}

	return pg, nil
}

func (pg *PGStorage) Stop() {
	if err := pg.conn.Close(); err != nil {
		pg.log.Error(err.Error())
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
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			pg.log.Error("migrate.Rollback", zap.Error(err))
		}
	}()

	if _, err := tx.Exec(migrationQuery); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	pg.log.Info("success migration")

	return nil
}

func (pg *PGStorage) SetMetric(ctx context.Context, metrics ...*models.Metric) ([]*models.Metric, error) {

	result := make([]*models.Metric, 0, len(metrics))

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
RETURNING
    name,
    type,
    delta,
    value
`
	tx, err := pg.conn.Beginx()
	if err != nil {
		return nil, fmt.Errorf("SetMetric.Begin: %w", err)
	}

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("SetMetric.PreparexContext: %w", err)
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
			return nil, err
		}

		if metric.MType == models.MetricTypeCounter {
			oldMetric, err := pg.GetMetricValue(ctx, metric.MType, metric.ID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("SetMetric: %w", err)
			}

			oldValue := int64(0)

			if oldMetric != nil {
				oldValue = *oldMetric.Delta
			}

			v := oldValue + *metric.Delta
			metric.Delta = &v
		}

		var m models.Metric

		row := stmt.QueryRowxContext(ctx, metric.ID, metric.MType, metric.Delta, metric.Value)

		if err := row.StructScan(&m); err != nil {
			fmt.Printf("%+v\n", *metric.Delta)
			return nil, fmt.Errorf("row.StructScan: %w", err)
		}

		if err := row.Err(); err != nil {
			return nil, err
		}

		result = append(result, &m)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("SetMetric.Commit: %w", err)
	}

	return result, nil
}

func (pg *PGStorage) GetMetricValue(ctx context.Context,
	metricType models.MetricType, metricName models.MetricName) (*models.Metric, error) {
	if err := checkMetricType(metricType); err != nil {
		return nil, err
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
		return nil, fmt.Errorf("GetMetricValue.QueryRowxContext: %w", err)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("GetMetricValue: %w", err)
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
		return nil, err
	}
	return result, nil
}

func (pg *PGStorage) SetAllMetrics(ctx context.Context, metrics []*models.Metric) error {
	for i := range metrics {
		if _, err := pg.SetMetric(ctx, metrics[i]); err != nil {
			return fmt.Errorf("SetAllMetrics: %w", err)
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
