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
		updatedGaugeValue   float64
		updatedCounterValue int64
	)

	metric, err := ParseMetricInfo(ctx)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	resp := Metrics{
		ID:    metric.Name(),
		MType: string(metric.Type()),
	}

	metricsRepo := s.GetMetricsRepo()

	switch metric.Type() {
	case models.Gauge:
		updatedGaugeValue, err = metricsRepo.UpdateGaugeMetric(metric.Name(), metric.GaugeValue())
		resp.Value = &updatedGaugeValue
	case models.Counter:
		updatedCounterValue, err = metricsRepo.UpdateCounterMetric(metric.Name(), metric.CounterValue())
		resp.Delta = &updatedCounterValue
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, resp, nil
}
