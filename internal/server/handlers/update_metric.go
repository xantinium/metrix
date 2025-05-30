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
// @Tags Metrics_Legacy
// @Summary Запрос на обновление метрик
// @Description Запрос на обновление метрик
// @ID updateMetricLegacy
// @Produce text/plain
// @Param metric_type path string true "Тип метрики"
// @Param metric_id path string true "Идентификатор метрики"
// @Param metric_value path string true "Значение метрики"
// @Success 200 {string} string
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /update/{metric_type}/{metric_id}/{metric_value} [post]
func UpdateMetricHandler(ctx *gin.Context, s interfaces.Server) (int, string, error) {
	req, err := parseUpdateMetricRequest(ctx)
	if err != nil {
		return http.StatusBadRequest, "", err
	}

	metricsRepo := s.GetMetricsRepo()

	switch req.metricType {
	case models.Gauge:
		_, err = metricsRepo.UpdateGaugeMetric(ctx, req.metricID, req.metricValue)
	case models.Counter:
		_, err = metricsRepo.UpdateCounterMetric(ctx, req.metricID, int64(req.metricValue))
	}

	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, "", nil
}

// updateMetricRequest структура запроса обновления метрик.
type updateMetricRequest struct {
	metricType  models.MetricType
	metricID    string
	metricValue float64
}

// parseUpdateMetricRequest парсит сырой HTTP-запрос в структуру запроса.
func parseUpdateMetricRequest(r *gin.Context) (updateMetricRequest, error) {
	var (
		err                                         error
		maybeMetricType, metricID, maybeMetricValue string
		metricType                                  models.MetricType
		metricValue                                 float64
	)

	maybeMetricType = r.Param("type")
	if maybeMetricType == "" {
		return updateMetricRequest{}, fmt.Errorf("metric type is missing")
	}

	metricType, err = models.ParseStringAsMetricType(maybeMetricType)
	if err != nil {
		return updateMetricRequest{}, err
	}

	metricID = r.Param("id")
	if metricID == "" {
		return updateMetricRequest{}, fmt.Errorf("metric id is missing")
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
		metricID:    metricID,
		metricValue: metricValue,
	}, nil
}
