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
	metric, err := ParseMetricInfo(ctx)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	metricsRepo := s.GetMetricsRepo()

	switch metric.Type() {
	case models.Gauge:
		return getGaugeMetricHandler(metricsRepo, metric.Name())
	case models.Counter:
		return getCounterMetricHandler(metricsRepo, metric.Name())
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

func ParseMetricInfo(ctx *gin.Context) (models.MetricInfo, error) {
	var (
		err        error
		bodyBytes  []byte
		req        Metrics
		metricType models.MetricType
	)

	bodyBytes, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		return models.MetricInfo{}, err
	}

	err = easyjson.Unmarshal(bodyBytes, &req)
	if err != nil {
		return models.MetricInfo{}, err
	}

	metricType, err = parseType(req.MType)
	if err != nil {
		return models.MetricInfo{}, err
	}

	if metricType == models.Gauge {
		if req.Value == nil {
			return models.MetricInfo{}, fmt.Errorf("value is missing")
		}

		return models.NewGaugeMetric(req.ID, *req.Value), nil
	}

	if metricType == models.Counter {
		if req.Delta == nil {
			return models.MetricInfo{}, fmt.Errorf("value is missing")
		}

		return models.NewCounterMetric(req.ID, *req.Delta), nil
	}

	return models.MetricInfo{}, fmt.Errorf("invalid request")
}

//easyjson:json
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// parseType парсит тип метрики.
func parseType(maybeMetricType string) (models.MetricType, error) {
	switch maybeMetricType {
	case string(models.Gauge):
		return models.Gauge, nil
	case string(models.Counter):
		return models.Counter, nil
	default:
		return "", fmt.Errorf("unknown metric type: %q", maybeMetricType)
	}
}
