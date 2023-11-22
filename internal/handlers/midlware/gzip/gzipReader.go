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
		return nil, fmt.Errorf("gzip.NewReader: %w", err)
	}

	return &gzipReader{
		ginReader: r,
		gzReader:  gzReader,
	}, nil
}

// Read implements io.Reader, reading uncompressed bytes from its underlying Reader.
func (c *gzipReader) Read(p []byte) (n int, err error) {
	return c.gzReader.Read(p)
}

// Close closes the Reader. It does not close the underlying io.Reader. In order for the GZIP checksum to be
// verified, the reader must be fully consumed until the io.EOF.
func (c *gzipReader) Close() error {
	if err := c.ginReader.Close(); err != nil {
		return err
	}
	return c.gzReader.Close()
}
