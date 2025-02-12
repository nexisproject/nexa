// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-22, by liasica

package ammeter

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server struct {
		Tcp struct {
			Bind string
		}
	}

	Kafka struct {
		Addresses []string
	}
}

var (
	cfg *Config
	k   = koanf.New(".")
)

func LoadConfig(p string) {
	cfg = new(Config)
	err := k.Load(file.Provider(p), yaml.Parser())
	if err != nil {
		panic(err)
	}

	err = k.UnmarshalWithConf(
		"",
		cfg,
		koanf.UnmarshalConf{
			Tag: "koanf",
			DecoderConfig: &mapstructure.DecoderConfig{
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToSliceHookFunc(","),
					mapstructure.TextUnmarshallerHookFunc()),
				Metadata:         nil,
				Result:           cfg,
				WeaklyTypedInput: true,
				Squash:           true,
			},
		},
	)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return cfg
}
