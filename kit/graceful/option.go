// Copyright (C) nexa. 2025-present.
//
// Created at 2025-07-08, by liasica

package graceful

import "time"

type option struct {
	timeout time.Duration // 优雅停止超时时间, 若小于等于0则不设置超时
}

type Option interface {
	apply(*option)
}

type optionFunc func(*option)

func (f optionFunc) apply(c *option) {
	f(c)
}

func WithTimeout(d time.Duration) Option {
	return optionFunc(func(c *option) {
		c.timeout = d
	})
}
