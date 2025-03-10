package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/server/interfaces"
	"github.com/xantinium/metrix/internal/tools"
)

// GetAllMetricHandler реализация хендлера для получения всех метрик в виде HTML.
func GetAllMetricHandler(ctx *gin.Context, s interfaces.Server) (int, string, error) {
	metrics, err := s.GetMetricsRepo().GetAllMetrics(ctx)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	b := strings.Builder{}

	for _, metric := range metrics {
		b.WriteString("<p>")
		b.WriteString("<strong>")
		b.WriteString(metric.ID())
		b.WriteString(": </strong>")
		b.WriteString("<span>")
		switch metric.Type() {
		case models.Gauge:
			b.WriteString(tools.FloatToStr(metric.GaugeValue()))
		case models.Counter:
			b.WriteString(tools.IntToStr(metric.CounterValue()))
		}
		b.WriteString(" (")
		b.WriteString(string(metric.Type()))
		b.WriteString(")</span></p>")
	}

	return http.StatusOK, b.String(), nil
}
