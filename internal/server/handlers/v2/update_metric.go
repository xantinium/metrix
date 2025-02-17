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
		err          error
		bodyBytes    []byte
		req          Metrics
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
