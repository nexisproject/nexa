// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-25, by liasica

package configure

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"nexis.run/nexa/kit"
)

type Configure struct {
	App         string          // 应用名称
	Environment kit.Environment // 环境变量
	Logger      *Logger         // 日志配置
}

type Configurable interface {
	GetApp() string
	GetEnvironment() kit.Environment
	GetLogger() *Logger
}

func (c Configure) GetApp() string {
	return c.App
}

func (c Configure) GetEnvironment() kit.Environment {
	return c.Environment
}

func (c Configure) GetLogger() *Logger {
	return c.Logger
}

type Logger struct {
	Name string // 日志名称

	Stdout bool // 是否输出到控制台

	// 输出至kafka
	Kafka *LoggerKafka
}

type LoggerKafka struct {
	Topic   string
	Brokers []string
}

func (l *Logger) IsVaild() (vaild bool) {
	if l == nil {
		return
	}

	// 如果没有配置kafka和stdout，返回false
	if !l.Stdout && l.Kafka == nil {
		return
	}

	// 如果配置了kafka，topic和name不能为空
	if l.Kafka != nil {
		return l.Kafka.Topic != "" && len(l.Kafka.Brokers) > 0
	}

	return true
}

func Load[T Configurable](p string) (c T, err error) {
	k := koanf.New(".")
	f := file.Provider(p)
	err = k.Load(f, yaml.Parser())
	if err != nil {
		return
	}
	err = k.UnmarshalWithConf(
		"",
		&c,
		koanf.UnmarshalConf{
			Tag: "koanf",
			DecoderConfig: &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToSliceHookFunc(","),
					mapstructure.TextUnmarshallerHookFunc()),
				Metadata:         nil,
				Result:           &c,
				WeaklyTypedInput: true,
				Squash:           true,
			},
		},
	)

	if c.GetApp() == "" {
		err = kit.ErrConfigMissName
	}

	if c.GetEnvironment() == "" || !c.GetEnvironment().IsValid() {
		err = kit.ErrConfigMissEnvironment
	}

	if c.GetLogger() == nil {
		err = kit.ErrConfigMissLogger
	}

	return
}
