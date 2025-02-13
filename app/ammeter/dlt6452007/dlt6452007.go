// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package dlt6452007

import (
	"bytes"
	"io"

	"nexis.run/nexa/pkg/pool"
)

var (
	suffix = []byte{0xFE, 0xFE, 0xFE, 0xFE}
)

type DLT6452007 struct {
	Address string  `json:"address"` // 地址域
	Control Control `json:"control"` // 控制码
	Data    Data    `json:"data"`    // 数据域
}

// Parse [DL/T 645 - 2007] 协议基本结构解析
// https://www.toky.com.cn/up_pic/2020_12_15_12243_142130.pdf
// https://blog.csdn.net/weixin_44451022/article/details/130793888
// https://zhuanlan.zhihu.com/p/630182168
// 前导符 4字节 0xFE, 0xFE, 0xFE, 0xFE
// 帧起始符 1字节 0x68
// 地址域BCD码(A0 - A5) 6字节, 在DLT645协议中规定，表号字段，数据字段都是逆序的，也就是与实际表号循序相反
// 帧起始符 1字节 0x68
// 控制码(C) 1字节
// 数据长度域(L) 读数据时 L≤200，写数据时 L≤50，L=0 表示无数据域
// 数据域(DATA) 0~200字节
// 校验码(CS) 1字节
// 结束符 1字节 0x16
func Parse(raw []byte) (out *DLT6452007, err error) {
	out = &DLT6452007{}

	r := bytes.NewReader(raw)
	_, _ = r.Seek(5, io.SeekCurrent)

	addr := make([]byte, 6)
	_, _ = r.Read(addr)
	out.Address = AddressFromBytes(addr)

	control := make([]byte, 1)
	_, _ = r.Seek(1, io.SeekCurrent)
	_, _ = r.Read(control)
	out.Control = ControlFromByte(control[0])

	l := make([]byte, 1)
	_, _ = r.Read(l)

	data := make([]byte, ParseLength(l))
	_, _ = r.Read(data)

	out.Data, err = ParseData(data)
	return
}

// ParseLength 解析长度
func ParseLength(b []byte) byte {
	return b[0] & 0xFF
}

func (d *DLT6452007) Bytes() []byte {
	w := pool.GetBuffer()
	defer pool.PutBuffer(w)

	// 前导符
	w.Write(suffix)

	// 帧起始符
	w.WriteByte(0x68)

	// 地址域
	w.Write(AddressToBytes(d.Address))

	// 帧起始符
	w.WriteByte(0x68)

	// 控制码
	w.WriteByte(d.Control.Byte())

	// 获取数据字节串
	data := d.Data.Bytes()

	// 数据长度域
	w.WriteByte(byte(len(data)))

	// 数据域
	w.Write(data)

	// 校验码，排除前导符
	w.WriteByte(Checksum(w.Bytes()[4:]))

	// 结束符
	w.WriteByte(0x16)

	return w.Bytes()
}

// Checksum 校验和
// 从第一个帧起始符开始到校验码之前的所有各字节的模 256 的和，即各字节二进制算术和
func Checksum(data []byte) byte {
	count := 0
	for i := 0; i < len(data); i++ {
		count += int(data[i])
	}
	return byte(count & 0xFF)
}

func New(address string, identifier Identifier, value []byte) *DLT6452007 {
	f := CFunctionRead
	if len(value) > 0 {
		f = CFunctionWrite
	}
	return &DLT6452007{
		Address: address,
		Control: Control{
			Function:  f,
			NextFrame: CNextEnd,
			Response:  CResponseSuccess,
			Direction: CDirectionMaster,
		},
		Data: Data{Identifier: identifier, Value: value},
	}
}
