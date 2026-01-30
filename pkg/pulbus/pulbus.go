// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-28, by liasica

package pulbus

import (
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"
	"go.uber.org/zap"
)

type Pulbus struct {
	client pulsar.Client
	admin  *Admin

	producers sync.Map // map[Topic]pulsar.Producer - 缓存 producer，避免重复创建
	consumers sync.Map // map[ConsumerKey]pulsar.Consumer - 缓存 consumer，避免重复创建
}

// Option Pulbus 配置选项
type Option func(bus *Pulbus)

// WithAdmin 配置 Pulsar Admin 客户端
func WithAdmin(webServiceURL string, opts ...AdminOption) Option {
	return func(bus *Pulbus) {
		admin, err := NewAdmin(webServiceURL, opts...)
		if err != nil {
			zap.L().Error("Pulsar Admin 创建失败", zap.String("webServiceURL", webServiceURL), zap.Error(err))
			return
		}

		bus.admin = admin
	}
}

func New(bookie string, opts ...Option) (bus *Pulbus, err error) {
	var client pulsar.Client
	client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL: bookie,
	})

	if err != nil {
		return
	}

	bus = &Pulbus{
		client:    client,
		producers: sync.Map{},
		consumers: sync.Map{},
	}

	// 应用选项
	for _, opt := range opts {
		opt(bus)
	}

	return bus, nil
}

// Close 关闭所有 producers、consumers 和 client
func (bus *Pulbus) Close() error {
	// 关闭所有 producers
	bus.producers.Range(func(key, value interface{}) bool {
		if producer, ok := value.(pulsar.Producer); ok {
			producer.Close()
		}
		return true
	})

	// 关闭所有 consumers
	bus.consumers.Range(func(key, value interface{}) bool {
		if consumer, ok := value.(pulsar.Consumer); ok {
			consumer.Close()
		}
		return true
	})

	// 关闭 client
	bus.client.Close()
	return nil
}

// GetAdmin 获取 Pulsar Admin 客户端
func (bus *Pulbus) GetAdmin() *Admin {
	return bus.admin
}
