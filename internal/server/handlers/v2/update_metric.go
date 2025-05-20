package v2handlers

import (
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
		ID:    req.Metric.ID(),
		MType: string(req.Metric.Type()),
	}

	metricsRepo := s.GetMetricsRepo()

	switch req.Metric.Type() {
	case models.Gauge:
		updatedGaugeValue, err = metricsRepo.UpdateGaugeMetric(ctx, req.Metric.ID(), req.Metric.GaugeValue())
		resp.Value = &updatedGaugeValue
	case models.Counter:
		updatedCounterValue, err = metricsRepo.UpdateCounterMetric(ctx, req.Metric.ID(), req.Metric.CounterValue())
		resp.Delta = &updatedCounterValue
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, resp, nil
}

type UpdateMetricRequest struct {
	Metric models.MetricInfo
}

// ParseUpdateMetricRequest парсит запрос на обновление метрики.
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

	req.Metric, err = parseMetric(rawReq)
	if err != nil {
		return UpdateMetricRequest{}, err
	}

	return req, nil
}
