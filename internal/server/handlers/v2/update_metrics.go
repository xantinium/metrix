package v2handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/server/interfaces"
)

// UpdateMetricsHandler реализация хендлера для батчевого обновления метрик.
// @Tags Metrics
// @Summary Батчевое обновление метрик
// @Description Батчевое обновление метрик
// @ID updateMetrics
// @Accept  json
// @Produce json
// @Param payload body MetricsBatch true "Тело запроса"
// @Success 200 {object} nil
// @Failure 400 {string} string "Неверный запрос"
// @Failure 404 {string} string "Метрика не найдена"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /updates [post]
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

// UpdateMetricsRequest запрос на батчевое обновление метрик.
type UpdateMetricsRequest struct {
	Metrics []models.MetricInfo
}

// ParseUpdateMetricsRequest парсит запрос на батчевое обновление метрик.
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
		var metricInfo models.MetricInfo

		metricInfo, err = parseMetric(metric)
		if err != nil {
			return UpdateMetricsRequest{}, err
		}

		req.Metrics[i] = metricInfo
	}

	return req, nil
}
