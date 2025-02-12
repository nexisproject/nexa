// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-14, by liasica

package dlt6452007

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"

	"orba.plus/nexa/pkg/convert"
)

type Unit int

const (
	UnitkWh Unit = iota + 1 // 电量 (kWh)
	UnitV                   // 电压 (V)
	UnitA                   // 电流 (A)
	UnitkW                  // 功率 (kW)
)

func (u Unit) String() string {
	switch u {
	default:
		return " - "
	case UnitkWh:
		return "kWh"
	case UnitV:
		return "V"
	case UnitA:
		return "A"
	case UnitkW:
		return "kW"
	}
}

type DataResult struct {
	Identifier Identifier `json:"identifier"` // 数据标识符
	Value      any        `json:"value"`
	Unit       Unit       `json:"unit"`
}

// Identifier 数据标识符
// 小端序Uint32
// 排序方式为: DI3 DI2 DI1 DI0
//
// DI3标识符	    对应数据类型
//
//	00	        电能量
//	01	        最大需量及发生时间
//	02	        变量数据 （遥测等）
//	03	        事件记录
//	04	        参变量数据
//	05	        冻结量
//	06	        负荷记录
type Identifier uint32

var (
	ICombinationPowerTotal Identifier = 0x00000000 // (当前)组合有功总电能（4字节，单位：0.01kWh）
	IForwardPowerTotal     Identifier = 0x00010000 // (当前)正向有功总电能（4字节，单位：0.01kWh）
	IForwardMaxDemand      Identifier = 0x01010000 // (当前)正向有功总最大需量及发生时间 // TODO: 暂未支持：待实现
	IVoltage               Identifier = 0x0201FF00 // 电压数据块，分别是 A、B、C三相电压（每相2字节，共6字节，单位：0.1V）
	IElectricCurrent       Identifier = 0x0202FF00 // 电流数据块，分别是 A、B、C三相电流（每相3字节，共9字节，单位：0.001A）
	IHavePowerRate         Identifier = 0x0203FF00 // 瞬时有功功率数据块，分别是 总、A相、B相、C相有功功率（每个3字节，共12字节，单位：0.0001kW）
)

func (di Identifier) Bytes() (b []byte) {
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(di))
	return
}

func (di Identifier) String() string {
	return fmt.Sprintf("%08x", uint32(di))
}

// 分别解析float64
// n: 解析值数量
// m: 每个值字节数量
// factor: 除数
func parseFloat64Separately(b []byte, n int, m int, factor float64) ([]float64, error) {
	arr := make([]float64, n)
	if len(b) != n*m {
		return nil, ErrorInvalidDataLength
	}

	for i := 0; i < n*m; i += m {
		str := hex.EncodeToString(convert.Reverse(b[i : i+m]))
		v, _ := strconv.ParseFloat(str, 64)
		arr[i/m] = v / factor
	}

	return arr, nil
}

// 解析千瓦时
func (di Identifier) valeOfkWh(b []byte) *DataResult {
	str := hex.EncodeToString(convert.Reverse(b))
	v, _ := strconv.ParseFloat(str, 64)
	// float64(bcd.ToUint64(b)) / 100.0
	return &DataResult{
		Identifier: di,
		Value:      v / 100.0,
		Unit:       UnitkWh,
	}
}

// 解析电压 (三相)
func (di Identifier) valueOfVoltageSeparately(b []byte) (*DataResult, error) {
	arr, err := parseFloat64Separately(b, 3, 2, 10.0)
	if err != nil {
		return nil, err
	}
	return &DataResult{
		Identifier: di,
		Value:      arr,
		Unit:       UnitV,
	}, nil
}

// 解析电流 (三相)
func (di Identifier) valueOfElectricCurrentSeparately(b []byte) (*DataResult, error) {
	arr, err := parseFloat64Separately(b, 3, 3, 1000.0)
	if err != nil {
		return nil, err
	}
	return &DataResult{
		Identifier: di,
		Value:      arr,
		Unit:       UnitA,
	}, nil
}

// 解析瞬时总有功功率
func (di Identifier) valueOfPowerRateSeparately(b []byte) (*DataResult, error) {
	arr, err := parseFloat64Separately(b, 4, 3, 10000.0)
	if err != nil {
		return nil, err
	}
	return &DataResult{
		Identifier: di,
		Value:      arr,
		Unit:       UnitkW,
	}, nil
}
