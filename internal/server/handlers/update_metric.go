package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/server/interfaces"
	"github.com/xantinium/metrix/internal/tools"
)

// UpdateMetricHandler реализация хендлера для обновления метрик.
func UpdateMetricHandler(ctx *gin.Context, s interfaces.Server) (int, string, error) {
	req, err := parseUpdateMetricRequest(ctx)
	if err != nil {
		return http.StatusBadRequest, "", err
	}

	metricsRepo := s.GetMetricsRepo()

	switch req.metricType {
	case models.Gauge:
		err = metricsRepo.UpdateGaugeMetric(req.metricName, req.metricValue)
	case models.Counter:
		err = metricsRepo.UpdateCounterMetric(req.metricName, int64(req.metricValue))
	}

	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, "", nil
}

// updateMetricRequest структура запроса обновления метрик.
type updateMetricRequest struct {
	metricType  models.MetricType
	metricName  string
	metricValue float64
}

// parseUpdateMetricRequest парсит сырой HTTP-запрос в структуру запроса.
func parseUpdateMetricRequest(r *gin.Context) (updateMetricRequest, error) {
	var (
		err                                           error
		maybeMetricType, metricName, maybeMetricValue string
		metricType                                    models.MetricType
		metricValue                                   float64
	)

	maybeMetricType = r.Param("type")
	if maybeMetricType == "" {
		return updateMetricRequest{}, fmt.Errorf("metric type is missing")
	}

	metricType, err = models.ParseStringAsMetricType(maybeMetricType)
	if err != nil {
		return updateMetricRequest{}, err
	}

	metricName = r.Param("name")
	if metricName == "" {
		return updateMetricRequest{}, fmt.Errorf("metric name is missing")
	}

	maybeMetricValue = r.Param("value")
	if maybeMetricValue == "" {
		return updateMetricRequest{}, fmt.Errorf("metric value is missing")
	}

	metricValue, err = tools.StrToFloat(maybeMetricValue)
	if err != nil {
		return updateMetricRequest{}, fmt.Errorf("invalid metric value")
	}

	return updateMetricRequest{
		metricType:  metricType,
		metricName:  metricName,
		metricValue: metricValue,
	}, nil
}
