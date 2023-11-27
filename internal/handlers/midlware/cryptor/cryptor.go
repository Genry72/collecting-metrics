package cryptor

import (
	"bytes"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type crypt struct {
	gin.ResponseWriter
	password string
	log      *zap.Logger
}

// Write добавляет заголовок HashSHA256 со значением хеша тела ответа
func (c *crypt) Write(b []byte) (int, error) {
	hashFromBody, err := cryptor.Encrypt(b, c.password)
	if err != nil {
		c.log.Error("cryptor.Encrypt", zap.Error(err))
		return 0, err
	}

	c.Header().Set(models.HeaderHash, hashFromBody)

	return c.ResponseWriter.Write(b)
}

/*
CheckHashFromHeader Проверка входящего хедера HashSHA256. Если передан, то сверяем сумму хэша с телом запроса.
Считаем хэш тела ответа и добавляем заголовок в ответ.
*/
func CheckHashFromHeader(log *zap.Logger, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerHash := c.Request.Header.Get(models.HeaderHash)
		if headerHash == "" {
			c.Next()
			return
		}

		body, err := c.GetRawData()

		if err != nil {
			log.Error("c.GetRawData", zap.Error(err))
			c.String(http.StatusBadRequest, models.ErrBadBody.Error())
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		hashFromBody, err := cryptor.Encrypt(body, password)
		if err != nil {
			log.Error("c.GetRawData", zap.Error(err))
			c.String(http.StatusInternalServerError, models.ErrBadBody.Error())
			return
		}

		if headerHash != hashFromBody {
			log.Error("headerHash != hashFromBody", zap.Error(models.ErrHashNotEqual))
			c.String(http.StatusBadRequest, models.ErrHashNotEqual.Error())
			return
		}

		t := &crypt{
			ResponseWriter: c.Writer,
			password:       password,
			log:            log,
		}

		c.Writer = t

		c.Next()

	}
}
