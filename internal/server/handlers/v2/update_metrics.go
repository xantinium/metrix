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

// UpdateMetricsHandler реализация хендлера для батчевого обновления метрик.
func UpdateMetricsHandler(ctx *gin.Context, s interfaces.Server) (int, easyjson.Marshaler, error) {
	req, err := ParseUpdateMetricsRequest(ctx)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	err = s.GetMetricsRepo().UpdateMetrics(ctx, req.Metrics)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}

type UpdateMetricsRequest struct {
	Metrics []models.MetricInfo
}

func ParseUpdateMetricsRequest(ctx *gin.Context) (UpdateMetricsRequest, error) {
	var (
		err       error
		bodyBytes []byte
		rawReq    MetricsBatch
		req       UpdateMetricsRequest
	)

	bodyBytes, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		return UpdateMetricsRequest{}, err
	}

	err = easyjson.Unmarshal(bodyBytes, &rawReq)
	if err != nil {
		return UpdateMetricsRequest{}, err
	}

	req.Metrics = make([]models.MetricInfo, len(rawReq))
	for i, metric := range rawReq {
		var (
			metricID   string
			metricType models.MetricType
			metricInfo models.MetricInfo
		)

		metricID = metric.ID
		if metricID == "" {
			return UpdateMetricsRequest{}, fmt.Errorf("metric id cannot be empty")
		}

		metricType, err = parseType(metric.MType)
		if err != nil {
			return UpdateMetricsRequest{}, err
		}

		switch metricType {
		case models.Gauge:
			if metric.Value == nil {
				err = fmt.Errorf("value is missing")
			} else {
				metricInfo = models.NewGaugeMetric(metricID, *metric.Value)
			}
		case models.Counter:
			if metric.Delta == nil {
				err = fmt.Errorf("value is missing")
			} else {
				metricInfo = models.NewCounterMetric(metricID, *metric.Delta)
			}
		default:
			err = fmt.Errorf("unknown metric type: %q", metricType)
		}
		if err != nil {
			return UpdateMetricsRequest{}, err
		}

		req.Metrics[i] = metricInfo
	}

	return req, nil
}
