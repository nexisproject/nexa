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

func (c *Clara) NewWriter() *kafka.Writer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(c.addresses...),
		Balancer:               &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		RequiredAcks:           kafka.RequireAll,    // ack模式
		Async:                  true,                // 异步
		AllowAutoTopicCreation: true,                // 自动创建topic
	}
	if c.topic != "" {
		w.Topic = c.topic
	}
	return w
}

func (c *Clara) WriteMessages(messages ...kafka.Message) (err error) {
	w := c.NewWriter()
	defer func(w *kafka.Writer) {
		_ = w.Close()
	}(w)

	for i := 0; i < c.retries; i++ {
		err = c.writeMessagesWithTimeout(w, messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, kafka.UnknownTopicOrPartition) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(c.sleep)
			continue
		}

		return
	}
	return
}

func (c *Clara) writeMessagesWithTimeout(w *kafka.Writer, messages ...kafka.Message) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return w.WriteMessages(ctx, messages...)
}
