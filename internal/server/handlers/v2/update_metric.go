package v2handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/server/interfaces"
)

// UpdateMetricHandler реализация хендлера для обновления метрик.
func UpdateMetricHandler(ctx *gin.Context, s interfaces.Server) (int, easyjson.Marshaler, error) {
	var (
		err          error
		bodyBytes    []byte
		req          updateMetricsRequest
		metricType   models.MetricType
		gaugeValue   float64
		counterValue int64
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
		gaugeValue, err = req.ParseGaugeValue()
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		err = metricsRepo.UpdateGaugeMetric(req.ID, gaugeValue)
	case models.Counter:
		counterValue, err = req.ParseCounterValue()
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		err = metricsRepo.UpdateCounterMetric(req.ID, counterValue)
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}

//easyjson:json
type updateMetricsRequest struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// ParseType парсит тип метрики.
func (req updateMetricsRequest) ParseType() (models.MetricType, error) {
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
func (req updateMetricsRequest) ParseGaugeValue() (float64, error) {
	if req.Value == nil {
		return 0, fmt.Errorf("value is missing")
	}

	return *req.Value, nil
}

// ParseGaugeValue парсит значение для метрики типа Counter.
func (req updateMetricsRequest) ParseCounterValue() (int64, error) {
	if req.Delta == nil {
		return 0, fmt.Errorf("value is missing")
	}

	return *req.Delta, nil
}
