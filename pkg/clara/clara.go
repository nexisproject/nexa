// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-09, by liasica

package clara

import (
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/segmentio/kafka-go"
)

const (
	DefaultRetries = 3                      // 默认重试次数
	DefaultTimeout = 3 * time.Second        // 默认超时时间
	DefaultSleep   = 250 * time.Millisecond // 默认重试间隔
)

type Clara struct {
	addresses []string
	writers   cmap.ConcurrentMap[string, *kafka.Writer]
}

func New(addresses []string) *Clara {
	c := &Clara{
		addresses: addresses,
		writers:   cmap.New[*kafka.Writer](),
	}

	return c
}
