// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-06, by liasica

package logger

import "orba.plus/nexa/kit"

type Config struct {
	// Name 日志名称
	Name string

	// Environment 环境
	Environment kit.Environment

	kafka []string
}

type Option interface {
	apply(*Config)
}

type optionFunc func(*Config)

func (f optionFunc) apply(l *Config) {
	f(l)
}

func WithKafka(addresses []string) Option {
	return optionFunc(func(l *Config) {
		l.kafka = addresses
	})
}
