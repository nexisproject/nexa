// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package ammeter

import (
	"github.com/panjf2000/gnet/v2"
)

type NetPacket[T any] interface {
	Decode(raw []byte) (*T, error)
	Bytes() []byte
}

type Codec interface {
	// Decode 解码，获得全部数据包
	Decode(conn gnet.Conn) (b []byte, err error)

	// Encode 编码，生成数据包
	Encode(data []byte) (b []byte)
}
