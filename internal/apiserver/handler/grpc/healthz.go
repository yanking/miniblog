package grpc

import (
	"context"
	"miniblog/internal/pkg/log"
	apiv1 "miniblog/pkg/api/apiserver/v1"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(ctx context.Context, rq *emptypb.Empty) (*apiv1.HealthzResponse, error) {
	log.W(ctx).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")
	return &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
