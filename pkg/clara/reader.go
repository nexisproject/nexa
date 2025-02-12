// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-09, by liasica

package clara

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	topic   string
	groupID string

	clara *Clara
	*kafka.Reader
}

type MessageListener func(message kafka.Message, err error) error

func (c *Clara) NewReader(topic, groupID string) *Reader {
	return &Reader{
		topic:   topic,
		groupID: groupID,
		clara:   c,
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  c.addresses,
			Topic:    topic,
			GroupID:  groupID,
			MaxBytes: 10e6, // 10MB
			// https://github.com/segmentio/kafka-go/issues/800#issuecomment-981855523
			WatchPartitionChanges:  true,
			PartitionWatchInterval: time.Second * 5,
		}),
	}
}

func (c *Clara) Listen(topic, groupID string, cb MessageListener) error {
	// 创建Reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.addresses,
		Topic:    topic,
		GroupID:  groupID,
		MaxBytes: 10e6, // 10MB
		// https://github.com/segmentio/kafka-go/issues/800#issuecomment-981855523
		WatchPartitionChanges:  true,
		PartitionWatchInterval: time.Second * 5,
	})
	// r.SetOffset(42) // 设置Offset

	// 接收消息
	for {
		err := cb(r.ReadMessage(context.Background()))
		if err != nil {
			return err
		}
	}
}
