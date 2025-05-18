package v2handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/models"
	v2handlers "github.com/xantinium/metrix/internal/server/handlers/v2"
)

func TestParseUpdateMetricsRequest(t *testing.T) {
	tests := []struct {
		name    string
		reqBody string
		want    v2handlers.UpdateMetricsRequest
		wantErr bool
	}{
		{
			name:    "Валидный json",
			reqBody: `[{"id":"Alloc","type":"gauge","value":10.5}]`,
			want:    v2handlers.UpdateMetricsRequest{Metrics: []models.MetricInfo{models.NewGaugeMetric("Alloc", 10.5)}},
		},
		{
			name:    "Невалидный json: отсутствует идентификатор",
			reqBody: `[{"type":"gauge","delta":8}]`,
			wantErr: true,
		},
		{
			name:    "Невалидный json: отсутствует тип",
			reqBody: `[{"id":"Alloc","delta":8}]`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &gin.Context{
				Request: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer([]byte(tt.reqBody))),
				},
			}

			got, err := v2handlers.ParseUpdateMetricsRequest(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUpdateMetricsRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseUpdateMetricsRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
