package http

import (
	"github.com/gin-gonic/gin"
	apiserverv1 "github.com/yanking/miniblog/api/proto/gen/apiserver/v1"
	"github.com/yanking/miniblog/internal/pkg/log"
	"github.com/yanking/miniblog/pkg/core"
	"time"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(c *gin.Context) {
	log.W(c.Request.Context()).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")

	core.WriteResponse(c, apiserverv1.HealthzResponse{
		Status:    apiserverv1.ServiceStatus_HEALTHY,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil)
}
