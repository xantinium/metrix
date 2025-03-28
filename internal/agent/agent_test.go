// Пакет agent содержит реализацию агента для сбора метрик.
package agent

import (
	"testing"

	"github.com/xantinium/metrix/internal/models"
)

func TestMetrixAgent_GetUpdateMetricHandlerUrl(t *testing.T) {
	tests := []struct {
		name       string
		serverAddr string
		metric     models.MetricInfo
		want       string
	}{
		{
			name:       "Создание URL-адреса для метрики типа Gauge",
			serverAddr: "localhost:8080",
			metric:     models.NewGaugeMetric("Alloc", 123.45),
			want:       "http://localhost:8080/update/gauge/Alloc/123.45/",
		},
		{
			name:       "Создание URL-адреса для метрики типа Counter",
			serverAddr: "localhost:8080",
			metric:     models.NewCounterMetric("PollCount", 7),
			want:       "http://localhost:8080/update/counter/PollCount/7/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := NewMetrixAgent(MetrixAgentOptions{ServerAddr: tt.serverAddr})
			if got := agent.getUpdateMetricHandlerURL(tt.metric); got != tt.want {
				t.Errorf("MetrixAgent.getUpdateMetricHandlerURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
