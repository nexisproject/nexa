// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-17, by liasica

package shrmdq

import (
	"bytes"
	"encoding/binary"
)

// Montage 协议数据域
type Montage struct {
	buf        *bytes.Buffer
	decomposed [][]byte // 解包后数据域
	unread     [][]byte // 未使用数据域
}

func NewMontage(raw []byte) *Montage {
	return &Montage{
		buf: bytes.NewBuffer(raw),
	}
}

// AddDecomposed 添加解包后数据域
func (m *Montage) AddDecomposed(items ...[]byte) *Montage {
	m.decomposed = append(m.decomposed, items...)
	return m
}

// AddUnread 添加未使用数据域
func (m *Montage) AddUnread(items ...[]byte) *Montage {
	m.unread = append(m.unread, items...)
	return m
}

// GetUnread 获取未使用数据域
func (m *Montage) GetUnread() [][]byte {
	return m.unread
}

// GetDecomposed 获取解包后数据域
func (m *Montage) GetDecomposed() [][]byte {
	return m.decomposed
}

// Decompose 解包数据域
func (m *Montage) Decompose() (*Montage, error) {
	l := make([]byte, 2)
	_, err := m.buf.Read(l)
	if err != nil {
		return m, err
	}

	n := int(binary.LittleEndian.Uint16(l))
	if n == 0 || m.buf.Len() < n {
		// if m.buf.Len() > 0 {
		// 	m.AddDecomposed(m.buf.Bytes())
		// }
		unread := make([]byte, 2+m.buf.Len())
		copy(unread, l)
		_, _ = m.buf.Read(unread[2:])
		m.AddUnread(unread)
		return m, err
	}

	// if m.buf.Len() < n {
	// 	return m, ammeter.ErrorInvalidPacket
	// }

	data := make([]byte, n)
	_, err = m.buf.Read(data)
	if err != nil {
		return m, err
	}

	m.AddDecomposed(data)

	if m.buf.Len() > 0 {
		return m.Decompose()
	}

	return m, nil
}

// Compose 组包
func (m *Montage) Compose() *Montage {
	// 清空缓冲区
	m.buf.Reset()

	// 2字节长度
	l := make([]byte, 2)

	// 按顺序写入数据域
	for i := 0; i < len(m.decomposed); i++ {
		binary.LittleEndian.PutUint16(l, uint16(len(m.decomposed[i])))
		m.buf.Write(l)
		m.buf.Write(m.decomposed[i])
	}

	for i := 0; i < len(m.unread); i++ {
		m.buf.Write(m.unread[i])
	}

	return m
}

// GetComposed 获取组包后数据
func (m *Montage) GetComposed() []byte {
	return m.buf.Bytes()
}
