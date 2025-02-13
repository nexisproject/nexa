// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package shrmdq

import (
	"encoding/binary"
	"log"

	"github.com/panjf2000/gnet/v2"

	"nexis.run/nexa/app/ammeter"
	"nexis.run/nexa/pkg/dump"
)

// 数据包定义
// 帧起始符 + 帧总长度 + 保留字段 + 帧序号 + 标识符 + 数据区 + CRC16校验和 + 帧结束符
// 帧头 byte 1 固定为0x64
// 帧总长度 uint16 2 小端序

// 帧定义
const (
	startFrame = 0x64
	endFrame   = 0x20
)

const (
	headerSize = 3 // 帧头+帧长度
)

type Codec struct{}

func (codec *Codec) Decode(conn gnet.Conn) ([]byte, error) {
	buf, _ := conn.Peek(headerSize)
	if len(buf) < headerSize {
		return nil, ammeter.ErrorIncompletePacket
	}

	// 获取数据包总长度
	msgLen := int(binary.LittleEndian.Uint16(buf[1:headerSize]))
	if conn.InboundBuffered() < msgLen {
		return nil, ammeter.ErrorIncompletePacket
	}

	// 获取数据包
	buf, _ = conn.Peek(msgLen)
	_, _ = conn.Discard(msgLen)

	log.Printf(">>>> N: %d, DATA: %s\n", msgLen, dump.Bytes(buf))
	return buf, nil
}

// Encode 编码，生成数据包
// 数据包格式：保留字段 + 帧序号 + 标识符 + 数据区
// 最终返回完整数据包
func (codec *Codec) Encode(raw []byte) []byte {
	// size := len(raw) + 6
	// b = make([]byte, size)
	//
	// // 放入 帧起始符
	// b[0] = startFrame
	//
	// // 放入 帧总长度
	// binary.LittleEndian.PutUint16(b[1:3], uint16(size))
	//
	// // 放入 保留字段 + 帧序号 + 标识符 + 数据区
	// copy(b[3:size-3], raw)
	//
	// // 放入 校验和
	// sum := checkSum(b[:size-3])
	// binary.LittleEndian.PutUint16(b[size-3:size-1], sum)
	//
	// // 放入 结束符
	// b[size-1] = endFrame

	return raw
}
