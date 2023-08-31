package postgre

import (
	"fmt"
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
