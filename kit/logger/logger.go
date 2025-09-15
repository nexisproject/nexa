// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-04, by liasica

package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"nexis.run/nexa/kit/configure"
)

func Setup(cfg *configure.Logger) {
	var cores []zapcore.Core

	// 配置级别
	consoleLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel) // 控制台可以接收所有日志级别
	kafkaLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)    // Kafka仅接收Info及以上级别

	// 配置编码器
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 判断是否需要输出到控制台
	shouldLogToConsole := cfg.Stdout || (cfg.Kafka == nil)
	if shouldLogToConsole {
		// 控制台输出core
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			consoleLevel,
		)
		cores = append(cores, consoleCore)
	}

	// 判断是否需要输出到Kafka
	if cfg.Kafka != nil && len(cfg.Kafka.Brokers) > 0 {
		// Kafka输出配置使用JSON格式
		kafkaEncoderConfig := zap.NewProductionEncoderConfig()
		kafkaEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		kafkaEncoder := zapcore.NewJSONEncoder(kafkaEncoderConfig)

		// 创建Kafka writer
		w := NewKafkaWriter(cfg.Kafka.Brokers, cfg.Kafka.Topic)

		// Kafka core只处理Info级别及以上的日志
		kafkaCore := zapcore.NewCore(
			kafkaEncoder,
			w,
			kafkaLevel,
		)

		cores = append(cores, kafkaCore)
	}

	// 组合所有cores
	l := zap.New(zapcore.NewTee(cores...), zap.AddCaller())

	// 设置日志名称
	if cfg.Name != "" {
		l = l.Named(cfg.Name)
	}

	// 替换全局logger
	zap.ReplaceGlobals(l)
}
