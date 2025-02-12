// Copyright (C) micros. 2025-present.
//
// Created at 2025-02-10, by liasica

package micro

import (
	"context"
	"fmt"
	"io"

	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go.uber.org/zap"
)

func init() {
	// 关闭kratos日志
	kratoslog.SetLogger(kratoslog.NewStdLogger(io.Discard))
}

// 防止静态检查工具误报
var _ = Run

type Handler func(s *grpc.Server)

func Run(name, address string, h Handler) (server *grpc.Server, ch chan error) {

	ctx := &Context{
		Name:    name,
		Context: context.WithValue(context.Background(), "name", name),
	}

	server = grpc.NewServer(
		grpc.Address(address),
		grpc.Middleware(
			RecoverMiddleware(),
		),
	)

	h(server)

	// 使用协程启动gRPC服务器
	ch = make(chan error, 1)
	go func() {
		if err := server.Start(ctx); err != nil {
			ch <- fmt.Errorf("gRPC 服务启动失败: %w", err)
		}
	}()

	zap.L().Info(fmt.Sprintf("⇨ gRPC server listening on %s", address))

	return
}
