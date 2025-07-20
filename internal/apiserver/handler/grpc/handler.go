package grpc

import (
	apiserverv1 "github.com/yanking/miniblog/api/proto/gen/apiserver/v1"
	"github.com/yanking/miniblog/internal/apiserver/biz"
)

type Handler struct {
	apiserverv1.UnimplementedMiniBlogServiceServer
	biz biz.IBiz
}

func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{biz: biz}
}
