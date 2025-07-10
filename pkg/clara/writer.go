// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-09, by liasica

package clara

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	DefaultRetries       = 3                      // 默认重试次数
	DefaultTimeout       = 3 * time.Second        // 默认超时时间
	DefaultRetryInterval = 250 * time.Millisecond // 默认重试间隔
)

type Writer struct {
	*kafka.Writer

	retries       int
	retryInterval time.Duration
	timeout       time.Duration
}

var _ = NewWriter

func NewWriter(brokers []string, topic string, opts ...Option) *Writer {
	c := New(brokers)

	w, exists := c.writers.Get(topic)
	if exists {
		return w
	}

	w = &Writer{
		Writer: &kafka.Writer{
			Addr:                   kafka.TCP(c.brokers...),
			RequiredAcks:           kafka.RequireAll, // ack模式
			Topic:                  topic,
			Async:                  true, // 异步
			AllowAutoTopicCreation: true, // 自动创建topic
			// Balancer:               &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		},
		retries:       DefaultRetries,
		retryInterval: DefaultRetryInterval,
		timeout:       DefaultTimeout,
	}

	for _, opt := range opts {
		opt.apply(w)
	}

	c.writers.Set(topic, w)

	return w
}

func (w *Writer) SendMessages(messages ...kafka.Message) (err error) {
	for i := 0; i < w.retries; i++ {
		err = w.writeMessagesWithTimeout(messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, kafka.UnknownTopicOrPartition) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(DefaultRetryInterval)
			continue
		}

		return
	}
	return
}

func (w *Writer) writeMessagesWithTimeout(messages ...kafka.Message) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return w.WriteMessages(ctx, messages...)
}
