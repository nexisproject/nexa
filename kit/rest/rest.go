// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-04, by liasica

package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// 防止静态检查工具误报
var _ = Run

type RouteHandler func(e *echo.Echo)

// Run 启动Rest服务
func Run(app, address string, r RouteHandler) (e *echo.Echo, ch chan error) {
	e = echo.New()

	// 隐藏banner
	e.HideBanner = true

	// 隐藏打印端口
	e.HidePort = true

	// 获取真实IP
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	// 默认json序列化工具
	e.JSONSerializer = NewDefaultJSONSerializer()

	// 绑定校验器
	e.Validator = NewValidator()

	// 默认错误处理
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		_ = GetContext(c).SendResponse(http.StatusInternalServerError, err)
	}

	// 未找到错误
	echo.NotFoundHandler = func(c echo.Context) error {
		return GetContext(c).SendResponse(http.StatusNotFound)
	}

	// 请求方式错误
	echo.MethodNotAllowedHandler = func(c echo.Context) error {
		routerAllowMethods, ok := c.Get(echo.ContextKeyHeaderAllow).(string)
		if ok && routerAllowMethods != "" {
			c.Response().Header().Set(echo.HeaderAllow, routerAllowMethods)
		}
		return GetContext(c).SendResponse(http.StatusMethodNotAllowed)
	}

	// 设置全局中间件
	e.Use(
		ContextMiddleware(app),
		RecoverMiddleware(),
	)

	// 设置路由
	r(e)

	// 使用协程启动HTTP Rest服务器
	ch = make(chan error, 1)
	go func() {
		if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ch <- fmt.Errorf("HTTP Rest 服务启动失败: %w", err)
		}
	}()

	zap.L().Info(fmt.Sprintf("⇨ HTTP Rest server listening on %s", address))

	return
}
