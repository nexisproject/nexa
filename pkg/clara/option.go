// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-23, by liasica

package clara

import "time"

type Option interface {
	apply(*Clara)
}

type optionFunc func(*Clara)

func (f optionFunc) apply(c *Clara) {
	f(c)
}

func WithRetries(retries int) Option {
	return optionFunc(func(c *Clara) {
		c.retries = retries
	})
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(c *Clara) {
		c.timeout = timeout
	})
}

func WithSleep(sleep time.Duration) Option {
	return optionFunc(func(c *Clara) {
		c.sleep = sleep
	})
}

func WithTopic(topic string) Option {
	return optionFunc(func(c *Clara) {
		c.topic = topic
	})
}
