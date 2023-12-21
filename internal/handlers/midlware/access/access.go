package access

import (
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
)

/*
CheckIPAddress Является промежуточной функцией для фреймворка Gin. Она проверяет, находится ли IP-адрес клиента
в доверенной подсети. Если IP-адрес находится в подсети, то запрос разрешается,
в противном случае возвращается статус "Forbidden".
*/
func CheckIPAddress(log *zap.Logger, trustedSubnet string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if trustedSubnet == "" {
			c.Next()
			return
		}

		header := c.Request.Header.Get(models.HeaderTrustedSubnet)

		if header == "" {
			log.Error("empty header " + models.HeaderTrustedSubnet)
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}

		ip := net.ParseIP(header)

		_, subnet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			log.Error("net.ParseCIDR", zap.Error(err))
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}

		if subnet.Contains(ip) {
			c.Next()
			return
		}

		c.Status(http.StatusForbidden)
		c.Abort()
	}
}
