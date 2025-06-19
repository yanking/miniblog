package http

import (
	apiv1 "miniblog/pkg/api/apiserver/v1"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	})
}
