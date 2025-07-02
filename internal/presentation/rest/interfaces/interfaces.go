// Package interfaces содержит интерфейс http-сервера.
package interfaces

import (
	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/repository/metrics"
)

// Server интерфейс сервера, доступного в хендлерах.
type Server interface {
	GetInternalRouter() *gin.Engine
	GetMetricsRepo() *metrics.MetricsRepository
}
