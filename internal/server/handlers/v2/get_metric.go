package v2handlers

import (
	"context"
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
		return getGaugeMetricHandler(ctx, metricsRepo, req.MetricID)
	case models.Counter:
		return getCounterMetricHandler(ctx, metricsRepo, req.MetricID)
	default:
		// Попасть сюда невозможно, из-за валидации запроса.
		return http.StatusInternalServerError, nil, fmt.Errorf("unknown metric type")
	}
}

func getGaugeMetricHandler(ctx context.Context, repo *metrics.MetricsRepository, id string) (int, easyjson.Marshaler, error) {
	value, err := repo.GetGaugeMetric(ctx, id)
	if err != nil {
		if err == models.ErrNotFound {
			return http.StatusNotFound, nil, err
		}

		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, Metrics{
		ID:    id,
		MType: string(models.Gauge),
		Value: &value,
	}, nil
}

func getCounterMetricHandler(ctx context.Context, repo *metrics.MetricsRepository, id string) (int, easyjson.Marshaler, error) {
	value, err := repo.GetCounterMetric(ctx, id)
	if err != nil {
		if err == models.ErrNotFound {
			return http.StatusNotFound, nil, err
		}

		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, Metrics{
		ID:    id,
		MType: string(models.Counter),
		Delta: &value,
	}, nil
}

type GetMetricsRequest struct {
	MetricID   string
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

	req.MetricID = rawReq.ID
	if req.MetricID == "" {
		return GetMetricsRequest{}, fmt.Errorf("metric id cannot be empty")
	}

	req.MetricType, err = parseType(rawReq.MType)
	if err != nil {
		return GetMetricsRequest{}, err
	}

	return req, nil
}
