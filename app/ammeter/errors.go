// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package ammeter

import "errors"

var (
	// ErrorIncompletePacket 数据包不完整
	ErrorIncompletePacket = errors.New("incomplete packet")

	// ErrorInvalidPacket 数据包无效
	ErrorInvalidPacket = errors.New("数据包长度错误")

	// ErrorInvalidDataLength 数据区长度不足
	ErrorInvalidDataLength = errors.New("数据区长度不足")

	// ErrorNotDecoded 数据包未解码
	ErrorNotDecoded = errors.New("数据包未解码")

	// ErrorInvalidData 数据无效
	ErrorInvalidData = errors.New("数据无效")

	// ErrorUnknownData 未知数据
	ErrorUnknownData = errors.New("未知数据")
)
