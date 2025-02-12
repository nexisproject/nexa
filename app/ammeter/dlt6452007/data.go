// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package dlt6452007

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"

	"orba.plus/nexa/pkg/pool"
)

type Data struct {
	Identifier Identifier `json:"identifier"` // 数据标识符
	Value      []byte     `json:"value"`      // 数据域
}

// ParseValue 获取数据域的值
// TODO: 组合有功、无功电能最高位是符号位，0正1负。取值范围：0.00～799999.99
func (d Data) ParseValue() (*DataResult, error) {
	switch d.Identifier {
	default:
		return nil, errors.Wrap(ErrorUnknownIdentifier, fmt.Sprintf("%d", d.Identifier))
	case ICombinationPowerTotal, IForwardPowerTotal:
		if len(d.Value) != 4 {
			return nil, ErrorInvalidDataLength
		}
		return d.Identifier.valeOfkWh(d.Value), nil
	case IVoltage:
		return d.Identifier.valueOfVoltageSeparately(d.Value)
	case IElectricCurrent:
		return d.Identifier.valueOfElectricCurrentSeparately(d.Value)
	case IHavePowerRate:
		return d.Identifier.valueOfPowerRateSeparately(d.Value)
	}
}

// ParseData 解析数据域
func ParseData(raw []byte) (data Data, err error) {
	if len(raw) < 4 {
		err = ErrorUnknownIdentifier
		return
	}

	buf := make([]byte, len(raw))
	for i := 0; i < len(raw); i++ {
		buf[i] = raw[i] - 0x33
	}

	data = Data{
		Identifier: Identifier(binary.LittleEndian.Uint32(buf[:4])),
		Value:      buf[4:],
	}

	return
}

func (d Data) Bytes() []byte {
	w := pool.GetBuffer()
	defer pool.PutBuffer(w)

	b := d.Identifier.Bytes()
	for i := 0; i < len(b); i++ {
		w.WriteByte(b[i] + 0x33)
	}

	for i := 0; i < len(d.Value); i++ {
		w.WriteByte(d.Value[i] + 0x33)
	}

	return w.Bytes()
}
