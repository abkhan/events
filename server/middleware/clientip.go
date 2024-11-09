package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const xffHeader = "X-Forwarded-For"

// DetectClientIP lookup client ip within X-Forwarded-For value according to next rules:
//   - last item is alb (aws proxy) if there are more then one ip within list, otherwise we consider
//     it as a client ip
//   - item before last is a client ip or client proxy (we consider it as client ip)
//   - any possible fake ip or real client ip but we avoid using other items
func DetectClientIP(c *gin.Context) {
	xff := strings.TrimSpace(c.Request.Header.Get(xffHeader))
	if xff == "" {
		return
	}

	if ips := strings.Split(xff, ","); len(ips) > 1 {
		c.Request.Header.Del(xffHeader)
		c.Request.Header.Set(xffHeader, strings.TrimSpace(ips[len(ips)-2]))
	}
}
