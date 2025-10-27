// Copyright (C) nexa. 2025-present.
//
// Created at 2025-10-27, by liasica

package authz

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gopkg.auroraride.com/rbac"

	"nexis.run/nexa/kit/micro"
)

var (
	testAddress       = ":5531"
	allowedPermission = "allowed_permission"
	existingProject   = "existing_project"
	testToken         = "test-token"
)

type testServer struct {
	rbac.UnimplementedRBACServiceServer
}

func (*testServer) GetRestrictedUser(_ context.Context, req *rbac.GetRestrictedUserRequest) (*rbac.GetRestrictedUserResponse, error) {
	hasPermission := req.PermissionKey == allowedPermission
	hasUser := req.ProjectCode == existingProject

	var user *rbac.User
	if hasUser {
		user = &rbac.User{
			Uid: uuid.New().String(),
		}
	}

	res := &rbac.GetRestrictedUserResponse{
		HasPermission:       hasPermission,
		UserInfo:            user,
		DataPermissionRules: nil,
	}
	return res, nil
}

func (*testServer) GetUser(context.Context, *rbac.GetUserRequest) (*rbac.GetUserResponse, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
	return nil, errors.New("not implemented")
}

func TestServer(t *testing.T) {
	micro.Run("test-app", testAddress, func(s *grpc.Server) {
		rbac.RegisterRBACServiceServer(s, &testServer{})
	})

	t.Logf("Test gRPC server running on %s", testAddress)

	select {}
}

func TestGetRestrictedUser(t *testing.T) {
	Setup(testAddress)

	ctx := context.Background()

	res, err := GetRestrictedUser(ctx, testToken, existingProject, allowedPermission)
	require.NoError(t, err)
	require.True(t, res.HasPermission)
	require.NotNil(t, res.UserInfo)

	res, err = GetRestrictedUser(ctx, testToken, existingProject, "forbidden_permission")
	require.NoError(t, err)
	require.False(t, res.HasPermission)
	require.NotNil(t, res.UserInfo)

	res, err = GetRestrictedUser(ctx, testToken, "non_existing_project", allowedPermission)
	require.NoError(t, err)
	require.True(t, res.HasPermission)
	require.Nil(t, res.UserInfo)
}
