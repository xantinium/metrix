package v2handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/repository/metrics"
	"github.com/xantinium/metrix/internal/server/interfaces"
)

// GetMetricHandler реализация хендлера для получения метрик.
func GetMetricHandler(ctx *gin.Context, s interfaces.Server) (int, easyjson.Marshaler, error) {
	var (
		err        error
		bodyBytes  []byte
		req        Metrics
		metricType models.MetricType
	)

	bodyBytes, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	err = easyjson.Unmarshal(bodyBytes, &req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	metricType, err = req.ParseType()
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	metricsRepo := s.GetMetricsRepo()

	switch metricType {
	case models.Gauge:
		return getGaugeMetricHandler(metricsRepo, req.ID)
	case models.Counter:
		return getCounterMetricHandler(metricsRepo, req.ID)
	default:
		// Попасть сюда невозможно, из-за валидации запроса.
		return http.StatusInternalServerError, nil, fmt.Errorf("unknown metric type")
	}
}

func getGaugeMetricHandler(repo *metrics.MetricsRepository, name string) (int, easyjson.Marshaler, error) {
	value, err := repo.GetGaugeMetric(name)
	if err != nil {
		if err == models.ErrNotFound {
			return http.StatusNotFound, nil, err
		}

		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, Metrics{
		ID:    name,
		MType: string(models.Gauge),
		Value: &value,
	}, nil
}

func getCounterMetricHandler(repo *metrics.MetricsRepository, name string) (int, easyjson.Marshaler, error) {
	value, err := repo.GetCounterMetric(name)
	if err != nil {
		if err == models.ErrNotFound {
			return http.StatusNotFound, nil, err
		}

		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, Metrics{
		ID:    name,
		MType: string(models.Counter),
		Delta: &value,
	}, nil
}

//easyjson:json
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// ParseType парсит тип метрики.
func (req Metrics) ParseType() (models.MetricType, error) {
	switch req.MType {
	case string(models.Gauge):
		return models.Gauge, nil
	case string(models.Counter):
		return models.Counter, nil
	default:
		return "", fmt.Errorf("unknown metric type")
	}
}

// ParseGaugeValue парсит значение для метрики типа Gauge.
func (req Metrics) ParseGaugeValue() (float64, error) {
	if req.Value == nil {
		return 0, fmt.Errorf("value is missing")
	}

	return *req.Value, nil
}

// ParseGaugeValue парсит значение для метрики типа Counter.
func (req Metrics) ParseCounterValue() (int64, error) {
	if req.Delta == nil {
		return 0, fmt.Errorf("value is missing")
	}

	return *req.Delta, nil
}
