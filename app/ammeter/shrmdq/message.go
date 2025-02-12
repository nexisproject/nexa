// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-12, by liasica

package shrmdq

import (
	"fmt"

	"orba.plus/nexa/app/ammeter"
	"orba.plus/nexa/app/ammeter/dlt6452007"
	"orba.plus/nexa/pkg/convert"
	"orba.plus/nexa/pkg/silk"
)

type Message interface {
	Decode(mo *Montage) (err error)
}

// RegisterMessage 注册数据
// DLT645-2007协议解析只需要地址，其他数据可无视
type RegisterMessage struct {
	Version string `json:"version"` // {17} 版本号
	IMEI    string `json:"imei"`    // {15} 模组号IMEI
	ICCID   string `json:"iccid"`   // {20} SIM卡号ICCID
	Signal  byte   `json:"signal"`  // 信号强度，0~31，31最强
	Address string `json:"address"` // 645地址
}

func (m *RegisterMessage) Decode(mo *Montage) (err error) {
	if len(mo.decomposed) < 4 || len(mo.decomposed[3]) != 1 {
		return ammeter.ErrorInvalidDataLength
	}

	m.Version = convert.Bytes2String(mo.decomposed[0])
	m.IMEI = convert.Bytes2String(mo.decomposed[1])
	m.ICCID = convert.Bytes2String(mo.decomposed[2])
	m.Signal = mo.decomposed[3][0]

	var data *dlt6452007.DLT6452007
	data, err = dlt6452007.Parse(mo.decomposed[4])
	if err != nil {
		return
	}
	m.Address = data.Address
	return
}

// HeartbeatMessage 心跳数据
// 只解析信号强度，其他数据可无视
type HeartbeatMessage struct {
	Signal byte `json:"signal"` // 信号强度，0~31，31最强
}

func (m *HeartbeatMessage) Decode(mo *Montage) (err error) {
	if len(mo.decomposed) < 1 || len(mo.decomposed[0]) != 1 {
		return ammeter.ErrorInvalidDataLength
	}

	m.Signal = mo.decomposed[0][0]
	return
}

// SignalMessage 信号强度数据
type SignalMessage struct {
	Intensity byte `json:"intensity"` // 信号强度，0~31，31最强
}

func (m *SignalMessage) Decode(mo *Montage) (err error) {
	if len(mo.decomposed) < 1 || len(mo.decomposed[0]) != 1 {
		return ammeter.ErrorInvalidDataLength
	}

	m.Intensity = mo.decomposed[0][0]
	return
}

// PassthroughMessage 透传数据
type PassthroughMessage struct {
	Response [][]byte                 `json:"response"` // 透传应答原始数据
	Parsed   []*dlt6452007.DLT6452007 `json:"parsed"`   // 透传应答解析后数据
}

func (m *PassthroughMessage) Decode(mo *Montage) (err error) {
	if len(mo.decomposed) < 1 {
		return ammeter.ErrorInvalidData
	}

	m.Response = make([][]byte, mo.decomposed[0][0]&0xFF)
	m.Parsed = make([]*dlt6452007.DLT6452007, len(m.Response))
	for i := 1; i < len(mo.decomposed); i++ {
		m.Response[i-1] = mo.decomposed[i]
		var data *dlt6452007.DLT6452007
		data, err = dlt6452007.Parse(mo.decomposed[i])
		if err != nil {
			return
		}
		m.Parsed[i-1] = data
	}

	for _, p := range m.Parsed {
		var r *dlt6452007.DataResult
		r, err = p.Data.ParseValue()
		if err != nil {
			continue
		}

		// 发送kafka消息
		km := &ammeter.KafkaMessage{
			No: p.Address,
		}
		switch r.Identifier {
		default:
			continue
		case dlt6452007.IForwardPowerTotal:
			km.PowerTotal = silk.Pointer(r.Value.(float64))
		case dlt6452007.IVoltage:
			km.Voltage = &ammeter.VoltageData{
				A: r.Value.([]float64)[0],
				B: r.Value.([]float64)[1],
				C: r.Value.([]float64)[2],
			}
		case dlt6452007.IElectricCurrent:
			km.Electric = &ammeter.ElectricData{
				A: r.Value.([]float64)[0],
				B: r.Value.([]float64)[1],
				C: r.Value.([]float64)[2],
			}
		}
		if km.PowerTotal != nil || km.Voltage != nil || km.Electric != nil {
			ammeter.SendKafkaMessage(km)
		}

		fmt.Printf("DLT645-2007: %s -> %v(%s)\n", p.Address, r.Value, r.Unit)
	}

	return
}
