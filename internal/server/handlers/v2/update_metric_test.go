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

func TestParseUpdateMetricRequest(t *testing.T) {
	tests := []struct {
		name    string
		reqBody string
		want    v2handlers.UpdateMetricRequest
		wantErr bool
	}{
		{
			name:    "Валидный json c типом Gauge",
			reqBody: `{"id":"Alloc","type":"gauge","value":10.5}`,
			want:    v2handlers.UpdateMetricRequest{MetricID: "Alloc", MetricType: models.Gauge, GaugeValue: 10.5},
		},
		{
			name:    "Валидный json c типом Counter",
			reqBody: `{"id":"PollCounter","type":"counter","delta":8}`,
			want:    v2handlers.UpdateMetricRequest{MetricID: "PollCounter", MetricType: models.Counter, CounterValue: 8},
		},
		{
			name:    "Невалидный json: отсутствует значение",
			reqBody: `{"id":"Alloc","type":"gauge","delta":8}`,
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

			got, err := v2handlers.ParseUpdateMetricRequest(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUpdateMetricRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseUpdateMetricRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
