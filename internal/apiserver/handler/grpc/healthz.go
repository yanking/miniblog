package grpc

import (
	"context"
	apiserverv1 "github.com/yanking/miniblog/api/proto/gen/apiserver/v1"
	"github.com/yanking/miniblog/internal/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func (h *Handler) Healthz(ctx context.Context, req *emptypb.Empty) (*apiserverv1.HealthzResponse, error) {
	log.W(ctx).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")
	return &apiserverv1.HealthzResponse{
		Status:    apiserverv1.ServiceStatus_HEALTHY,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
