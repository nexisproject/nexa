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

func (c *Clara) NewWriter(topic string) *kafka.Writer {
	w, exists := c.writers.Get(topic)
	if exists {
		return w
	}

	w = &kafka.Writer{
		Addr:                   kafka.TCP(c.addresses...),
		RequiredAcks:           kafka.RequireAll, // ack模式
		Topic:                  topic,
		Async:                  true, // 异步
		AllowAutoTopicCreation: true, // 自动创建topic
		// Balancer:               &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
	}
	c.writers.Set(topic, w)

	return w
}

func (c *Clara) WriteMessages(topic string, messages ...kafka.Message) (err error) {
	w := c.NewWriter(topic)

	for i := 0; i < DefaultRetries; i++ {
		err = c.writeMessagesWithTimeout(w, messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, kafka.UnknownTopicOrPartition) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(DefaultSleep)
			continue
		}

		return
	}
	return
}

func (c *Clara) writeMessagesWithTimeout(w *kafka.Writer, messages ...kafka.Message) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return w.WriteMessages(ctx, messages...)
}
