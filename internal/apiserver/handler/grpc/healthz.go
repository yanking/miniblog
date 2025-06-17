package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	apiv1 "miniblog/pkg/api/apiserver/v1"
	"time"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(ctx context.Context, rq *emptypb.Empty) (*apiv1.HealthzResponse, error) {
	return &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
