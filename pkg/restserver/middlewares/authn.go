package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yanking/miniblog/internal/apiserver/model"
	"github.com/yanking/miniblog/internal/pkg/contextx"
	"github.com/yanking/miniblog/internal/pkg/errno"
	"github.com/yanking/miniblog/internal/pkg/log"
	"github.com/yanking/miniblog/pkg/core"
	"github.com/yanking/miniblog/pkg/token"
)

// UserRetriever 用于根据用户名获取用户的接口.
type UserRetriever interface {
	// GetUser 根据用户ID获取用户信息
	GetUser(ctx context.Context, userID string) (*model.UserM, error)
}

// AuthnMiddleware 是一个认证中间件，用于从 gin.Context 中提取 token 并验证 token 是否合法.
func AuthnMiddleware(retriever UserRetriever) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析 JWT Token
		userID, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, nil, errno.ErrTokenInvalid.WithMessage(err.Error()))
			c.Abort()
			return
		}

		log.Debugw("Token parsing successful", "userID", userID)

		user, err := retriever.GetUser(c, userID)
		if err != nil {
			core.WriteResponse(c, nil, errno.ErrUnauthenticated.WithMessage(err.Error()))
			c.Abort()
			return
		}

		ctx := contextx.WithUserID(c.Request.Context(), user.UserID)
		ctx = contextx.WithUsername(ctx, user.Username)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
