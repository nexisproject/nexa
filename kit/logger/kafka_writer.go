// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-06, by liasica

package logger

import (
	"context"

	"github.com/segmentio/kafka-go"

	"nexis.run/nexa/pkg/clara"
)

type KafkaWriter struct {
	*clara.Writer
}

func NewKafkaWriter(brokers []string, topic string) *KafkaWriter {
	return &KafkaWriter{
		Writer: clara.NewWriter(brokers, topic),
	}
}

func (w *KafkaWriter) Write(p []byte) (n int, err error) {
	// 创建一个副本以避免数据竞争
	safeCopy := make([]byte, len(p))
	copy(safeCopy, p)

	err = w.SendMessages(context.Background(), kafka.Message{
		Value: safeCopy,
	})
	if err != nil {
		return
	}
	n = len(p)
	return
}

func (w *KafkaWriter) Sync() error {
	return nil
}
