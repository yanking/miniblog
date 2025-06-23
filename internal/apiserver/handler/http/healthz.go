package http

import (
	"miniblog/internal/pkg/log"
	apiv1 "miniblog/pkg/api/apiserver/v1"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/onexstack/onexstack/pkg/core"
)

func (h *Handler) Healthz(c *gin.Context) {
	log.W(c.Request.Context()).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")
	core.WriteResponse(c, apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil)
}
