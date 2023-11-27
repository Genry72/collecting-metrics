package gzip

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Gzip Возвращает обработчик gin, который проверяет, поддерживает ли клиент сжатие gzip, и если да,
// то оборачивает ответ сервера в gzip.Writer. Также проверяется, отправляет ли клиент данные в сжатом виде,
// и если да, то оборачивает тело запроса в gzip.Reader с поддержкой декомпрессии.
func Gzip(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := c.Request.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip {
			gz, err := newGzipWriter(c.Writer)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			defer func() {
				if err := gz.Close(); err != nil {
					log.Error(err.Error())
				}
			}()

			c.Writer = gz
		}

		// проверяем что клиент отправил данные в сжатом виде
		contentEncoding := c.Request.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newGzipReader(c.Request.Body)
			if err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			c.Request.Body = cr
			defer func() {
				if err := cr.Close(); err != nil {
					log.Error(err.Error())
				}
			}()
		}
		c.Next()
	}
}
