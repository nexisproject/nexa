// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-04, by liasica

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"nexis.run/nexa/kit/configure"
	"nexis.run/nexa/pkg/clara"
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
	if len(cfg.Kafka) > 0 {
		prod := zap.NewProductionEncoderConfig()
		prod.EncodeTime = zapcore.ISO8601TimeEncoder

		w := NewKafkaWriter(clara.New(cfg.Kafka, clara.WithTopic(cfg.Name+"-log")))
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

	// 设置全局logger
	zap.ReplaceGlobals(l)
}
