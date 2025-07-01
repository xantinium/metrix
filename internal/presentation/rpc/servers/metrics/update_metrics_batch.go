package metrics

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/presentation/rpc/gen"
)

func (server *MetricsServer) UpdateMetricsBatch(ctx context.Context, rawReq *gen.UpdateMetricsBatchRequest) (*empty.Empty, error) {
	req, err := parseUpdateMetricsBatchRequest(rawReq)
	if err != nil {
		return nil, err
	}

	err = server.repo.UpdateMetrics(ctx, req.Metrics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return nil, nil
}

type updateMetricsBatchRequest struct {
	Metrics []models.MetricInfo
}

func parseUpdateMetricsBatchRequest(rawReq *gen.UpdateMetricsBatchRequest) (updateMetricsBatchRequest, error) {
	var (
		err error
		req updateMetricsBatchRequest
	)

	metrics := rawReq.GetMetrics()
	if metrics == nil {
		return req, status.Error(codes.InvalidArgument, "metrics data is missing")
	}

	req.Metrics = make([]models.MetricInfo, len(metrics))

	for i := range metrics {
		req.Metrics[i], err = parseMetric(metrics[i])
		if err != nil {
			return req, err
		}
	}

	return req, nil
}
