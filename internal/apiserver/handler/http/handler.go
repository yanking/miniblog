package http

import "miniblog/internal/apiserver/biz"

type Handler struct {
	biz biz.IBiz
}

func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{
		biz: biz,
	}
}
