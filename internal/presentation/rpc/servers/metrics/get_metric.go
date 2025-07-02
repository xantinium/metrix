package metrics

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/presentation/rpc/gen"
)

func (server *MetricsServer) GetMetric(ctx context.Context, rawReq *gen.GetMetricRequest) (*gen.GetMetricReply, error) {
	req, err := parseGetMetricRequest(rawReq)
	if err != nil {
		return nil, err
	}

	switch req.MetricType {
	case models.Gauge:
		return server.getGaugeMetric(ctx, req.MetricID)
	case models.Counter:
		return server.getCounterMetric(ctx, req.MetricID)
	default:
		// Попасть сюда невозможно, из-за валидации запроса.
		return nil, status.Errorf(codes.Internal, "unknown metric type %q", req.MetricType)
	}
}

func (server *MetricsServer) getGaugeMetric(ctx context.Context, id string) (*gen.GetMetricReply, error) {
	value, err := server.repo.GetGaugeMetric(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetMetricReply{
		Metric: &gen.Metric{
			Id:    id,
			MType: gen.MetricType_GAUGE,
			Value: value,
		},
	}, nil
}

func (server *MetricsServer) getCounterMetric(ctx context.Context, id string) (*gen.GetMetricReply, error) {
	value, err := server.repo.GetCounterMetric(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetMetricReply{
		Metric: &gen.Metric{
			Id:    id,
			MType: gen.MetricType_COUNTER,
			Delta: value,
		},
	}, nil
}

type getMetricRequest struct {
	MetricID   string
	MetricType models.MetricType
}

func parseGetMetricRequest(rawReq *gen.GetMetricRequest) (getMetricRequest, error) {
	var (
		err error
		req getMetricRequest
	)

	req.MetricID = rawReq.GetMetricId()
	if req.MetricID == "" {
		return req, status.Error(codes.InvalidArgument, "metric id cannot be empty")
	}

	req.MetricType, err = parseMetricType(rawReq.GetMetricType())
	if err != nil {
		return req, err
	}

	return req, nil
}
