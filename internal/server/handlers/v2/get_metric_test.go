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

func Test_ParseMetricInfo(t *testing.T) {
	tests := []struct {
		name    string
		reqBody string
		want    models.MetricInfo
		wantErr bool
	}{
		{
			name:    "Валидный json c типом Gauge",
			reqBody: `{"id":"Alloc","type":"gauge","value":10.5}`,
			want:    models.NewGaugeMetric("Alloc", 10.5),
			wantErr: false,
		},
		{
			name:    "Валидный json c типом Counter",
			reqBody: `{"id":"Counter","type":"counter","delta":8}`,
			want:    models.NewCounterMetric("Counter", 8),
			wantErr: false,
		},
		{
			name:    "Невалидный json",
			reqBody: `{"id":"Counter"}`,
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

			got, err := v2handlers.ParseMetricInfo(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMetricInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMetricInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
