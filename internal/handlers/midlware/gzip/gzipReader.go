package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
)

// Разорхивируем запросы от клиента
type gzipReader struct {
	ginReader io.ReadCloser
	gzReader  *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("gzip: %w", err)
	}

	return &gzipReader{
		ginReader: r,
		gzReader:  gzReader,
	}, nil
}

func (c *gzipReader) Read(p []byte) (n int, err error) {
	return c.gzReader.Read(p)
}

func (c *gzipReader) Close() error {
	if err := c.ginReader.Close(); err != nil {
		return err
	}
	return c.gzReader.Close()
}
