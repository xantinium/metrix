package metrics

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/presentation/rpc/gen"
)

func (server *MetricsServer) UpdateMetric(ctx context.Context, rawReq *gen.UpdateMetricRequest) (*gen.UpdateMetricReply, error) {
	req, err := parseUpdateMetricRequest(rawReq)
	if err != nil {
		return nil, err
	}

	var updatedMetric models.MetricInfo

	switch req.Metric.Type() {
	case models.Gauge:
		updatedMetric, err = server.updateGaugeMetric(ctx, req.Metric.ID(), req.Metric.GaugeValue())
	case models.Counter:
		updatedMetric, err = server.updateCounterMetric(ctx, req.Metric.ID(), req.Metric.CounterValue())
	default:
		// Попасть сюда невозможно, из-за валидации запроса.
		err = fmt.Errorf("unknown metric type %q", req.Metric.Type())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.UpdateMetricReply{
		Metric: metricToGen(updatedMetric),
	}, nil
}

func (server *MetricsServer) updateGaugeMetric(ctx context.Context, id string, value float64) (models.MetricInfo, error) {
	updatedValue, err := server.repo.UpdateGaugeMetric(ctx, id, value)
	if err != nil {
		return models.MetricInfo{}, err
	}

	return models.NewGaugeMetric(id, updatedValue), nil
}

func (server *MetricsServer) updateCounterMetric(ctx context.Context, id string, value int64) (models.MetricInfo, error) {
	updatedValue, err := server.repo.UpdateCounterMetric(ctx, id, value)
	if err != nil {
		return models.MetricInfo{}, err
	}

	return models.NewCounterMetric(id, updatedValue), nil
}

type updateMetricRequest struct {
	Metric models.MetricInfo
}

func parseUpdateMetricRequest(rawReq *gen.UpdateMetricRequest) (updateMetricRequest, error) {
	var (
		err error
		req updateMetricRequest
	)

	req.Metric, err = parseMetric(rawReq.Metric)
	if err != nil {
		return req, err
	}

	return req, nil
}
