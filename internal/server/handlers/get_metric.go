package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/repository/metrics"
	"github.com/xantinium/metrix/internal/server/interfaces"
	"github.com/xantinium/metrix/internal/tools"
)

// GetMetricHandler реализация хендлера для получения метрик.
// @Tags Metrics_Legacy
// @Summary Запрос на получение метрики по ID
// @Description Запрос на получение метрики по ID
// @ID getMetricLegacy
// @Produce text/plain
// @Param metric_type path string true "Тип метрики"
// @Param metric_id path string true "Идентификатор метрики"
// @Success 200 {string} string
// @Failure 404 {string} string "Метрика не найдена"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /value/{metric_type}/{metric_id} [get]
func GetMetricHandler(ctx *gin.Context, s interfaces.Server) (int, string, error) {
	req, err := parseGetMetricRequest(ctx)
	if err != nil {
		return http.StatusBadRequest, "", err
	}

	metricsRepo := s.GetMetricsRepo()

	switch req.metricType {
	case models.Gauge:
		return getGaugeMetricHandler(ctx, metricsRepo, req.metricID)
	case models.Counter:
		return getCounterMetricHandler(ctx, metricsRepo, req.metricID)
	default:
		// Попасть сюда невозможно, из-за валидации запроса.
		return http.StatusInternalServerError, "", fmt.Errorf("unknown metric type")
	}
}

func getGaugeMetricHandler(ctx context.Context, repo *metrics.MetricsRepository, id string) (int, string, error) {
	value, err := repo.GetGaugeMetric(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return http.StatusNotFound, "", err
		}

		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, tools.FloatToStr(value), nil
}

func getCounterMetricHandler(ctx context.Context, repo *metrics.MetricsRepository, id string) (int, string, error) {
	value, err := repo.GetCounterMetric(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return http.StatusNotFound, "", err
		}

		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, tools.IntToStr(value), nil
}

// getMetricRequest структура запроса обновления метрик.
type getMetricRequest struct {
	metricType models.MetricType
	metricID   string
}

// parseGetMetricRequest парсит сырой HTTP-запрос в структуру запроса.
func parseGetMetricRequest(r *gin.Context) (getMetricRequest, error) {
	var (
		err                       error
		maybeMetricType, metricID string
		metricType                models.MetricType
	)

	maybeMetricType = r.Param("type")
	if maybeMetricType == "" {
		return getMetricRequest{}, fmt.Errorf("metric type is missing")
	}

	metricType, err = models.ParseStringAsMetricType(maybeMetricType)
	if err != nil {
		return getMetricRequest{}, err
	}

	metricID = r.Param("id")
	if metricID == "" {
		return getMetricRequest{}, fmt.Errorf("metric id is missing")
	}

	return getMetricRequest{
		metricType: metricType,
		metricID:   metricID,
	}, nil
}
