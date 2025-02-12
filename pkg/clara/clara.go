// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-09, by liasica

package clara

import "time"

const (
	DefaultRetries = 3                      // 默认重试次数
	DefaultTimeout = 3 * time.Second        // 默认超时时间
	DefaultSleep   = 250 * time.Millisecond // 默认重试间隔
)

type Clara struct {
	addresses []string      // kafka地址
	topic     string        // topic
	retries   int           // 重试次数
	sleep     time.Duration // 重试间隔
	timeout   time.Duration // 超时时间
}

func New(addresses []string, options ...Option) *Clara {
	c := &Clara{
		addresses: addresses,
		retries:   DefaultRetries,
		sleep:     DefaultSleep,
		timeout:   DefaultTimeout,
	}

	for _, option := range options {
		option.apply(c)
	}

	return c
}
