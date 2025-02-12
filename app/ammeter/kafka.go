// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-21, by liasica

package ammeter

import (
	"github.com/json-iterator/go"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"orba.plus/nexa/pkg/clara"
)

const (
	KafkaTopic string = "ammeter"
)

type VoltageData struct {
	A float64 `json:"a"` // A相电压
	B float64 `json:"b"` // B相电压
	C float64 `json:"c"` // C相电压
}

type ElectricData struct {
	A float64 `json:"a"` // A相电流
	B float64 `json:"b"` // B相电流
	C float64 `json:"c"` // C相电流
}

type KafkaMessage struct {
	No string `json:"no"` // 表号

	Voltage    *VoltageData  `json:"voltage,omitempty"`    // 电压数据
	Electric   *ElectricData `json:"electric,omitempty"`   // 电流数据
	PowerTotal *float64      `json:"powerTotal,omitempty"` // (当前)正向有功总电能
}

var c *clara.Clara

func NewKafka(addresses []string) {
	c = clara.New(addresses, clara.WithTopic(KafkaTopic))
}

func SendKafkaMessage(m *KafkaMessage) {
	b, _ := jsoniter.Marshal(m)
	err := c.WriteMessages(kafka.Message{
		Value: b,
	})

	if err != nil {
		zap.L().Error("kafka消息发送失败", zap.Error(err))
	}
}
