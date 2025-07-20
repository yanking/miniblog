package http

import (
	"github.com/gin-gonic/gin"
	apiserverv1 "github.com/yanking/miniblog/api/proto/gen/apiserver/v1"
	"time"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(c *gin.Context) {
	// 返回 JSON 响应
	c.JSON(200, &apiserverv1.HealthzResponse{
		Status:    apiserverv1.ServiceStatus_HEALTHY,
		Timestamp: time.Now().Format(time.DateTime),
	})
}
