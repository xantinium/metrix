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
	req, err := ParseGetMetricRequest(ctx)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	metricsRepo := s.GetMetricsRepo()

	switch req.MetricType {
	case models.Gauge:
		return getGaugeMetricHandler(metricsRepo, req.MetricName)
	case models.Counter:
		return getCounterMetricHandler(metricsRepo, req.MetricName)
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

type GetMetricsRequest struct {
	MetricName string
	MetricType models.MetricType
}

func ParseGetMetricRequest(ctx *gin.Context) (GetMetricsRequest, error) {
	var (
		err       error
		bodyBytes []byte
		rawReq    Metrics
		req       GetMetricsRequest
	)

	bodyBytes, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		return GetMetricsRequest{}, err
	}

	err = easyjson.Unmarshal(bodyBytes, &rawReq)
	if err != nil {
		return GetMetricsRequest{}, err
	}

	req.MetricName = rawReq.ID
	if req.MetricName == "" {
		return GetMetricsRequest{}, fmt.Errorf("metric id cannot be empty")
	}

	req.MetricType, err = parseType(rawReq.MType)
	if err != nil {
		return GetMetricsRequest{}, err
	}

	return req, nil
}
