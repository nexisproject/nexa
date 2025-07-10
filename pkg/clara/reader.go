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

	clara  *Clara
	reader *kafka.Reader
}

type MessageListener func(message kafka.Message, err error) error

var _ = NewReader

// NewReader 创建一个新的Kafka Reader
func NewReader(brokers []string, topic, groupID string) *Reader {
	c := New(brokers)

	return &Reader{
		topic:   topic,
		groupID: groupID,
		clara:   c,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  c.brokers,
			Topic:    topic,
			GroupID:  groupID,
			MaxBytes: 10e6, // 10MB
			// https://github.com/segmentio/kafka-go/issues/800#issuecomment-981855523
			WatchPartitionChanges:  true,
			PartitionWatchInterval: time.Second * 5,
		}),
	}
}

// Listen 监听消息回调
func (r *Reader) Listen(cb MessageListener) error {
	// r.SetOffset(42) // 设置Offset

	// 接收消息
	for {
		err := cb(r.reader.ReadMessage(context.Background()))
		if err != nil {
			return err
		}
	}
}
