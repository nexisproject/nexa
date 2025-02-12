// Copyright (C) micros. 2025-present.
//
// Created at 2025-02-10, by liasica

package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 防止静态检查工具误报
var _ Gracefully = (*Server)(nil)
var _ = Run

type Gracefully interface {
	Start()
	Stop(ctx context.Context)
}

type Server struct{}

func (s *Server) Start() {}

func (s *Server) Stop(_ context.Context) {}

// Run 启动服务
func Run(s Gracefully) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	// 启动服务
	go s.Start()

	// 当中断信号发生时，关闭服务器并返回 (20秒超时)
	<-ctx.Done()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	s.Stop(ctx)
}
