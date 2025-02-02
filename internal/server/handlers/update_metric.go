package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/xantinium/metrix/internal/models"
)

// UpdateMetric реализация хендлера для обновления метрик.
func UpdateMetric(s server, r *http.Request) (int, []byte, error) {
	req, err := parseUpdateMetricRequest(r)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	metricsRepo := s.GetMetricsRepo()

	switch req.metricType {
	case models.GAUGE:
		err = metricsRepo.UpdateGaugeMetric(req.metricName, req.metricValue)
	case models.COUNTER:
		err = metricsRepo.UpdateCounterMetric(req.metricName, int64(req.metricValue))
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}

// updateMetricRequest структура запроса обновления метрик.
type updateMetricRequest struct {
	metricType  models.MetricType
	metricName  string
	metricValue float64
}

// parseUpdateMetricRequest парсит сырой HTTP-запрос в структуру запроса.
func parseUpdateMetricRequest(r *http.Request) (updateMetricRequest, error) {
	var (
		err                                           error
		maybeMetricType, metricName, maybeMetricValue string
		metricType                                    models.MetricType
		metricValue                                   float64
	)

	maybeMetricType = r.PathValue("type")
	if maybeMetricType == "" {
		return updateMetricRequest{}, fmt.Errorf("metric type is missing")
	}

	metricType, err = models.ParseStringAsMetricType(maybeMetricType)
	if err != nil {
		return updateMetricRequest{}, err
	}

	metricName = r.PathValue("name")
	if metricName == "" {
		return updateMetricRequest{}, fmt.Errorf("metric name is missing")
	}

	maybeMetricValue = r.PathValue("value")
	if maybeMetricValue == "" {
		return updateMetricRequest{}, fmt.Errorf("metric value is missing")
	}

	metricValue, err = strconv.ParseFloat(maybeMetricValue, 64)
	if err != nil {
		return updateMetricRequest{}, fmt.Errorf("invalid metric value")
	}

	return updateMetricRequest{
		metricType:  metricType,
		metricName:  metricName,
		metricValue: metricValue,
	}, nil
}
