// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package shrmdq

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"time"
)

type Code byte

// 标识符
// 上行 ^ 0x80 = 下行
// 下行 ^ 0x80 = 上行
// 上行： 服务端发送至设备
// 下行： 设备发送至服务端
const (
	CodeRegisterD Code = 0x01 // 注册下行
	CodeRegisterU Code = 0x81 // 注册上行

	CodeHeartbeatD Code = 0x03 // 心跳下行
	CodeHeartbeatU Code = 0x83 // 心跳上行

	CodeSignalD Code = 0x07 // 查询信号强度下行
	CodeSignalU Code = 0x87 // 查询信号强度上行

	CodePassthroughD Code = 0x0A // 透传数据包下行
	CodePassthroughU Code = 0x8A // 透传数据包上行
)

// Reverse 反转标识符
func (c Code) Reverse() Code {
	return c ^ 0x80
}

// NeedAnswer 是否需要应答
func (c Code) NeedAnswer() bool {
	return c == CodeRegisterU || c == CodeHeartbeatU
}

// Packet 数据包
type Packet struct {
	raw     []byte  // 原始字节组
	decoded bool    // 是否已解码
	message Message // 数据区解析后消息

	Start byte    `json:"start,omitempty"` // 1字节 帧起始符 0x64
	Len   [2]byte `json:"len,omitempty"`   // 2字节 帧总长度
	Keep  byte    `json:"keep,omitempty"`  // 1字节 保留字段
	Sn    [4]byte `json:"sn,omitempty"`    // 4字节 帧序号
	Code  Code    `json:"code,omitempty"`  // 1字节 标识符
	Data  []byte  `json:"data,omitempty"`  // N字节 数据区
	Sum   [2]byte `json:"sum,omitempty"`   // 2字节 CRC16校验
	End   byte    `json:"end,omitempty"`   // 1字节 帧结束符 0x20
}

// NewPacket 创建空数据包
func NewPacket() *Packet {
	return &Packet{}
}

// Answer 应答
// 仅包含保留字段、帧序号、标识符
// 返回应答数据和是否需要应答
func (p *Packet) Answer() (b []byte, need bool) {
	need = p.Code.NeedAnswer()
	if !need {
		return
	}

	if p.Code < 0x80 {
		return
	}

	b = NewPacket().
		SetKeep(p.Keep ^ CommandDefaultKeep).
		SetSn(p.Sn[:]).
		SetCode(p.Code.Reverse()).
		Build()

	return
}

// SetKeep 链式 - 设置保留字段
func (p *Packet) SetKeep(keep byte) *Packet {
	p.Keep = keep
	return p
}

// SetSn 链式 - 设置帧序号
func (p *Packet) SetSn(params ...[]byte) *Packet {
	if len(params) == 0 {
		sn := rand.New(rand.NewSource(time.Now().UnixNano())).Uint32()
		binary.LittleEndian.PutUint32(p.Sn[:], sn)
	} else {
		copy(p.Sn[:], params[0])
	}
	return p
}

// SetCode 链式 - 设置标识符
func (p *Packet) SetCode(code Code) *Packet {
	p.Code = code
	return p
}

// SetData 链式 - 设置数据区
func (p *Packet) SetData(data []byte) *Packet {
	p.Data = data
	return p
}

// Compare 比较数据包
func (p *Packet) Compare(other *Packet) bool {
	return p.Start == other.Start && p.Len == other.Len && p.Keep == other.Keep && p.Sn == other.Sn && p.Code == other.Code && bytes.Equal(p.Data, other.Data) && p.Sum == other.Sum && p.End == other.End
}

// Build 编码构建数据包
func (p *Packet) Build() []byte {
	buf := bytes.NewBuffer(nil)
	p.Start = startFrame
	p.End = endFrame

	// 计算长度
	binary.LittleEndian.PutUint16(p.Len[:], uint16(len(p.Data)+12))

	// 写入帧起始符
	buf.WriteByte(p.Start)

	// 写入长度
	buf.Write(p.Len[:])

	// 写入保留字段
	buf.WriteByte(p.Keep)

	// 写入帧序号
	buf.Write(p.Sn[:])

	// 写入标识符
	buf.WriteByte(byte(p.Code))

	// 写入数据区
	buf.Write(p.Data)

	// 计算校验和
	p.Sum = sumBytes(buf.Bytes())
	// 写入校验和
	buf.Write(p.Sum[:])

	// 写入帧结束符
	buf.WriteByte(p.End)

	return buf.Bytes()
}

// Decode 链式 - 解码数据包
func (p *Packet) Decode(raw []byte) (out *Packet, err error) {
	p.raw = raw

	r := bytes.NewReader(raw)
	p.Start, err = r.ReadByte()
	if err != nil {
		return
	}

	_, err = r.Read(p.Len[:])
	if err != nil {
		return
	}

	p.Keep, err = r.ReadByte()
	if err != nil {
		return
	}

	_, err = r.Read(p.Sn[:])
	if err != nil {
		return
	}

	var code byte
	code, err = r.ReadByte()
	if err != nil {
		return
	}
	p.Code = Code(code)

	p.Data = make([]byte, r.Len()-3)
	_, err = r.Read(p.Data)
	if err != nil {
		return
	}

	_, err = r.Read(p.Sum[:])
	if err != nil {
		return
	}

	p.End, err = r.ReadByte()
	if err != nil {
		return
	}

	p.decoded = true
	return p, nil
}

// GetMessage 获取消息
func (p *Packet) GetMessage() (Message, error) {
	mo, err := NewMontage(p.Data).Decompose()
	if err != nil {
		return nil, err
	}

	if p.message == nil {
		switch p.Code {
		case CodeRegisterU:
			p.message = &RegisterMessage{}
		case CodeHeartbeatU:
			p.message = &HeartbeatMessage{}
		case CodeSignalU:
			p.message = &SignalMessage{}
		case CodePassthroughU:
			p.message = &PassthroughMessage{}
		}
	}

	err = p.message.Decode(mo)
	if err != nil {
		return nil, err
	}

	return p.message, nil
}
