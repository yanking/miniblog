package grpc

import (
	"context"
	"miniblog/internal/apiserver/model"
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/errno"
	"miniblog/internal/pkg/known"
	"miniblog/internal/pkg/log"
	"miniblog/pkg/token"

	"google.golang.org/grpc"
)

// UserRetriever 用于根据用户名获取用户信息的接口.
type UserRetriever interface {
	// GetUser 根据用户 ID 获取用户信息
	GetUser(ctx context.Context, userID string) (*model.UserM, error)
}

// AuthnInterceptor 是一个 gRPC 拦截器，用于进行认证.
func AuthnInterceptor(retriever UserRetriever) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 解析 JWT Token
		userID, err := token.ParseRequest(ctx)
		if err != nil {
			log.Errorw("Failed to parse request", "err", err)
			return nil, errno.ErrTokenInvalid.WithMessage(err.Error())
		}

		log.Debugw("Token parsing successful", "userID", userID)

		user, err := retriever.GetUser(ctx, userID)
		if err != nil {
			return nil, errno.ErrUnauthenticated.WithMessage(err.Error())
		}

		// 将用户信息存入上下文
		//nolint: staticcheck
		ctx = context.WithValue(ctx, known.XUsername, user.Username)
		//nolint: staticcheck
		ctx = context.WithValue(ctx, known.XUserID, userID)

		// 供 log 和 contextx 使用
		ctx = contextx.WithUserID(ctx, user.UserID)
		ctx = contextx.WithUsername(ctx, user.Username)

		// 继续处理请求
		return handler(ctx, req)
	}
}
