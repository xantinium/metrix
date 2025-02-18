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
		updatedGaugeValue   float64
		updatedCounterValue int64
	)

	req, err := ParseUpdateMetricRequest(ctx)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	resp := Metrics{
		ID:    req.MetricName,
		MType: string(req.MetricType),
	}

	metricsRepo := s.GetMetricsRepo()

	switch req.MetricType {
	case models.Gauge:
		updatedGaugeValue, err = metricsRepo.UpdateGaugeMetric(req.MetricName, req.GaugeValue)
		resp.Value = &updatedGaugeValue
	case models.Counter:
		updatedCounterValue, err = metricsRepo.UpdateCounterMetric(req.MetricName, req.CounterValue)
		resp.Delta = &updatedCounterValue
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, resp, nil
}

type UpdateMetricRequest struct {
	MetricName   string
	MetricType   models.MetricType
	GaugeValue   float64
	CounterValue int64
}

func ParseUpdateMetricRequest(ctx *gin.Context) (UpdateMetricRequest, error) {
	var (
		err       error
		bodyBytes []byte
		rawReq    Metrics
		req       UpdateMetricRequest
	)

	bodyBytes, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		return UpdateMetricRequest{}, err
	}

	err = easyjson.Unmarshal(bodyBytes, &rawReq)
	if err != nil {
		return UpdateMetricRequest{}, err
	}

	req.MetricName = rawReq.ID
	if req.MetricName == "" {
		return UpdateMetricRequest{}, fmt.Errorf("metric id cannot be empty")
	}

	req.MetricType, err = parseType(rawReq.MType)
	if err != nil {
		return UpdateMetricRequest{}, err
	}

	if req.MetricType == models.Gauge {
		if rawReq.Value == nil {
			return UpdateMetricRequest{}, fmt.Errorf("value is missing")
		} else {
			req.GaugeValue = *rawReq.Value
		}
	}

	if req.MetricType == models.Counter {
		if rawReq.Delta == nil {
			return UpdateMetricRequest{}, fmt.Errorf("value is missing")
		} else {
			req.CounterValue = *rawReq.Delta
		}
	}

	return req, nil
}
