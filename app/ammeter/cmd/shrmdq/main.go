// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-10, by liasica

package main

import (
	"flag"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"

	"orba.plus/nexa/app/ammeter"
	"orba.plus/nexa/app/ammeter/shrmdq"
)

func main() {
	var cfg string
	flag.StringVar(&cfg, "config", "config/config.yaml", "配置文件")
	flag.Parse()

	ammeter.Boot(cfg)

	l, err := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	if err != nil {
		panic(err)
	}

	err = gnet.Run(
		ammeter.NewHandler(&shrmdq.Codec{}, shrmdq.New()),
		"tcp://"+ammeter.GetConfig().Server.Tcp.Bind,
		gnet.WithMulticore(true),
		gnet.WithReuseAddr(true),
	)
	if err != nil {
		zap.L().Fatal("启动失败", zap.Error(err))
	}
}
