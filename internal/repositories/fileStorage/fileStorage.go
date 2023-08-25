package fileStorage

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sync"
)

type FileStorage struct {
	mx       sync.RWMutex
	fileName string
	log      *zap.Logger
	file     *os.File
	writer   *bufio.Writer
	reader   *bufio.Scanner
	started  bool
	conf     *StorageConf
}

type StorageConf struct {
	StoreInterval   int
	FileStorageFile string
	Restore         bool
}

func NewPermanentStorageConf(storeInterval int, fileStorageFile string, restore bool) *StorageConf {
	return &StorageConf{
		StoreInterval:   storeInterval,
		FileStorageFile: fileStorageFile,
		Restore:         restore,
	}
}

func NewFileStorage(conf *StorageConf, log *zap.Logger) *FileStorage {
	return &FileStorage{
		fileName: conf.FileStorageFile,
		conf:     conf,
		log:      log,
	}
}

func (fs *FileStorage) GetConfig() *StorageConf {
	return fs.conf
}

func (fs *FileStorage) Stop() error {
	if !fs.IsStarted() {
		return nil
	}

	fs.mx.Lock()
	defer fs.mx.Unlock()

	if err := fs.file.Close(); err != nil {
		return err
	}
	fs.file = nil
	fs.writer = nil
	fs.reader = nil
	fs.started = false

	fs.log.Info("file storage success stopped")
	return nil
}

func (fs *FileStorage) Start() error {
	if fs.IsStarted() {
		return nil
	}

	fs.mx.Lock()
	defer fs.mx.Unlock()

	file, err := os.OpenFile(fs.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("NewFileStorage: %w", err)
	}
	fs.file = file
	fs.writer = bufio.NewWriter(file)
	fs.reader = bufio.NewScanner(file)
	fs.started = true

	fs.log.Info("file storage success started")
	return nil
}

func (fs *FileStorage) IsStarted() bool {
	fs.mx.RLock()
	defer fs.mx.RUnlock()
	return fs.started
}
