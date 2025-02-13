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
	Name        string
	Environment kit.Environment
	Logger      *Logger
}

type Logger struct {
	Kafka  []string `koanf:"kafka"`  // 输出至kafka地址列表
	Stdout bool     `koanf:"stdout"` // 是否输出到控制台
	Name   string   `koanf:"name"`   // 日志名称
}

type Configurable interface {
	GetName() string
	GetEnvironment() kit.Environment
	GetLogger() *Logger
}

func (c Configure) GetName() string {
	return c.Name
}

func (c Configure) GetEnvironment() kit.Environment {
	return c.Environment
}

func (c Configure) GetLogger() *Logger {
	return c.Logger
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

	if c.GetName() == "" {
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
