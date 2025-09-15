// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-04, by liasica

package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"nexis.run/nexa/kit/configure"
)

var (
	kafkaWriter *KafkaWriter
	once        sync.Once
)

func Setup(cfg *configure.Logger) {
	var cores []zapcore.Core

	// 配置级别
	consoleLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	kafkaLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// 配置编码器 - 明确区分控制台和Kafka的编码器
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 判断是否需要输出到控制台
	shouldLogToConsole := cfg.Stdout || (cfg.Kafka == nil)
	if shouldLogToConsole {
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout), // 明确使用控制台输出
			consoleLevel,
		)
		cores = append(cores, consoleCore)
	}

	// 判断是否需要输出到Kafka
	if cfg.Kafka != nil && len(cfg.Kafka.Brokers) > 0 {
		// Kafka输出配置使用JSON格式 - 明确使用不同的配置
		kafkaEncoderConfig := zap.NewProductionEncoderConfig()
		kafkaEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		kafkaEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		kafkaEncoder := zapcore.NewJSONEncoder(kafkaEncoderConfig)

		// 使用单例模式确保只有一个Kafka writer实例
		once.Do(func() {
			kafkaWriter = NewKafkaWriter(cfg.Kafka.Brokers, cfg.Kafka.Topic)
		})

		// 确保Kafka core只处理JSON格式的日志
		kafkaCore := zapcore.NewCore(
			kafkaEncoder,
			zapcore.AddSync(kafkaWriter), // 使用AddSync包装
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
