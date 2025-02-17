package v2handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/server/interfaces"
)

// UpdateMetricHandler реализация хендлера для обновления метрик.
func UpdateMetricHandler(ctx *gin.Context, s interfaces.Server) (int, easyjson.Marshaler, error) {
	var (
		err        error
		req        Metrics
		metricType models.MetricType
	)

	metric, err := ParseMetricInfo(ctx)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	metricsRepo := s.GetMetricsRepo()

	switch metricType {
	case models.Gauge:
		*req.Value, err = metricsRepo.UpdateGaugeMetric(req.ID, metric.GaugeValue())
	case models.Counter:
		*req.Delta, err = metricsRepo.UpdateCounterMetric(req.ID, metric.CounterValue())
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, req, nil
}
