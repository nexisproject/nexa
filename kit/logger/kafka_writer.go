// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-06, by liasica

package logger

import (
	"github.com/segmentio/kafka-go"

	"nexis.run/nexa/pkg/clara"
)

type KafkaWriter struct {
	*clara.Clara
}

func NewKafkaWriter(c *clara.Clara) *KafkaWriter {
	return &KafkaWriter{
		Clara: c,
	}
}

func (w *KafkaWriter) Write(p []byte) (n int, err error) {
	err = w.WriteMessages(kafka.Message{
		Value: p,
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
