package filestorage

import (
	"bufio"
	"context"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"go.uber.org/zap"
	"os"
	"sync"
)

type FileStorage struct {
	mx     sync.RWMutex
	log    *zap.Logger
	file   *os.File
	writer *bufio.Writer
	reader *bufio.Scanner
	conf   *StorageConf
}

type StorageConf struct {
	StoreInterval   int
	FileStorageFile string
	Restore         bool
	Enabled         bool // Флаг, указывающий, нужно ли сохранять метрики в storage
}

// NewPermanentStorageConf возвращает конфигурацию файлового хранилища
func NewPermanentStorageConf(storeInterval int, fileStorageFile string, restore bool) *StorageConf {
	return &StorageConf{
		StoreInterval:   storeInterval,
		FileStorageFile: fileStorageFile,
		Restore:         restore,
		Enabled:         true,
	}
}

// NewFileStorage фозвращает файловое хранилище
func NewFileStorage(conf *StorageConf, log *zap.Logger) (*FileStorage, error) {
	file, err := os.OpenFile(conf.FileStorageFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}
	fs := &FileStorage{
		file:   file,
		writer: bufio.NewWriter(file),
		reader: bufio.NewScanner(file),
		conf:   conf,
		log:    log,
	}

	return fs, nil
}

// GetConfig получение конфигурации
func (fs *FileStorage) GetConfig() *StorageConf {
	return fs.conf
}

// Stop остановка работы хранилища
func (fs *FileStorage) Stop() {

	fs.mx.Lock()
	defer fs.mx.Unlock()

	if err := fs.file.Close(); err != nil {
		fs.log.Error("fs.file.Close", zap.Error(err))
		return
	}

	fs.log.Info("file storage success stopped")

}

func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return models.ErrDeadlineContext
	default:
	}
	return nil
}
