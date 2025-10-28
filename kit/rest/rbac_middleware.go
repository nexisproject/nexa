// Copyright (C) nexa. 2025-present.
//
// Created at 2025-10-25, by liasica

package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ew "github.com/labstack/echo/v4/middleware"
	"gopkg.auroraride.com/rbac"

	"nexis.run/nexa/kit/authz"
)

// RBACMiddlewareConfig 权限控制中间件配置
type RBACMiddlewareConfig struct {
	EnableRemoteAuth bool       // 是否启用远程权限验证
	StaticUser       *rbac.User // 静态用户信息（当不使用远程验证时）
	Skipper          ew.Skipper // 跳过函数
}

type RBACMiddlewareOption func(*RBACMiddlewareConfig)

var _ = WithRBACRemoteAuth

// WithRBACRemoteAuth 设置是否启用远程权限验证
func WithRBACRemoteAuth(enable bool) RBACMiddlewareOption {
	return func(cfg *RBACMiddlewareConfig) {
		cfg.EnableRemoteAuth = enable
	}
}

var _ = WithRBACStaticUser

// WithRBACStaticUser 设置静态用户信息
func WithRBACStaticUser(user *rbac.User) RBACMiddlewareOption {
	return func(cfg *RBACMiddlewareConfig) {
		cfg.StaticUser = user
	}
}

var _ = WithRBACSkipper

// WithRBACSkipper 设置跳过函数
func WithRBACSkipper(skipper ew.Skipper) RBACMiddlewareOption {
	return func(cfg *RBACMiddlewareConfig) {
		cfg.Skipper = skipper
	}
}

var _ = RBACMiddleware

// RBACMiddleware 权限控制中间件
func RBACMiddleware(opts ...RBACMiddlewareOption) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cfg := &RBACMiddlewareConfig{
				EnableRemoteAuth: true,
			}

			for _, opt := range opts {
				opt(cfg)
			}

			ctx := GetContext(c)

			// 是否跳过权限检查
			skip := cfg.Skipper != nil && cfg.Skipper(c)

			// 获取用户token
			token := c.Request().Header.Get(HeaderAuthToken)
			projectCode := c.Request().Header.Get(HeaderProjectCode)
			permissionKey := c.Request().Header.Get(HeaderPermissionKey)

			var (
				user          *rbac.User
				hasPermission bool
			)

			// 获取用户信息和权限
			if cfg.EnableRemoteAuth && token != "" && projectCode != "" && permissionKey != "" {
				authed, err := authz.GetRestrictedUser(c.Request().Context(), token, projectCode, permissionKey)
				if err != nil {
					return err
				}

				user = authed.UserInfo
				hasPermission = authed.HasPermission
			}

			// 如果未使用远程验证且配置了静态用户信息
			if !cfg.EnableRemoteAuth && cfg.StaticUser != nil {
				user = cfg.StaticUser
				hasPermission = true
			}

			// 如果用户信息不为空
			if user != nil {
				// 设置用户信息到上下文
				ctx.User = user
				c.Set(ContextKeyUser, user)
			}

			// 检查用户信息是否跳过
			if !skip && user == nil {
				return WrapError(http.StatusUnauthorized, authz.ErrUnauthorized)
			}

			// 检查权限是否跳过
			if !skip && !hasPermission {
				return WrapError(http.StatusForbidden, authz.ErrForbidden)
			}

			return next(ctx)
		}
	}
}
