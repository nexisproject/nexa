// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package ammeter

type CRCSum interface {
	~uint16 | ~uint32
}

type CRC[T CRCSum] interface {
	CheckSum(data []byte) T
}
