// Package metrics содержит реализацию сервера метрик.
package metrics

import (
	"github.com/xantinium/metrix/internal/presentation/rpc/gen"
	"github.com/xantinium/metrix/internal/repository/metrics"
)

// New создаёт новый сервер метрик.
func New(repo *metrics.MetricsRepository) *MetricsServer {
	return &MetricsServer{
		repo: repo,
	}
}

// MetricsServer реализация сервера метрик.
type MetricsServer struct {
	gen.UnimplementedMetricsServer

	repo *metrics.MetricsRepository
}
