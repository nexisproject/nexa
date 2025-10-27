// Copyright (C) nexa. 2025-present.
//
// Created at 2025-10-25, by liasica

package authz

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gopkg.auroraride.com/rbac"
)

var instance rbac.RBACServiceClient

var _ = Setup

// Setup 初始化 rbac gRPC 客户端
// 如果初始化失败, 会直接抛出致命错误
func Setup(address string) {
	sync.OnceFunc(func() {
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zap.L().Fatal("rbac rpc连接失败", zap.Error(err))
			return
		}
		instance = rbac.NewRBACServiceClient(conn)
	})()
}

// GetRBACContext 取得权限服务上下文并添加认证信息
func GetRBACContext(ctx context.Context, token string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"Authorization": "Bearer " + token}))
}

// GetRestrictedUser 检查权限
func GetRestrictedUser(ctx context.Context, token string, projectCode string, permissionKey string, opts ...Option) (*rbac.GetRestrictedUserResponse, error) {
	o := defaultOption
	for _, opt := range opts {
		opt.apply(o)
	}

	res, err := instance.GetRestrictedUser(GetRBACContext(ctx, token), &rbac.GetRestrictedUserRequest{
		PermissionKey: permissionKey,
		ProjectCode:   projectCode,
	})
	if o.errorHandler != nil {
		err = o.errorHandler(err)
	}

	return res, err
}

var _ = GetUserFromUid

// GetUserFromUid 获取用户信息
func GetUserFromUid(ctx context.Context, uid string, opts ...Option) (*rbac.User, error) {
	o := defaultOption
	for _, opt := range opts {
		opt.apply(o)
	}

	res, err := instance.GetUser(ctx, &rbac.GetUserRequest{
		Uid: uid,
	})
	if o.errorHandler != nil {
		err = o.errorHandler(err)
	}

	if err != nil {
		return nil, err
	}

	return res.UserInfo, nil
}
