// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-23, by liasica

package ammeter

import (
	"os"
	"time"
)

func Boot(p string) {
	// 设置全局时区
	tz := "Asia/Shanghai"
	_ = os.Setenv("TZ", tz)
	loc, _ := time.LoadLocation(tz)
	time.Local = loc

	// 加载配置
	LoadConfig(p)

	// 初始化Kafka
	NewKafka(cfg.Kafka.Addresses)
}
