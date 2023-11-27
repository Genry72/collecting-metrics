package gzip

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

// Сжимаем ответы от сервера
type gzipWriter struct {
	ginRW gin.ResponseWriter
	gzWR  *gzip.Writer
}

func newGzipWriter(w gin.ResponseWriter) (*gzipWriter, error) {
	gzWR, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("gzip.NewWriterLevel: %w", err)
	}

	return &gzipWriter{
		ginRW: w,
		gzWR:  gzWR,
	}, nil
}

// Write Имплементация ResponseWriter. Добавляеет заголовок Content = gzip
func (g *gzipWriter) Write(b []byte) (int, error) {
	//if len(b) > 1400 {
	//	g.Header().Set("Content-Encoding", "gzip")
	//	return g.gzWR.Write(b)
	//}
	//
	//return g.ginRW.Write(b)
	g.Header().Set("Content-Encoding", "gzip")
	return g.gzWR.Write(b)
}

func (g *gzipWriter) Close() error {
	return g.gzWR.Close()
}

func (g *gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return g.ginRW.Hijack()
}

func (g *gzipWriter) WriteHeaderNow() {
	g.ginRW.WriteHeaderNow()
}

func (g *gzipWriter) CloseNotify() <-chan bool {
	return g.ginRW.CloseNotify()
}

func (g *gzipWriter) Status() int {
	return g.ginRW.Status()
}

func (g *gzipWriter) Size() int {
	return g.ginRW.Size()
}

func (g *gzipWriter) Written() bool {
	return g.ginRW.Written()
}

func (g *gzipWriter) Pusher() http.Pusher {
	return g.ginRW.Pusher()
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.ginRW.Write([]byte(s))
}

func (g *gzipWriter) WriteHeader(code int) {
	g.ginRW.WriteHeader(code)
}

func (g *gzipWriter) Header() http.Header {
	return g.ginRW.Header()
}

func (g *gzipWriter) Flush() {
	g.ginRW.Flush()
}
