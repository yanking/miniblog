package server

import (
	"context"
	"github.com/yanking/miniblog/internal/pkg/log"
	genericoptions "github.com/yanking/miniblog/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

// GRPCServer 代表一个 GRPC 服务器.
type GRPCServer struct {
	srv *grpc.Server
	lis net.Listener
}

// NewGRPCServer 创建一个新的 GRPC 服务器实例.
func NewGRPCServer(
	grpcOptions *genericoptions.GRPCOptions,
	registerServer func(grpc.ServiceRegistrar),
) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", grpcOptions.Addr)
	if err != nil {
		log.Errorw("Failed to listen", "err", err)
		return nil, err
	}

	grpcsrv := grpc.NewServer()

	registerServer(grpcsrv)
	registerHealthServer(grpcsrv)
	reflection.Register(grpcsrv)

	return &GRPCServer{
		srv: grpcsrv,
		lis: lis,
	}, nil
}

// RunOrDie 启动 GRPC 服务器并在出错时记录致命错误.
func (s *GRPCServer) RunOrDie() {
	log.Infow("Start to listening the incoming requests", "protocol", "grpc", "addr", s.lis.Addr().String())
	if err := s.srv.Serve(s.lis); err != nil {
		log.Fatalw("Failed to serve grpc server", "err", err)
	}
}

// GracefulStop 优雅地关闭 GRPC 服务器.
func (s *GRPCServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop grpc server")
	s.srv.GracefulStop()
}

// registerHealthServer 注册健康检查服务.
func registerHealthServer(grpcsrv *grpc.Server) {
	// 创建健康检查服务实例
	healthServer := health.NewServer()

	// 设定服务的健康状态
	healthServer.SetServingStatus("MiniBlog", grpc_health_v1.HealthCheckResponse_SERVING)

	// 注册健康检查服务
	grpc_health_v1.RegisterHealthServer(grpcsrv, healthServer)
}
