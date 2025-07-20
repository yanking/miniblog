package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/onexstack/onexstack/pkg/store/where"
	"github.com/yanking/miniblog/internal/apiserver/model"
	"github.com/yanking/miniblog/internal/apiserver/pkg/conversion"
	"github.com/yanking/miniblog/internal/apiserver/store"
	"github.com/yanking/miniblog/internal/pkg/contextx"
	"github.com/yanking/miniblog/internal/pkg/errno"
	"github.com/yanking/miniblog/internal/pkg/known"
	"github.com/yanking/miniblog/internal/pkg/log"
	"github.com/yanking/miniblog/pkg/auth"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"

	apiv1 "github.com/yanking/miniblog/api/proto/gen/apiserver/v1"
)

// UserBiz 定义处理用户请求所需的方法.
type UserBiz interface {
	Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error)
	Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error)
	List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)

	UserExpansion
}

// UserExpansion 定义用户操作的扩展方法.
type UserExpansion interface {
	Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error)
	RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error)
	ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)
}

// userBiz 是 UserBiz 接口的实现.
type userBiz struct {
	store store.IStore
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ UserBiz = (*userBiz)(nil)

func New(store store.IStore) *userBiz {
	return &userBiz{store: store}
}

// Login 实现 UserBiz 接口中的 Login 方法.
func (b *userBiz) Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	// 获取登录用户的所有信息
	whr := where.F("username", rq.GetUsername())
	userM, err := b.store.User().Get(ctx, whr)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	if err := auth.Compare(userM.Password, rq.GetPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}

	// TODO：实现 Token 签发逻辑

	return &apiv1.LoginResponse{Token: "<placeholder>", ExpireAt: timestamppb.New(time.Now().Add(2 * time.Hour))}, nil
}

// RefreshToken 用于刷新用户的身份验证令牌.
// 当用户的令牌即将过期时，可以调用此方法生成一个新的令牌.
func (b *userBiz) RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error) {
	// TODO：实现 Token 签发逻辑
	return &apiv1.RefreshTokenResponse{Token: "<placeholder>", ExpireAt: timestamppb.New(time.Now().Add(2 * time.Hour))}, nil
}

// ChangePassword 实现 UserBiz 接口中的 ChangePassword 方法.
func (b *userBiz) ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	if err := auth.Compare(userM.Password, rq.GetOldPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}

	userM.Password, _ = auth.Encrypt(rq.GetNewPassword())
	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}

	return &apiv1.ChangePasswordResponse{}, nil
}

// Create 实现 UserBiz 接口中的 Create 方法.
func (b *userBiz) Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
	var userM model.UserM
	_ = copier.Copy(&userM, rq)

	if err := b.store.User().Create(ctx, &userM); err != nil {
		return nil, err
	}

	return &apiv1.CreateUserResponse{UserId: userM.UserID}, nil
}

// Update 实现 UserBiz 接口中的 Update 方法.
func (b *userBiz) Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	if rq.Username != nil {
		userM.Username = rq.GetUsername()
	}
	if rq.Email != nil {
		userM.Email = rq.GetEmail()
	}
	if rq.Nickname != nil {
		userM.Nickname = rq.GetNickname()
	}
	if rq.Phone != nil {
		userM.Phone = rq.GetPhone()
	}

	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}

	return &apiv1.UpdateUserResponse{}, nil
}

// Delete 实现 UserBiz 接口中的 Delete 方法.
func (b *userBiz) Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	// 只有 `root` 用户可以删除用户，并且可以删除其他用户
	// 所以这里不用 where.T()，因为 where.T() 会查询 `root` 用户自己
	if err := b.store.User().Delete(ctx, where.F("userID", rq.UserId)); err != nil {
		return nil, err
	}

	return &apiv1.DeleteUserResponse{}, nil
}

// Get 实现 UserBiz 接口中的 Get 方法.
func (b *userBiz) Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	return &apiv1.GetUserResponse{User: conversion.UserModelToUserV1(userM)}, nil
}

// List 实现 UserBiz 接口中的 List 方法.
func (b *userBiz) List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))
	if contextx.Username(ctx) != known.AdminUsername {
		whr.T(ctx)
	}

	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)

	// 设置最大并发数量为常量 MaxConcurrency
	eg.SetLimit(known.MaxErrGroupConcurrency)

	// 使用 goroutine 提高接口性能
	for _, user := range userList {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				count, _, err := b.store.Post().List(ctx, where.T(ctx))
				if err != nil {
					return err
				}

				converted := conversion.UserModelToUserV1(user)
				converted.PostCount = count
				m.Store(user.ID, converted)

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.W(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	users := make([]*apiv1.User, 0, len(userList))
	for _, item := range userList {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*apiv1.User))
	}

	log.W(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &apiv1.ListUserResponse{TotalCount: count, Users: users}, nil
}

// ListWithBadPerformance 是性能较差的实现方式（已废弃）.
func (b *userBiz) ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))
	if contextx.Username(ctx) != known.AdminUsername {
		whr.T(ctx)
	}

	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	users := make([]*apiv1.User, 0, len(userList))
	for _, user := range userList {
		count, _, err := b.store.Post().List(ctx, where.T(ctx))
		if err != nil {
			return nil, err
		}

		converted := conversion.UserModelToUserV1(user)
		converted.PostCount = count
		users = append(users, converted)
	}

	log.W(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &apiv1.ListUserResponse{TotalCount: count, Users: users}, nil
}
