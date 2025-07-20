package grpc

import apiserverv1 "github.com/yanking/miniblog/api/proto/gen/apiserver/v1"

type Handler struct {
	apiserverv1.UnimplementedMiniBlogServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}
