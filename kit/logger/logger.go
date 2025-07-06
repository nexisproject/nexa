// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-04, by liasica

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"nexis.run/nexa/kit/configure"
)

func Setup(cfg *configure.Logger) {
	var cores []zapcore.Core

	// 开发环境输出到控制台
	if cfg.Stdout {
		d, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		cores = append(cores, d.Core())
	}

	// 输出到kafka
	if cfg.Kafka != nil {
		prod := zap.NewProductionEncoderConfig()
		prod.EncodeTime = zapcore.ISO8601TimeEncoder

		w := NewKafkaWriter(cfg.Kafka.Brokers, cfg.Kafka.Topic)
		cores = append(
			cores,
			zapcore.NewCore(
				zapcore.NewJSONEncoder(prod),
				w,
				zapcore.InfoLevel,
			),
		)
	}

	l := zap.New(zapcore.NewTee(cores...), zap.AddCaller())

	// 如果配置了日志名称，则设置日志名称
	if cfg.Name != "" {
		l = l.Named(cfg.Name)
	}

	// 设置全局logger
	zap.ReplaceGlobals(l)
}
