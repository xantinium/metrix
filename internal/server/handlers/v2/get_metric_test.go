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

func TestParseGetMetricRequest(t *testing.T) {
	tests := []struct {
		name    string
		reqBody string
		want    v2handlers.GetMetricsRequest
		wantErr bool
	}{
		{
			name:    "Валидный json для типа Gauge",
			reqBody: `{"id":"Alloc","type":"gauge"}`,
			want:    v2handlers.GetMetricsRequest{MetricName: "Alloc", MetricType: models.Gauge},
		},
		{
			name:    "Валидный json для типа Counter",
			reqBody: `{"id":"PollCount","type":"counter"}`,
			want:    v2handlers.GetMetricsRequest{MetricName: "PollCount", MetricType: models.Counter},
		},
		{
			name:    "Невалидный json: пустой id",
			reqBody: `{"id":""}`,
			wantErr: true,
		},
		{
			name:    "Невалидный json: пустой type",
			reqBody: `{"id":"Alloc","type":""}`,
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

			got, err := v2handlers.ParseGetMetricRequest(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGetMetricRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGetMetricRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
