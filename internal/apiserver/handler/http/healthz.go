package http

import (
	"miniblog/internal/pkg/log"
	apiv1 "miniblog/pkg/api/apiserver/v1"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Healthz(c *gin.Context) {
	log.W(c.Request.Context()).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")
	c.JSON(http.StatusOK, apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	})
}
