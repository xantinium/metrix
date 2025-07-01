package middlewares

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/tools"
)

// NetGuardMiddleware мидлварь для проверки IP-адреса клиента,
// при помощи переданной доверенной подсети.
func NetGuardMiddleware(trustedSubnet *net.IPNet) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		realIP := ctx.GetHeader(tools.HeaderXRealIP)
		if realIP == "" {
			logger.Infof("%s header is missing, but required", tools.HeaderXRealIP)
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ip := net.ParseIP(realIP)
		if ip == nil {
			logger.Infof("invalid ip address in %s header", tools.HeaderXRealIP)
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !trustedSubnet.Contains(ip) {
			logger.Infof("ip address %q not in trusted subnet", ip.String())
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
