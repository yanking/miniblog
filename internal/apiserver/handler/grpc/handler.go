package grpc

import (
	"miniblog/internal/apiserver/biz"
	apiv1 "miniblog/pkg/api/apiserver/v1"
)

type Handler struct {
	apiv1.UnimplementedMiniBlogServer

	biz biz.IBiz
}

func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{
		biz: biz,
	}
}
