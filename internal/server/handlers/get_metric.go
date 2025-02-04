package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/repository/metrics"
)

// GetMetricHandler реализация хендлера для получения метрик.
func GetMetricHandler(ctx *gin.Context, s server) (int, []byte, error) {
	req, err := parseGetMetricRequest(ctx)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	metricsRepo := s.GetMetricsRepo()

	switch req.metricType {
	case models.Gauge:
		return getGaugeMetricHandler(metricsRepo, req.metricName)
	case models.Counter:
		return getCounterMetricHandler(metricsRepo, req.metricName)
	default:
		// Попасть сюда невозможно, из-за валидации запроса.
		return http.StatusInternalServerError, nil, fmt.Errorf("unknown metric type")
	}
}

func getGaugeMetricHandler(repo *metrics.MetricsRepository, name string) (int, []byte, error) {
	value, err := repo.GetGaugeMetric(name)
	if err != nil {
		if err == models.ErrNotFound {
			return http.StatusNotFound, nil, err
		}

		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, []byte(floatToStr(value)), nil
}

func getCounterMetricHandler(repo *metrics.MetricsRepository, name string) (int, []byte, error) {
	value, err := repo.GetCounterMetric(name)
	if err != nil {
		if err == models.ErrNotFound {
			return http.StatusNotFound, nil, err
		}

		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, []byte(intToStr(value)), nil
}

func floatToStr(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func intToStr(v int64) string {
	return strconv.FormatInt(v, 10)
}

// getMetricRequest структура запроса обновления метрик.
type getMetricRequest struct {
	metricType models.MetricType
	metricName string
}

// parseGetMetricRequest парсит сырой HTTP-запрос в структуру запроса.
func parseGetMetricRequest(r *gin.Context) (getMetricRequest, error) {
	var (
		err                         error
		maybeMetricType, metricName string
		metricType                  models.MetricType
	)

	maybeMetricType = r.Param("type")
	if maybeMetricType == "" {
		return getMetricRequest{}, fmt.Errorf("metric type is missing")
	}

	metricType, err = models.ParseStringAsMetricType(maybeMetricType)
	if err != nil {
		return getMetricRequest{}, err
	}

	metricName = r.Param("name")
	if metricName == "" {
		return getMetricRequest{}, fmt.Errorf("metric name is missing")
	}

	return getMetricRequest{
		metricType: metricType,
		metricName: metricName,
	}, nil
}
