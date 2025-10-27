// Copyright (C) nexa. 2025-present.
//
// Created at 2025-10-25, by liasica

package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ew "github.com/labstack/echo/v4/middleware"

	"nexis.run/nexa/kit/authz"
)

// RBACMiddleware 权限控制中间件
func RBACMiddleware(skipper ew.Skipper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper != nil && skipper(c) {
				return next(c)
			}

			// 获取用户token
			token := c.Request().Header.Get(HeaderAuthToken)
			projectCode := c.Request().Header.Get(HeaderProjectCode)
			permissionKey := c.Request().Header.Get(HeaderPermissionKey)

			if token != "" && projectCode != "" && permissionKey != "" {
				authed, err := authz.GetRestrictedUser(c.Request().Context(), token, projectCode, permissionKey)
				if err != nil {
					return err
				}

				// 检查权限
				if !authed.HasPermission {
					return WrapError(http.StatusForbidden, authz.ErrForbidden)
				}

				// 检查用户信息
				if authed.UserInfo == nil {
					return WrapError(http.StatusUnauthorized, authz.ErrUnauthorized)
				}

				// 设置用户信息到上下文
				ctx := GetContext(c)
				ctx.User = authed.UserInfo
				c.Set(ContextKeyUser, authed.UserInfo)

				return next(ctx)
			}

			return next(c)
		}
	}
}
