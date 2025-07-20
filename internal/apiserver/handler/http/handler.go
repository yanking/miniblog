package http

import "github.com/yanking/miniblog/internal/apiserver/biz"

// Handler 处理博客模块的请求.
type Handler struct {
	biz biz.IBiz
}

// NewHandler 创建新的 Handler 实例.
func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{biz: biz}
}
