package serverinterceptors

import (
	"context"
	"github.com/onexstack/onexstack/pkg/log"
	"google.golang.org/grpc"
	"runtime/debug"
)

func StreamCrashInterceptor(svr interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v\n \n %s", r, debug.Stack())
	})

	return handler(svr, stream)
}

func UnaryCrashInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("%+v\n \n %s", r, debug.Stack())
	})

	return handler(ctx, req)
}

func handleCrash(hanlder func(interface{})) {
	if r := recover(); r != nil {
		hanlder(r)
	}
}
