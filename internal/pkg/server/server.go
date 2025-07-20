package server

import (
	"context"
	"net/http"
)

// Server 定义所有服务器类型的接口.
type Server interface {
	// RunOrDie 运行服务器，如果运行失败会退出程序（OrDie的含义所在）.
	RunOrDie()
	// GracefulStop 方法用来优雅关停服务器。关停服务器时需要处理 context 的超时时间.
	GracefulStop(ctx context.Context)
}

// protocolName 从 http.Server 中获取协议名.
func protocolName(server *http.Server) string {
	if server.TLSConfig != nil {
		return "https"
	}
	return "http"
}
