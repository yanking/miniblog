package grpc

import (
	"context"
	apiserverv1 "github.com/yanking/miniblog/api/proto/gen/apiserver/v1"
)

// Login 用户登录.
func (h *Handler) Login(ctx context.Context, rq *apiserverv1.LoginRequest) (*apiserverv1.LoginResponse, error) {
	return h.biz.UserV1().Login(ctx, rq)
}

// RefreshToken 刷新令牌.
func (h *Handler) RefreshToken(ctx context.Context, rq *apiserverv1.RefreshTokenRequest) (*apiserverv1.RefreshTokenResponse, error) {
	return h.biz.UserV1().RefreshToken(ctx, rq)
}

// ChangePassword 修改用户密码.
func (h *Handler) ChangePassword(ctx context.Context, rq *apiserverv1.ChangePasswordRequest) (*apiserverv1.ChangePasswordResponse, error) {
	return h.biz.UserV1().ChangePassword(ctx, rq)
}

// CreateUser 创建新用户.
func (h *Handler) CreateUser(ctx context.Context, rq *apiserverv1.CreateUserRequest) (*apiserverv1.CreateUserResponse, error) {
	return h.biz.UserV1().Create(ctx, rq)
}

// UpdateUser 更新用户信息.
func (h *Handler) UpdateUser(ctx context.Context, rq *apiserverv1.UpdateUserRequest) (*apiserverv1.UpdateUserResponse, error) {
	return h.biz.UserV1().Update(ctx, rq)
}

// DeleteUser 删除用户.
func (h *Handler) DeleteUser(ctx context.Context, rq *apiserverv1.DeleteUserRequest) (*apiserverv1.DeleteUserResponse, error) {
	return h.biz.UserV1().Delete(ctx, rq)
}

// GetUser 获取用户信息.
func (h *Handler) GetUser(ctx context.Context, rq *apiserverv1.GetUserRequest) (*apiserverv1.GetUserResponse, error) {
	return h.biz.UserV1().Get(ctx, rq)
}

// ListUser 列出用户.
func (h *Handler) ListUser(ctx context.Context, rq *apiserverv1.ListUserRequest) (*apiserverv1.ListUserResponse, error) {
	return h.biz.UserV1().List(ctx, rq)
}
